package doc

import "github.com/invopop/gobl/bill"

// InvoiceReference is used for invoice modification (reference other invoice)
type InvoiceReference struct {
	OriginalInvoiceNumber string `xml:"originalInvoiceNumber"`
	ModifyWithoutMaster   bool   `xml:"modifyWithoutMaster"`
	ModificationIndex     string `xml:"modificationIndex"`
}

func newInvoiceReference(reference *bill.Preceding) *InvoiceReference {
	return &InvoiceReference{
		OriginalInvoiceNumber: reference.Code,
		ModifyWithoutMaster:   false,
		ModificationIndex:     reference.Series,
	}
}

/*
ModifyWithoutMaster value is only true in specific cases:
See the p.99 of https://onlineszamla-test.nav.gov.hu/files/container/download/Online%20Invoice%20System%203.0%20Interface%20Specification.pdf
This can be handled with an extension code
*/
