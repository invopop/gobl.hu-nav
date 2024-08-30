package doc

import (
	"reflect"
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func TestNewTaxNumber(t *testing.T) {
	tests := []struct {
		name          string
		party         *org.Party
		expectedMain  *TaxNumber
		expectedGroup *TaxNumber
		expectedErr   error
	}{
		{
			name: "Non-Hungarian TaxID",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.ES.Tax(), // Non-Hungarian country code
					Code:    "1234567890",
				},
			},
			expectedMain: &TaxNumber{
				TaxPayerID: "ES1234567890",
			},
			expectedGroup: nil,
			expectedErr:   nil,
		},
		{
			name: "Hungarian TaxID with VatCode 5 and valid group member code",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.HU.Tax(),
					Code:    "12345678501",
				},
				Identities: []*org.Identity{
					{Code: "12345678402"},
				},
			},
			expectedMain: &TaxNumber{
				TaxPayerID: "12345678",
				VatCode:    "5",
				CountyCode: "01",
			},
			expectedGroup: &TaxNumber{
				TaxPayerID: "12345678",
				VatCode:    "4",
				CountyCode: "02",
			},
			expectedErr: nil,
		},
		{
			name: "Hungarian TaxID with other VatCode",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.HU.Tax(),
					Code:    "12345678301",
				},
			},
			expectedMain: &TaxNumber{
				TaxPayerID: "12345678",
				VatCode:    "3",
				CountyCode: "01",
			},
			expectedGroup: nil,
			expectedErr:   nil,
		},
		{
			name: "Hungarian TaxID with short code",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.HU.Tax(),
					Code:    "12345678",
				},
			},
			expectedMain: &TaxNumber{
				TaxPayerID: "12345678",
			},
			expectedGroup: nil,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mainTaxNum, groupTaxNum := newTaxNumber(tt.party)

			if !reflect.DeepEqual(mainTaxNum, tt.expectedMain) {
				t.Errorf("expected mainTaxNum %v, got %v", tt.expectedMain, mainTaxNum)
			}
			if !reflect.DeepEqual(groupTaxNum, tt.expectedGroup) {
				t.Errorf("expected groupTaxNum %v, got %v", tt.expectedGroup, groupTaxNum)
			}
		})
	}
}
