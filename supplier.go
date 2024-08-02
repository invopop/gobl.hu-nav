package nav

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

type SupplierInfo struct {
	SupplierTaxNumber *TaxNumber `xml:"supplierTaxNumber"`
	//GroupMemberTaxNumber TaxNumber `xml:"groupMemberTaxNumber,omitempty"`
	//CommunityVATNumber   string    `xml:"communityVATNumber,omitempty"`
	SupplierName    string   `xml:"supplierName"`
	SupplierAddress *Address `xml:"supplierAddress"`
	//SupplierBankAccount  string    `xml:"supplierBankAccount,omitempty"`
	//IndividualExemption bool `xml:"individualExemption,omitempty"`
	//ExciseLicenceNum string `xml:"exciseLicenceNum,omitempty"`
}

func NewSupplierInfo(supplier *org.Party) (*SupplierInfo, error) {
	if supplier.TaxID.Country != l10n.HU.Tax() {
		return nil, ErrNotHungarian
	}
	return &SupplierInfo{
		SupplierTaxNumber: NewTaxNumber(supplier.TaxID),
		SupplierName:      supplier.Name,
		SupplierAddress:   NewAddress(supplier.Addresses[0]),
	}, nil
}
