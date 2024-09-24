package doc

import (
	"reflect"
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func TestNewCustomerInfo(t *testing.T) {
	tests := []struct {
		name            string
		customer        *org.Party
		expectedStatus  string
		expectedVatData *VatData
		expectedName    string
		expectedErr     error
	}{
		{
			name: "Customer with no TaxID",
			customer: &org.Party{
				Name: "John Doe",
				Addresses: []*org.Address{
					{
						Country:  l10n.HU.ISO(),
						Code:     "1234",
						Locality: "Budapest",
						Street:   "Main Street",
					},
				},
			},
			expectedStatus:  "PRIVATE_PERSON",
			expectedVatData: nil,
			expectedName:    "John Doe",
			expectedErr:     nil,
		},
		{
			name: "Hungarian Private Person TaxID",
			customer: &org.Party{
				Name: "Private Person",
				TaxID: &tax.Identity{
					Country: l10n.HU.Tax(),
					Code:    "8123456789", // Indicates a private person
				},
				Addresses: []*org.Address{
					{
						Country:  l10n.HU.ISO(),
						Code:     "1234",
						Locality: "Budapest",
						Street:   "Main Street",
					},
				},
			},
			expectedStatus:  "PRIVATE_PERSON",
			expectedVatData: nil,
			expectedName:    "Private Person",
			expectedErr:     nil,
		},
		{
			name: "Hungarian Domestic VAT Status",
			customer: &org.Party{
				Name: "Hungarian Company",
				TaxID: &tax.Identity{
					Country: l10n.HU.Tax(),
					Code:    "12345678501",
				},
				Identities: []*org.Identity{
					{Code: "12345678402"},
				},
				Addresses: []*org.Address{
					{
						Country:  l10n.HU.ISO(),
						Code:     "1234",
						Locality: "Budapest",
						Street:   "Main Street",
					},
				},
			},
			expectedStatus: "DOMESTIC",
			expectedVatData: &VatData{
				CustomerTaxNumber: &CustomerTaxNumber{
					TaxPayerID: "12345678",
					VatCode:    "5",
					CountyCode: "01",
					GroupMemberTaxNumber: &TaxNumber{
						TaxPayerID: "12345678",
						VatCode:    "4",
						CountyCode: "02",
					},
				},
			},
			expectedName: "Hungarian Company",
			expectedErr:  nil,
		},
		{
			name: "EU Country VAT Status",
			customer: &org.Party{
				Name: "European Company",
				TaxID: &tax.Identity{
					Country: l10n.DE.Tax(),
					Code:    "123456789",
				},
				Addresses: []*org.Address{
					{
						Country:  l10n.HU.ISO(),
						Code:     "1234",
						Locality: "Budapest",
						Street:   "Main Street",
					},
				},
			},
			expectedStatus: "OTHER",
			expectedVatData: &VatData{
				CommunityVATNumber: "DE123456789",
			},
			expectedName: "European Company",
			expectedErr:  nil,
		},
		{
			name: "Non-EU (Third State) VAT Status",
			customer: &org.Party{
				Name: "Non-EU Company",
				TaxID: &tax.Identity{
					Country: l10n.US.Tax(),
					Code:    "123456789",
				},
				Addresses: []*org.Address{
					{
						Country:  l10n.HU.ISO(),
						Code:     "1234",
						Locality: "Budapest",
						Street:   "Main Street",
					},
				},
			},
			expectedStatus: "OTHER",
			expectedVatData: &VatData{
				ThirdStateTaxID: "US123456789",
			},
			expectedName: "Non-EU Company",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerInfo := newCustomerInfo(tt.customer)

			if customerInfo.CustomerVatStatus != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, customerInfo.CustomerVatStatus)
			}
			if !reflect.DeepEqual(customerInfo.CustomerVatData, tt.expectedVatData) {
				t.Errorf("expected VAT data %v, got %v", tt.expectedVatData, customerInfo.CustomerVatData)
			}
			if customerInfo.CustomerName != tt.expectedName {
				t.Errorf("expected name %v, got %v", tt.expectedName, customerInfo.CustomerName)
			}
		})
	}
}
