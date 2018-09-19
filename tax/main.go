package tax

type Tax struct {
	Code int
	Type string
}

const (
	FOOD_TAX_CODE          = 1
	TOBACCO_TAX_CODE       = 2
	ENTERTAINMENT_TAX_CODE = 3
)

func (tax *Tax) GetTaxType() string {
	var taxType string
	if tax.Code == FOOD_TAX_CODE {
		taxType = "Food"
	} else if tax.Code == TOBACCO_TAX_CODE {
		taxType = "Tobacco"
	} else if tax.Code == ENTERTAINMENT_TAX_CODE {
		taxType = "Entertainment"
	}

	return taxType
}

func (tax *Tax) CalculateTax(amount float64) float64 {
	var taxAmount float64
	if tax.Code == FOOD_TAX_CODE {
		taxAmount = 0.1 * amount //10% of amount
	} else if tax.Code == TOBACCO_TAX_CODE {
		taxAmount = 10 + (0.02 * amount) //10 + (2% of value)
	} else if tax.Code == ENTERTAINMENT_TAX_CODE {
		if amount >= 100 {
			taxAmount = 0.01 * (amount - 100) //1% of (value - 100)
		} else {
			taxAmount = 0
		}
	}
	return taxAmount
}
