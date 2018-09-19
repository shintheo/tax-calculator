package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

const (
	FOOD_TAX_CODE          = 1
	TOBACCO_TAX_CODE       = 2
	ENTERTAINMENT_TAX_CODE = 3
)

type ResponseJSON struct {
	Data interface{} `json:"data"`
}

type OrderDetail struct {
	Name        string `db:item_name`
	TaxCode     int    `db:item_tax_code`
	Type        string
	Amount      int `db:item_amount`
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

	db, err := sql.Open("mysql", "root:password@tcp(mysql-db:3306)/tax-calc")
	if err != nil {
		panic(fmt.Sprintf("Error opening db connection. Error : %s", err.Error()))
	}

	defer db.Close()

	rows, err := db.Query("SELECT item_name, item_tax_code, item_amount FROM order_detail")
	if err != nil {
		log.Println("Error querying select to db. Error :", err)
	}

	for rows.Next() {
		var detail OrderDetail
		if err = rows.Scan(&detail.Name, &detail.TaxCode, &detail.Amount); err != nil {
			log.Println("Error scan query result. Error :", err)
		}

		detail.Type = getTaxType(detail.TaxCode)
		detail.TaxAmount = calculateTax(detail.TaxCode, float64(detail.Amount))
		detail.TotalAmount = float64(detail.Amount) + detail.TaxAmount

		order.Detail = append(order.Detail, detail)
		order.Total = order.Total + detail.Amount
		order.TotalTax = order.TotalTax + detail.TaxAmount
		order.GrandTotal = order.GrandTotal + detail.TotalAmount
	}

	fmt.Println("Order :", order)
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, order)
}

func doCalculation(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@tcp(mysql-db:3306)/tax-calc")
	if err != nil {
		panic(fmt.Sprintf("Error opening db connection. Error : %s", err.Error()))
	}

	defer db.Close()

	itemName := r.FormValue("item-name")
	itemTaxCode := r.FormValue("item-tax-code")
	itemAmount := r.FormValue("item-amount")

	query := `INSERT INTO order_detail(item_name, item_tax_code, item_amount) VALUES (?, ?, ?)`
	_, err = db.Query(query, itemName, itemTaxCode, itemAmount)
	if err != nil {
		log.Println("Error querying insert to db. Error :", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := ResponseJSON{
		Data: "OK",
	}
	e, _ := json.Marshal(response)
	w.Write(e)
}

func main() {
	http.HandleFunc("/", handleCalculation)
	http.HandleFunc("/calculate", doCalculation)
	http.ListenAndServe(":80", nil)
}

func getTaxType(taxCode int) string {
	var taxType string
	if taxCode == FOOD_TAX_CODE {
		taxType = "Food"
	} else if taxCode == TOBACCO_TAX_CODE {
		taxType = "Tobacco"
	} else if taxCode == ENTERTAINMENT_TAX_CODE {
		taxType = "Entertainment"
	}
	return taxType
}

func calculateTax(taxCode int, amount float64) float64 {
	var taxAmount float64
	if taxCode == FOOD_TAX_CODE {
		taxAmount = 0.1 * amount //10% of amount
	} else if taxCode == TOBACCO_TAX_CODE {
		taxAmount = 10 + (0.02 * amount) //10 + (2% of value)
	} else if taxCode == ENTERTAINMENT_TAX_CODE {
		if amount >= 100 {
			taxAmount = 0.01 * (amount - 100) //1% of (value - 100)
		} else {
			fmt.Println("If entertainment is tax-free : ", taxAmount)
		}
	}
	return taxAmount
}
