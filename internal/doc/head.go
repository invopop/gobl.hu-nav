package doc

import "github.com/invopop/gobl/bill"

// InvoiceHead contains data pertaining to the invoice as a whole.
type InvoiceHead struct {
	SupplierInfo *SupplierInfo `xml:"supplierInfo"`
	CustomerInfo *CustomerInfo `xml:"customerInfo,omitempty"`
	//FiscalRepresentativeInfo FiscalRepresentativeInfo `xml:"fiscalRepresentativeInfo,omitempty"`

	InvoiceDetail *InvoiceDetail `xml:"invoiceDetail"`
}

func NewInvoiceHead(inv *bill.Invoice) (*InvoiceHead, error) {
	supplierInfo, err := NewSupplierInfo(inv.Supplier)
	if err != nil {
		return nil, err
	}

	customerInfo, err := NewCustomerInfo(inv.Customer)
	if err != nil {
		return nil, err
	}

	detail, err := NewInvoiceDetail(inv)
	if err != nil {
		return nil, err
	}
	return &InvoiceHead{
		SupplierInfo:  supplierInfo,
		CustomerInfo:  customerInfo,
		InvoiceDetail: detail,
	}, nil
}
