package tax

import "testing"

func TestSetTaxType(t *testing.T) {
	var dummyTax Tax
	expectedTaxType := map[int]string{
		1: "Food",
		2: "Tobacco",
		3: "Entertainment",
	}

	for code, taxType := range expectedTaxType {
		dummyTax.Code = code
		dummyTax.SetTaxType()

		if dummyTax.Type != taxType {
			t.Fatalf("Expected %s, got %s", taxType, dummyTax.Type)
		}
	}
}

func TestCalculateTax(t *testing.T) {
	var dummyTax Tax
	var dummyAmount float64 = 1000

	//expected result with assumption amount is 1000 for all code
	expectedResult := map[int]float64{
		1: 100, // 10% of amount
		2: 30,  // 10 + (2% of amount)
		3: 9,   // 1% of (amount - 100), because 1000 is greater than 100
	}

	for code, tax := range expectedResult {
		dummyTax.Code = code
		result := dummyTax.CalculateTax(dummyAmount)

		if result != tax {
			t.Fatalf("Expected %f, got %f", tax, result)
		}
	}

	//Test case for entertainment tax which amount is below 100
	dummyTax.Code = ENTERTAINMENT_TAX_CODE
	dummyAmount = 99
	if result := dummyTax.CalculateTax(dummyAmount); result != 0 {
		t.Fatalf("Expected %d, got %f", 0, result)
	}

}
