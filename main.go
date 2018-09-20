package main

import (
	"app/tax"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = initDatabase()

type ResponseJSON struct {
	Header ResponseHeader `json:"header"`
	Data   interface{}    `json:"data"`
}

type ResponseHeader struct {
	Message   string `json:"message"`
	Reason    string `json:"reason"`
	ErrorCode int    `json:"error_code"`
}

type OrderDetail struct {
	tax.Tax
	Name        string
	Amount      int
	TaxAmount   float64
	TotalAmount float64
}

type Order struct {
	Detail     []OrderDetail
	Total      int
	TotalTax   float64
	GrandTotal float64
}

func handleCalculation(w http.ResponseWriter, r *http.Request) {
	var order Order

	rows, err := db.Query("SELECT item_name, item_tax_code, item_amount FROM order_detail")
	if err != nil {
		log.Println("Error querying select to db. Error :", err)
	}

	for rows.Next() {
		var detail OrderDetail
		if err = rows.Scan(&detail.Name, &detail.Code, &detail.Amount); err != nil {
			log.Println("Error scan query result. Error :", err)
		}

		detail.SetTaxType()
		detail.TaxAmount = detail.CalculateTax(float64(detail.Amount))
		detail.TotalAmount = float64(detail.Amount) + detail.TaxAmount

		order.Detail = append(order.Detail, detail)
		order.Total = order.Total + detail.Amount
		order.TotalTax = order.TotalTax + detail.TaxAmount
		order.GrandTotal = order.GrandTotal + detail.TotalAmount
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, order)
}

func doCalculation(w http.ResponseWriter, r *http.Request) {
	itemName := r.FormValue("item-name")
	if itemName == "" {
		writeResponse(w, nil, http.StatusBadRequest, "Invalid Item Name", "Item name must not empty")
		return
	}

	itemTaxCode, err := strconv.ParseInt(r.FormValue("item-tax-code"), 10, 64)
	if err != nil {
		writeResponse(w, nil, http.StatusBadRequest, "Invalid Tax Code", "Tax code must be in numeric format")
		return
	}

	itemAmount, err := strconv.ParseInt(r.FormValue("item-amount"), 10, 64)
	if err != nil {
		writeResponse(w, nil, http.StatusBadRequest, "Invalid Amount", "Amount must be in numeric format")
		return
	}

	query := `INSERT INTO order_detail(item_name, item_tax_code, item_amount) VALUES (?, ?, ?)`
	_, err = db.Query(query, itemName, itemTaxCode, itemAmount)
	if err != nil {
		log.Println("Error querying insert to db. Error :", err)
		writeResponse(w, nil, http.StatusInternalServerError, "Unable to insert data into database", err.Error())
		return
	}

	response := ResponseJSON{
		Data: "Data successfully inserted into database",
	}

	writeResponse(w, response, http.StatusOK, "", "")
}

func initDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(mysql-db:3306)/tax-calc")
	if err != nil {
		panic(fmt.Sprintf("Error opening db connection. Error : %s", err.Error()))
	}

	return db
}

func main() {
	http.HandleFunc("/", handleCalculation)
	http.HandleFunc("/calculate", doCalculation)
	http.ListenAndServe(":80", nil)
}

//This is util function for writing api response
func writeResponse(w http.ResponseWriter, data interface{}, statusCode int, message string, reason string) {
	w.Header().Set("Content-Type", "application/json")
	response := ResponseJSON{
		Header: ResponseHeader{
			Message:   message,
			Reason:    reason,
			ErrorCode: statusCode,
		},
		Data: data,
	}
	resp, err := json.Marshal(response)
	if err != nil {
		log.Println("[writeResponse] Error marshalling response. Error :", err)
	}
	w.WriteHeader(statusCode)
	w.Write(resp)
}
