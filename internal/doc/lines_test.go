package doc

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
	// Set up test data
	invoice := &bill.Invoice{
		Currency: currency.HUF,
		Lines: []*bill.Line{
			{
				Index:    1,
				Quantity: num.MakeAmount(2, 0),
				Item: &org.Item{
					Name:  "Test Product",
					Key:   "PRODUCT",
					Price: num.MakeAmount(10000, 2),
					Unit:  org.UnitPiece,
					Identities: []*org.Identity{
						{Type: "VTSZ", Code: cbc.Code("1234")},
					},
				},
				Sum: num.MakeAmount(20000, 2),
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Percent: num.NewPercentage(27, 4)},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason: "Seasonal Discount",
						Amount: num.MakeAmount(500, 2),
					},
				},
				Total: num.MakeAmount(19500, 2),
			},
		},
	}

	// Execute the function under test
	invoiceLines, err := NewInvoiceLines(invoice)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, invoiceLines)
	assert.False(t, invoiceLines.MergedItemIndicator)
	assert.Len(t, invoiceLines.Lines, 1)

	line := invoiceLines.Lines[0]
	assert.Equal(t, 1, line.LineNumber)
	assert.Equal(t, "Test Product", line.LineDescription)
	assert.Equal(t, "PRODUCT", line.LineNatureIndicator)
	assert.Equal(t, 2.0, line.Quantity)
	assert.Equal(t, 100.00, line.UnitPrice)
	assert.Equal(t, "PIECE", line.UnitOfMeasure)

	// Check Product Codes
	assert.NotNil(t, line.ProductCodes)
	assert.Len(t, line.ProductCodes.ProductCode, 1)
	assert.Equal(t, "VTSZ", line.ProductCodes.ProductCode[0].ProductCodeCategory)
	assert.Equal(t, "1234", line.ProductCodes.ProductCode[0].ProductCodeValue)

	// Check Discount Data
	assert.NotNil(t, line.LineDiscountData)
	assert.Equal(t, "Seasonal Discount. ", line.LineDiscountData.DiscountDescription)
	assert.Equal(t, 5.00, line.LineDiscountData.DiscountValue)

	// Check VAT and Amounts
	assert.NotNil(t, line.LineAmountsNormal)
	assert.Equal(t, 195.00, line.LineAmountsNormal.LineNetAmountData.LineNetAmount)
	assert.Equal(t, 247.65, line.LineAmountsNormal.LineGrossAmountData.LineGrossAmount) // Assuming 27% VAT
}

func TestNewLine_SimplifiedInvoice(t *testing.T) {
	// Set up test data for a simplified invoice
	invoice := &bill.Invoice{
		Lines: []*bill.Line{
			{
				Index: 1,
				Item: &org.Item{
					Name:  "Simplified Service",
					Key:   "SERVICE",
					Price: num.MakeAmount(10000, 2),
					Unit:  org.UnitHour,
				},
				Quantity: num.MakeAmount(3, 0),
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Percent: num.NewPercentage(18, 4)},
				},
				Total: num.MakeAmount(30000, 2),
			},
		},
	}

	info := &taxInfo{simplifiedInvoice: true}
	rate := 1.0

	// Execute the function under test
	line, err := NewLine(invoice.Lines[0], info, rate)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, line)

	// Check line data for a simplified invoice
	assert.Equal(t, "Simplified Service", line.LineDescription)
	assert.Equal(t, "SERVICE", line.LineNatureIndicator)
	assert.Equal(t, 3.0, line.Quantity)
	assert.Equal(t, 100.00, line.UnitPrice)
	assert.Equal(t, "HOUR", line.UnitOfMeasure)

	// Check VAT and Amounts
	assert.NotNil(t, line.LineAmountsSimplified)
	assert.Equal(t, 300.00, line.LineAmountsSimplified.LineGrossAmountSimplified)
}
