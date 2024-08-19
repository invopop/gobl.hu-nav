package nav

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNewInvoiceLines(t *testing.T) {
	invoice := &bill.Invoice{
		Currency: currency.HUF,
		Lines: []*bill.Line{
			{
				Index: 1,
				Item: &org.Item{
					Identities: []*org.Identity{
						{Type: "VTSZ", Code: cbc.Code("12345")},
					},
					Unit: org.UnitKilogram,
					Key:  "PRODUCT",
				},
				Total: num.AmountFromFloat64(1000, 2),
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Percent: num.NewPercentage(27, 4)},
				},
			},
		},
		Tax: &bill.Tax{Tags: []cbc.Key{tax.TagSimplified}},
	}

	invoiceLines, err := NewInvoiceLines(invoice)
	assert.NoError(t, err)
	assert.NotNil(t, invoiceLines)
	assert.False(t, invoiceLines.MergedItemIndicator)
	assert.Len(t, invoiceLines.Lines, 1)

	line := invoiceLines.Lines[0]
	assert.Equal(t, 1, line.LineNumber)
	assert.NotNil(t, line.ProductCodes)
	assert.Equal(t, "KILOGRAM", line.UnitOfMeasure)
	assert.Equal(t, "PRODUCT", line.LineNatureIndicator)
	assert.NotNil(t, line.LineAmountsSimplified)
	assert.Nil(t, line.LineAmountsNormal)
}

func TestNewLine_NormalInvoice(t *testing.T) {
	line := &bill.Line{
		Index: 1,
		Item: &org.Item{
			Identities: []*org.Identity{
				{Type: "VTSZ", Code: cbc.Code("12345")},
			},
			Unit: org.UnitKilogram,
			Key:  "PRODUCT",
		},
		Total: num.AmountFromFloat64(1000, 2),
		Taxes: tax.Set{
			{Category: tax.CategoryVAT, Percent: num.NewPercentage(27, 4)},
		},
	}

	taxInfo := &taxInfo{}
	rate := 1.0

	lineNav, err := NewLine(line, taxInfo, rate)
	assert.NoError(t, err)
	assert.NotNil(t, lineNav)
	assert.Equal(t, 1, lineNav.LineNumber)
	assert.Equal(t, "KILOGRAM", lineNav.UnitOfMeasure)
	assert.Equal(t, "PRODUCT", lineNav.LineNatureIndicator)
	assert.NotNil(t, lineNav.LineAmountsNormal)
	assert.Nil(t, lineNav.LineAmountsSimplified)
}

func TestNewLine_SimplifiedInvoice(t *testing.T) {
	line := &bill.Line{
		Index: 1,
		Item: &org.Item{
			Identities: []*org.Identity{
				{Type: "VTSZ", Code: cbc.Code("12345")},
			},
			Unit: org.UnitKilogram,
			Key:  "PRODUCT",
		},
		Total: num.AmountFromFloat64(1000, 2),
		Taxes: tax.Set{
			{Category: tax.CategoryVAT, Percent: num.NewPercentage(27, 4)},
		},
	}

	taxInfo := &taxInfo{simplifiedInvoice: true}
	rate := 1.0

	lineNav, err := NewLine(line, taxInfo, rate)
	assert.NoError(t, err)
	assert.NotNil(t, lineNav)
	assert.Equal(t, 1, lineNav.LineNumber)
	assert.Equal(t, "KILOGRAM", lineNav.UnitOfMeasure)
	assert.Equal(t, "PRODUCT", lineNav.LineNatureIndicator)
	assert.NotNil(t, lineNav.LineAmountsSimplified)
	assert.Nil(t, lineNav.LineAmountsNormal)
}

func TestNewProductCodes(t *testing.T) {
	identities := []*org.Identity{
		{Type: "VTSZ", Code: cbc.Code("12345")},
		{Type: "OWN", Code: cbc.Code("OWN123")},
	}

	productCodes := NewProductCodes(identities)
	assert.NotNil(t, productCodes)
	assert.Len(t, productCodes.ProductCode, 2)

	assert.Equal(t, "VTSZ", productCodes.ProductCode[0].ProductCodeCategory)
	assert.Equal(t, "12345", productCodes.ProductCode[0].ProductCodeValue)

	assert.Equal(t, "OWN", productCodes.ProductCode[1].ProductCodeCategory)
	assert.Equal(t, "OWN123", productCodes.ProductCode[1].ProductCodeOwnValue)
}

func TestNewProductCode(t *testing.T) {
	identity := &org.Identity{Type: "VTSZ", Code: cbc.Code("12345")}
	productCode := NewProductCode(identity)
	assert.NotNil(t, productCode)
	assert.Equal(t, "VTSZ", productCode.ProductCodeCategory)
	assert.Equal(t, "12345", productCode.ProductCodeValue)

	identity = &org.Identity{Type: "OWN", Code: cbc.Code("OWN123")}
	productCode = NewProductCode(identity)
	assert.NotNil(t, productCode)
	assert.Equal(t, "OWN", productCode.ProductCodeCategory)
	assert.Equal(t, "OWN123", productCode.ProductCodeOwnValue)
}

func TestAmountToHUF(t *testing.T) {
	amount := num.AmountFromFloat64(1000, 2)
	exchangeRate := 300.0

	convertedAmount := amountToHUF(amount, exchangeRate)
	expectedAmount := num.AmountFromFloat64(300000, 2)

	assert.Equal(t, expectedAmount, convertedAmount)
}
