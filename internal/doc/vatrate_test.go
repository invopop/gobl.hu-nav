package doc

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNewVatRate(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		info      *taxInfo
		expected  *VatRate
		expectErr bool
	}{
		{
			name: "RateTotal with Percentage",
			input: &tax.RateTotal{
				Percent: num.NewPercentage(27, 4),
			},
			info: &taxInfo{},
			expected: &VatRate{
				VatPercentage: "0.27",
			},
			expectErr: false,
		},
		{
			name: "RateTotal with Simplified Invoice",
			input: &tax.RateTotal{
				Amount: num.MakeAmount(100, 0),
			},
			info: &taxInfo{simplifiedInvoice: true},
			expected: &VatRate{
				VatContent: "100.00",
			},
			expectErr: false,
		},
		{
			name: "RateTotal with Exemption Code",
			input: &tax.RateTotal{
				Ext: tax.Extensions{
					"hu-exemption-code": "AAM",
				},
			},
			info: &taxInfo{},
			expected: &VatRate{
				VatExemption: &DetailedReason{
					Case:   "AAM",
					Reason: "Exempt",
				},
			},
			expectErr: false,
		},
		{
			name: "RateTotal with Out of Scope Code",
			input: &tax.RateTotal{
				Ext: tax.Extensions{
					"hu-vat-out-of-scope-code": "ATK",
				},
			},
			info: &taxInfo{},
			expected: &VatRate{
				VatOutOfScope: &DetailedReason{
					Case:   "ATK",
					Reason: "Out of Scope",
				},
			},
			expectErr: false,
		},
		{
			name:  "RateTotal with Domestic Reverse Charge",
			input: &tax.RateTotal{},
			info:  &taxInfo{domesticReverseCharge: true},
			expected: &VatRate{
				VatDomesticReverseCharge: true,
			},
			expectErr: false,
		},
		{
			name:  "RateTotal with No VAT Charge",
			input: &tax.RateTotal{},
			info:  &taxInfo{},
			expected: &VatRate{
				NoVatCharge: true,
			},
			expectErr: false,
		},
		{
			name:      "Invalid Type Input",
			input:     "invalid",
			info:      &taxInfo{},
			expected:  nil,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vatRate, err := newVatRate(tt.input, tt.info)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, vatRate)
			}
		})
	}
}
