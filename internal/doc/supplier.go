package doc

import (
	"github.com/invopop/gobl/org"
)

// SupplierInfo contains data related to the issuer of the invoice (supplier).
type SupplierInfo struct {
	SupplierTaxNumber    *TaxNumber `xml:"supplierTaxNumber"`
	GroupMemberTaxNumber *TaxNumber `xml:"groupMemberTaxNumber,omitempty"`
	//CommunityVATNumber   string    `xml:"communityVATNumber,omitempty"` // This is just the same as Supplier Number with HU prefix
	SupplierName    string   `xml:"supplierName"`
	SupplierAddress *Address `xml:"supplierAddress"`
	//SupplierBankAccount  string    `xml:"supplierBankAccount,omitempty"` // Not generally used
	//IndividualExemption bool `xml:"individualExemption,omitempty"` // Value is "true" if the seller has individual VAT exempt status
	//ExciseLicenceNum string `xml:"exciseLicenceNum,omitempty"` // Number of supplier’s tax warehouse license or excise license (Act LXVIII of 2016)
}

func newSupplierInfo(supplier *org.Party) (*SupplierInfo, error) {
	taxNumber, groupNumber, err := newTaxNumber(supplier)
	if err != nil {
		return nil, err
	}
	if groupNumber != nil {
		return &SupplierInfo{
			SupplierTaxNumber:    taxNumber,
			GroupMemberTaxNumber: groupNumber,
			SupplierName:         supplier.Name,
			SupplierAddress:      newAddress(supplier.Addresses[0]),
		}, nil
	}
	return &SupplierInfo{
		SupplierTaxNumber: taxNumber,
		SupplierName:      supplier.Name,
		SupplierAddress:   newAddress(supplier.Addresses[0]),
	}, nil

}
