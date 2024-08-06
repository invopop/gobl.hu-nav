package nav

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
)

type TaxNumber struct {
	TaxPayerID string `xml:"base:taxpayerId"`
	VatCode    string `xml:"base:vatCode,omitempty"`
	CountyCode string `xml:"base:countyCode,omitempty"`
}

// Have to look at the vatcodes for the regime

// NewTaxNumber creates a new TaxNumber from a taxid
func NewTaxNumber(taxid *tax.Identity) *TaxNumber {
	if taxid.Country == l10n.HU.Tax() {
		// Validate here or in validation: Only valid vat codes are 1,2,3 and 5 for the tax id (for the group could be 4)
		return NewHungarianTaxNumber(taxid.Code.String())
	} else {
		return &TaxNumber{
			TaxPayerID: taxid.String(),
		}
	}
}

func NewHungarianTaxNumber(code string) *TaxNumber {
	return &TaxNumber{
		TaxPayerID: code[:8],
		VatCode:    code[8:9],
		CountyCode: code[9:11],
	}
}
