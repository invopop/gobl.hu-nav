package doc

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

type TaxNumber struct {
	TaxPayerID string `xml:"base:taxpayerId"`
	VatCode    string `xml:"base:vatCode,omitempty"`
	CountyCode string `xml:"base:countyCode,omitempty"`
}

// Have to look at the vatcodes for the regime

// NewTaxNumber creates a new TaxNumber from a taxid
func NewTaxNumber(party *org.Party) (*TaxNumber, *TaxNumber, error) {
	taxID := party.TaxID
	if taxID.Country == l10n.HU.Tax() {
		// If the vat code is 5, then the group member code should be 4
		if len(taxID.Code) == 11 {
			if taxID.Code.String()[8:9] == "5" {
				groupMemberCode := party.Identities[0].Code.String()
				if len(groupMemberCode) != 11 || groupMemberCode[8:9] != "4" {
					return nil, nil, ErrInvalidGroupMemberCode
				}
				return NewHungarianTaxNumber(taxID.Code.String()),
					NewHungarianTaxNumber(groupMemberCode), nil
			}
			return NewHungarianTaxNumber(taxID.Code.String()), nil, nil
		}
		return &TaxNumber{
			TaxPayerID: taxID.Code.String(),
		}, nil, nil
	}
	return &TaxNumber{
		TaxPayerID: taxID.String(),
	}, nil, nil

}

func NewHungarianTaxNumber(code string) *TaxNumber {
	return &TaxNumber{
		TaxPayerID: code[:8],
		VatCode:    code[8:9],
		CountyCode: code[9:11],
	}
}
