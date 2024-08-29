package doc

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// TaxNumber contains the tax number or group identification number of a party.
type TaxNumber struct {
	TaxPayerID string `xml:"base:taxpayerId"`
	VatCode    string `xml:"base:vatCode,omitempty"`
	CountyCode string `xml:"base:countyCode,omitempty"`
}

func newTaxNumber(party *org.Party) (*TaxNumber, *TaxNumber, error) {
	taxID := party.TaxID
	if taxID.Country == l10n.HU.Tax() {
		if len(taxID.Code) == 11 {
			// If the 9th character of the tax number is 5, then it is a group member tax number
			if taxID.Code.String()[8:9] == "5" {
				groupMemberCode := party.Identities[0].Code.String()
				return newHungarianTaxNumber(taxID.Code.String()),
					newHungarianTaxNumber(groupMemberCode), nil
			}
			return newHungarianTaxNumber(taxID.Code.String()), nil, nil
		}
		// If the tax number is not 11 characters long, then it is 8 characters long
		return &TaxNumber{
			TaxPayerID: taxID.Code.String(),
		}, nil, nil
	}
	// If it is not a Hungarian tax number, then return the tax number with the country code
	return &TaxNumber{
		TaxPayerID: taxID.String(),
	}, nil, nil

}

func newHungarianTaxNumber(code string) *TaxNumber {
	return &TaxNumber{
		TaxPayerID: code[:8],
		VatCode:    code[8:9],
		CountyCode: code[9:11],
	}
}
