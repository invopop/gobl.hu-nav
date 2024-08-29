package doc

import "github.com/invopop/gobl/bill"

// InvoiceMain contains the invoice data.
// It can be used for both invoice creation and modification.
// It can have 2 values: Invoice and batchInvoice
type InvoiceMain struct {
	Invoice *Invoice `xml:"invoice"`
	//BatchInvoice *BatchInvoice `xml:"batchInvoice"` // Used only for batch modifications
}

// Invoice is the main invoice data structure.
type Invoice struct {
	//InvoiceReference  InvoiceReference  `xml:"invoiceReference,omitempty"` // Used for invoice modification (reference other invoice)
	InvoiceHead  *InvoiceHead  `xml:"invoiceHead"`
	InvoiceLines *InvoiceLines `xml:"invoiceLines,omitempty"`
	//ProductFeeSummary ProductFeeSummary `xml:"productFeeSummary,omitempty"`

	InvoiceSummary *InvoiceSummary `xml:"invoiceSummary"`
}

func newInvoiceMain(inv *bill.Invoice) (*InvoiceMain, error) {
	invoiceHead, err := newInvoiceHead(inv)
	if err != nil {
		return nil, err
	}

	invoiceLines, err := newInvoiceLines(inv)
	if err != nil {
		return nil, err
	}

	invoiceSummary, err := newInvoiceSummary(inv)
	if err != nil {
		return nil, err
	}

	return &InvoiceMain{
		Invoice: &Invoice{
			InvoiceHead:    invoiceHead,
			InvoiceLines:   invoiceLines,
			InvoiceSummary: invoiceSummary,
		},
	}, nil
}
