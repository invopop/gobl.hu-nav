package doc

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNewInvoiceDetail_NormalInvoice(t *testing.T) {
	invoice := &bill.Invoice{
		Currency: currency.USD,
		ExchangeRates: []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.HUF,
				Amount: num.AmountFromFloat64(358.3543210123, 5),
			},
		},
		IssueDate: *cal.NewDate(2023, 8, 10),
		Payment: &bill.Payment{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{Date: cal.NewDate(2023, 8, 20)},
				},
			},
			Instructions: &pay.Instructions{
				Key: "card",
			},
		},
	}

	detail, err := newInvoiceDetail(invoice)
	assert.NoError(t, err)
	assert.Equal(t, "NORMAL", detail.InvoiceCategory)
	assert.Equal(t, "2023-08-10", detail.InvoiceDeliveryDate)
	assert.Equal(t, "", detail.InvoiceDeliveryPeriodStart)
	assert.Equal(t, "", detail.InvoiceDeliveryPeriodEnd)
	assert.Equal(t, "USD", detail.CurrencyCode)
	assert.Equal(t, 358.35432, detail.ExchangeRate)
	assert.Equal(t, "CARD", detail.PaymentMethod)
	assert.Equal(t, "2023-08-20", detail.PaymentDate)
	assert.Equal(t, "EDI", detail.InvoiceAppearance)
}

func TestNewInvoiceDetail_SimplifiedInvoice(t *testing.T) {
	invoice := &bill.Invoice{
		Currency: currency.USD,
		ExchangeRates: []*currency.ExchangeRate{
			{
				From:   currency.USD,
				To:     currency.HUF,
				Amount: num.AmountFromFloat64(358.35, 5),
			},
		},
		IssueDate: *cal.NewDate(2023, 7, 15),
		Tax: &bill.Tax{
			Tags: []cbc.Key{tax.TagSimplified},
		},
	}

	detail, err := newInvoiceDetail(invoice)
	assert.NoError(t, err)
	assert.Equal(t, "SIMPLIFIED", detail.InvoiceCategory)
}

func TestNewInvoiceDetail_NoExchangeRate(t *testing.T) {
	invoice := &bill.Invoice{
		Currency: currency.JPY,
		ExchangeRates: []*currency.ExchangeRate{
			{
				From:   currency.JPY,
				To:     currency.USD,
				Amount: num.AmountFromFloat64(0.0068, 5),
			},
		},
		IssueDate: *cal.NewDate(2023, 7, 15),
	}

	_, err := newInvoiceDetail(invoice)
	assert.Error(t, err)
	assert.Equal(t, ErrNoExchangeRate, err)
}
