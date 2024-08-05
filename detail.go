package nav

import (
	"math"
	"strconv"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
)

type InvoiceDetail struct {
	InvoiceCategory     string `xml:"invoiceCategory"` //NORMAL, SIMPLIFIED, AGGREGATE
	InvoiceDeliveryDate string `xml:"invoiceDeliveryDate"`
	//InvoiceDeliveryPeriodStart    string `xml:"invoiceDeliveryPeriodStart,omitempty"`
	//InvoiceDeliveryPeriodEnd      string `xml:"invoiceDeliveryPeriodEnd,omitempty"`
	//InvoiceAccountingDeliveryDate string `xml:"invoiceAccountingDeliveryDate,omitempty"`
	//PeriodicalSettlement          bool   `xml:"periodicalSettlement,omitempty"`
	//SmallBusinessIndicator        bool   `xml:"smallBusinessIndicator,omitempty"`
	CurrencyCode string  `xml:"currencyCode"`
	ExchangeRate float64 `xml:"exchangeRate"`
	//UtilitySettlementIndicator bool   `xml:"utilitySettlementIndicator,omitempty"`
	//SelfBillingIndicator       bool   `xml:"selfBillingIndicator,omitempty"`
	//PaymentMethod              string `xml:"paymentMethod,omitempty"`
	//PaymentDate                string `xml:"paymentDate,omitempty"`
	//CashAccountingIndicator    bool   `xml:"cashAccountingIndicator,omitempty"`
	InvoiceAppearance string `xml:"invoiceAppearance"` // PAPER, ELECTRONIC, EDI, UNKNOWN
	//Some more optional data
}

// NewInvoiceDetail creates a new InvoiceDetail from an invoice
func NewInvoiceDetail(inv *bill.Invoice) (*InvoiceDetail, error) {
	rate, err := getInvoiceRate(inv)
	if err != nil {
		return nil, err
	}

	return &InvoiceDetail{
		InvoiceCategory:     "NORMAL",
		InvoiceDeliveryDate: inv.OperationDate.String(),
		CurrencyCode:        inv.Currency.String(),
		ExchangeRate:        formatRate(rate),
		InvoiceAppearance:   "EDI",
	}, nil
}

func getInvoiceRate(inv *bill.Invoice) (float64, error) {
	if inv.Currency == currency.HUF {
		return 1.0, nil
	}

	for _, ex := range inv.ExchangeRates {
		if ex.To == currency.HUF {
			return ex.Amount.Float64(), nil
		}
	}

	return -1.0, ErrNoExchangeRate
}

func formatRate(value float64) float64 {
	// Check if the float64 number has more than 6 decimal digits
	if hasMoreThanSixDecimalDigits(value) {
		return math.Round(value*1000000) / 1000000
	}

	// Convert the float64 number to a string without trailing zeros
	return value
}

// hasMoreThanSixDecimalDigits checks if a float64 number has more than 6 decimal digits
func hasMoreThanSixDecimalDigits(value float64) bool {
	// Separate the fractional part from the integer part
	fractionalPart := value - math.Floor(value)

	// Convert the fractional part to a string
	fractionalStr := strconv.FormatFloat(fractionalPart, 'f', -1, 64)

	// Remove the leading "0." from the string representation
	fractionalStr = strings.TrimPrefix(fractionalStr, "0.")

	// Check the length of the fractional part
	return len(fractionalStr) > 6
}
