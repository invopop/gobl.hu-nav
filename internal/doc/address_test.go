package doc

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestNewAddress(t *testing.T) {
	tests := []struct {
		name           string
		input          *org.Address
		expectedOutput *Address
	}{
		{
			name: "Detailed address with all fields",
			input: &org.Address{
				Country:    l10n.HU.ISO(),
				Region:     "Budapest",
				Code:       "1234",
				Locality:   "Budapest",
				Street:     "Main",
				StreetType: "street",
				Number:     "10",
				Block:      "B",
				Floor:      "2",
				Door:       "5",
			},
			expectedOutput: &Address{
				DetailedAddress: &DetailedAddress{
					CountryCode:         "HU",
					Region:              "Budapest",
					PostalCode:          "1234",
					City:                "Budapest",
					StreetName:          "Main",
					Number:              "10",
					Building:            "B",
					Floor:               "2",
					Door:                "5",
					PublicPlaceCategory: "street",
				},
			},
		},
		{
			name: "Detailed address with missing optional fields",
			input: &org.Address{
				Country:  l10n.HU.ISO(),
				Code:     "1234",
				Locality: "Budapest",
				Street:   "Main Street",
				Number:   "10",
			},
			expectedOutput: &Address{
				SimpleAddress: &SimpleAddress{
					CountryCode:             "HU",
					PostalCode:              "1234",
					City:                    "Budapest",
					AdditionalAddressDetail: "Main Street, 10",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newAddress(tt.input)
			assert.Equal(t, tt.expectedOutput, result)
		})
	}
}
