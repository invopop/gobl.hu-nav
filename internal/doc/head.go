package doc

import "github.com/invopop/gobl/bill"

// InvoiceHead contains data pertaining to the invoice as a whole.
type InvoiceHead struct {
	SupplierInfo *SupplierInfo `xml:"supplierInfo"`
	CustomerInfo *CustomerInfo `xml:"customerInfo,omitempty"`
	//FiscalRepresentativeInfo FiscalRepresentativeInfo `xml:"fiscalRepresentativeInfo,omitempty"`

	InvoiceDetail *InvoiceDetail `xml:"invoiceDetail"`
}

func newInvoiceHead(inv *bill.Invoice) (*InvoiceHead, error) {
	supplierInfo := newSupplierInfo(inv.Supplier)

	customerInfo := newCustomerInfo(inv.Customer)

	detail, err := newInvoiceDetail(inv)
	if err != nil {
		return nil, err
	}
	return &InvoiceHead{
		SupplierInfo:  supplierInfo,
		CustomerInfo:  customerInfo,
		InvoiceDetail: detail,
	}, nil
}
