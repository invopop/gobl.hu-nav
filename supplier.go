package nav

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
	//IndividualExemption bool `xml:"individualExemption,omitempty"`
	//ExciseLicenceNum string `xml:"exciseLicenceNum,omitempty"`
}

func NewSupplierInfo(supplier *org.Party) (*SupplierInfo, error) {
	taxId := supplier.TaxID
	if taxId.Country != l10n.HU.Tax() {
		return nil, ErrNotHungarian
	}
	if taxId.Code.String()[8:9] != "5" {
		return &SupplierInfo{
			SupplierTaxNumber: NewTaxNumber(taxId),
			SupplierName:      supplier.Name,
			SupplierAddress:   NewAddress(supplier.Addresses[0]),
		}, nil
	}
	groupMemberCode := supplier.Identities[0].Code.String()
	return &SupplierInfo{
		SupplierTaxNumber:    NewTaxNumber(taxId),
		GroupMemberTaxNumber: NewHungarianTaxNumber(groupMemberCode),
		SupplierName:         supplier.Name,
		SupplierAddress:      NewAddress(supplier.Addresses[0]),
	}, nil

}
