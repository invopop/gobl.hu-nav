package doc

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

type SupplierInfo struct {
	SupplierTaxNumber    *TaxNumber `xml:"supplierTaxNumber"`
	GroupMemberTaxNumber *TaxNumber `xml:"groupMemberTaxNumber,omitempty"`
	//CommunityVATNumber   string    `xml:"communityVATNumber,omitempty"` // This is just the same as Supplier Number is HU
	SupplierName    string   `xml:"supplierName"`
	SupplierAddress *Address `xml:"supplierAddress"`
	//SupplierBankAccount  string    `xml:"supplierBankAccount,omitempty"` // Not generally used
	//IndividualExemption bool `xml:"individualExemption,omitempty"` // Value is "true" if the seller has individual VAT exempt status
	//ExciseLicenceNum string `xml:"exciseLicenceNum,omitempty"`
}

func NewSupplierInfo(supplier *org.Party) (*SupplierInfo, error) {
	taxId := supplier.TaxID
	if taxId.Country != l10n.HU.Tax() {
		return nil, ErrNotHungarian
	}
	taxNumber, groupNumber, err := NewTaxNumber(supplier)
	if err != nil {
		return nil, err
	}
	if groupNumber != nil {
		return &SupplierInfo{
			SupplierTaxNumber:    taxNumber,
			GroupMemberTaxNumber: groupNumber,
			SupplierName:         supplier.Name,
			SupplierAddress:      NewAddress(supplier.Addresses[0]),
		}, nil
	}
	return &SupplierInfo{
		SupplierTaxNumber: taxNumber,
		SupplierName:      supplier.Name,
		SupplierAddress:   NewAddress(supplier.Addresses[0]),
	}, nil

}
