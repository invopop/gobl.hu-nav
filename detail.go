package nav

type InvoiceDetail struct {
	InvoiceCategory               string `xml:"invoiceCategory"`
	InvoiceDeliveryDate           string `xml:"invoiceDeliveryDate"`
	InvoiceDeliveryPeriodStart    string `xml:"invoiceDeliveryPeriodStart,omitempty"`
	InvoiceDeliveryPeriodEnd      string `xml:"invoiceDeliveryPeriodEnd,omitempty"`
	InvoiceAccountingDeliveryDate string `xml:"invoiceAccountingDeliveryDate,omitempty"`
	PeriodicalSettlement          bool   `xml:"periodicalSettlement,omitempty"`
	SmallBusinessIndicator        bool   `xml:"smallBusinessIndicator,omitempty"`
	CurrencyCode                  string `xml:"currencyCode"`
	ExchangeRate                  string `xml:"exchangeRate,omitempty"`
	UtilitySettlementIndicator    bool   `xml:"utilitySettlementIndicator,omitempty"`
	SelfBillingIndicator          bool   `xml:"selfBillingIndicator,omitempty"`
	PaymentMethod                 string `xml:"paymentMethod,omitempty"`
	PaymentDate                   string `xml:"paymentDate,omitempty"`
	CashAccountingIndicator       bool   `xml:"cashAccountingIndicator,omitempty"`
	InvoiceAppearance             string `xml:"invoiceAppearance"`
	//Some more optional data
}
