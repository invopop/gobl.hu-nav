package doc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/tax"
)

// InvoiceDetail contains the invoice detail data
type InvoiceDetail struct {
	InvoiceCategory            string `xml:"invoiceCategory"` //NORMAL, SIMPLIFIED, AGGREGATE
	InvoiceDeliveryDate        string `xml:"invoiceDeliveryDate"`
	InvoiceDeliveryPeriodStart string `xml:"invoiceDeliveryPeriodStart,omitempty"`
	InvoiceDeliveryPeriodEnd   string `xml:"invoiceDeliveryPeriodEnd,omitempty"`
	//InvoiceAccountingDeliveryDate string `xml:"invoiceAccountingDeliveryDate,omitempty"`
	//PeriodicalSettlement          bool   `xml:"periodicalSettlement,omitempty"`
	//SmallBusinessIndicator        bool   `xml:"smallBusinessIndicator,omitempty"`
	CurrencyCode string  `xml:"currencyCode"`
	ExchangeRate float64 `xml:"exchangeRate"`
	//UtilitySettlementIndicator bool   `xml:"utilitySettlementIndicator,omitempty"`
	//SelfBillingIndicator       bool   `xml:"selfBillingIndicator,omitempty"`
	PaymentMethod string `xml:"paymentMethod,omitempty"`
	PaymentDate   string `xml:"paymentDate,omitempty"`
	//CashAccountingIndicator    bool   `xml:"cashAccountingIndicator,omitempty"`
	InvoiceAppearance string `xml:"invoiceAppearance"` // PAPER, ELECTRONIC, EDI, UNKNOWN
	//Some more optional data
}

// NewInvoiceDetail creates a new InvoiceDetail from an invoice
func NewInvoiceDetail(inv *bill.Invoice) (*InvoiceDetail, error) {
	category := "NORMAL"
	if inv.Tax.ContainsTag(tax.TagSimplified) {
		category = "SIMPLIFIED"
	}

	date := &inv.IssueDate
	periodStart := ""
	periodEnd := ""
	if inv.Delivery != nil {
		if inv.Delivery.Date != nil {
			date = inv.Delivery.Date
		}
		if inv.Delivery.Period != nil {
			periodStart = inv.Delivery.Period.Start.String()
			periodEnd = inv.Delivery.Period.End.String()
		}
	}

	dueDate := ""
	paymentMethod := ""
	if inv.Payment != nil {
		if inv.Payment.Terms != nil {
			if inv.Payment.Terms.DueDates != nil {
				if len(inv.Payment.Terms.DueDates) > 0 {
					dueDate = inv.Payment.Terms.DueDates[0].Date.String()
				}
			}
		}
		if inv.Payment.Instructions != nil {
			switch inv.Payment.Instructions.Key {
			case "cash":
				paymentMethod = "CASH"
			case "credit-transfer":
				paymentMethod = "TRANSFER"
			case "debit-transfer":
				paymentMethod = "TRANSFER"
			case "card":
				paymentMethod = "CARD"
			// There is one case that is VOUCHER
			default:
				paymentMethod = "OTHER"
			}
		}
	}

	rate, err := getInvoiceRate(inv)
	if err != nil {
		return nil, err
	}

	return &InvoiceDetail{
		InvoiceCategory:            category,
		InvoiceDeliveryDate:        date.String(),
		InvoiceDeliveryPeriodStart: periodStart,
		InvoiceDeliveryPeriodEnd:   periodEnd,
		CurrencyCode:               inv.Currency.String(),
		ExchangeRate:               rate,
		PaymentMethod:              paymentMethod,
		PaymentDate:                dueDate,
		InvoiceAppearance:          "EDI",
	}, nil
}

func getInvoiceRate(inv *bill.Invoice) (float64, error) {
	if inv.Currency == currency.HUF {
		return 1.0, nil
	}

	for _, ex := range inv.ExchangeRates {
		if ex.To == currency.HUF {
			return ex.Amount.Rescale(6).Float64(), nil
		}
	}

	return -1.0, ErrNoExchangeRate
}
