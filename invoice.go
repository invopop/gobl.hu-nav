package nav

import "github.com/invopop/gobl/bill"

// InvoiceMain can have 2 values: Invoice and batchInvoice
// For the moment, we are only going to focus on invoice
type InvoiceMain struct {
	Invoice *Invoice `xml:"invoice"`
	//BatchInvoice *BatchInvoice `xml:"batchInvoice"` // Used only for batch modifications
}

type Invoice struct {
	//InvoiceReference  InvoiceReference  `xml:"invoiceReference,omitempty"` // Used for invoice modification
	InvoiceHead  *InvoiceHead  `xml:"invoiceHead"`
	InvoiceLines *InvoiceLines `xml:"invoiceLines,omitempty"`
	//ProductFeeSummary ProductFeeSummary `xml:"productFeeSummary,omitempty"`

	InvoiceSummary *InvoiceSummary `xml:"invoiceSummary"`
}

func NewInvoiceMain(inv *bill.Invoice) (*InvoiceMain, error) {
	invoiceHead, err := NewInvoiceHead(inv)
	if err != nil {
		return nil, err
	}

	invoiceLines, err := NewInvoiceLines(inv)
	if err != nil {
		return nil, err
	}

	//invoiceSummary, err := NewInvoiceSummary(inv)
	if err != nil {
		return nil, err
	}

	return &InvoiceMain{
		Invoice: &Invoice{
			InvoiceHead:  invoiceHead,
			InvoiceLines: invoiceLines,
			//InvoiceSummary: invoiceSummary,
		},
	}, nil
}
