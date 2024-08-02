package nav

import "github.com/invopop/gobl/tax"

type TaxNumber struct {
	TaxPayerID string `xml:"base:taxpayerId"`
	//VatCode    string `xml:"base:vatCode,omitempty"`
	//CountyCode string `xml:"base:countyCode,omitempty"`
}

// Have to look at the vatcodes for the regime

// NewTaxNumber creates a new TaxNumber from a taxid
func NewTaxNumber(taxid *tax.Identity) *TaxNumber {
	return &TaxNumber{
		TaxPayerID: taxid.Code.String(),
	}
}
