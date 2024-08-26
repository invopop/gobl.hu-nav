package doc

import (
	"github.com/invopop/gobl/bill"
	//"github.com/invopop/gobl/regimes/hu"
	"github.com/invopop/gobl/tax"
)

//"github.com/invopop/gobl/regimes/hu"

// Vat Rate may contain exactly one of the 8 possible fields
type VatRate struct {
	VatPercentage            float64            `xml:"vatPercentage,omitempty"`
	VatContent               float64            `xml:"vatContent,omitempty"` //VatContent is only for simplified invoices
	VatExemption             *DetailedReason    `xml:"vatExemption,omitempty"`
	VatOutOfScope            *DetailedReason    `xml:"vatOutOfScope,omitempty"`
	VatDomesticReverseCharge bool               `xml:"vatDomesticReverseCharge,omitempty"`
	MarginSchemeIndicator    string             `xml:"marginSchemeIndicator,omitempty"`
	VatAmountMismatch        *VatAmountMismatch `xml:"vatAmountMismatch,omitempty"`
	NoVatCharge              bool               `xml:"noVatCharge,omitempty"`
}

type DetailedReason struct {
	Case   string `xml:"case"`
	Reason string `xml:"reason"`
}

type VatAmountMismatch struct {
	VatRate float64 `xml:"vatRate"`
	Case    string  `xml:"case"`
}

type taxInfo struct {
	simplifiedInvoice     bool
	domesticReverseCharge bool
	travelAgency          bool
	secondHand            bool
	art                   bool
	antique               bool
}

func NewVatRate(obj any, info *taxInfo) (*VatRate, error) {
	switch obj := obj.(type) {
	case *tax.RateTotal:
		return newVatRateTotal(obj, info)
	case *tax.Combo:
		return newVatRateCombo(obj, info)
	}
	return nil, nil
}

// NewVatRate creates a new VatRate from a taxid
func newVatRateTotal(rate *tax.RateTotal, info *taxInfo) (*VatRate, error) {
	// First if it is not exent or simplified invoice we can return the percentage
	if rate.Percent != nil {
		return &VatRate{VatPercentage: rate.Percent.Amount().Rescale(4).Float64()}, nil
	}

	// If it is a simplified invoice we can return the content
	if info.simplifiedInvoice {
		return &VatRate{VatContent: rate.Amount.Rescale(4).Float64()}, nil
	}

	// Check if in the rate extensions there is extkeyexemptioncode or extkeyvatoutofscopecode
	for k, v := range rate.Ext {
		if k == "hu-exemption-code" { //hu.ExtKeyExemptionCode {
			return &VatRate{VatExemption: &DetailedReason{Case: v.String(), Reason: "Exempt"}}, nil
		}

		if k == "hu-vat-out-of-scope-code" { //hu.ExtKeyVatOutOfScopeCode {
			return &VatRate{VatOutOfScope: &DetailedReason{Case: v.String(), Reason: "Out of Scope"}}, nil
		}
	}

	// Check if it is a domestic reverse charge
	if info.domesticReverseCharge {
		return &VatRate{VatDomesticReverseCharge: true}, nil
	}

	// Check the margin scheme indicators

	if info.travelAgency {
		return &VatRate{MarginSchemeIndicator: "TRAVEL_AGENCY"}, nil
	}
	if info.secondHand {
		return &VatRate{MarginSchemeIndicator: "SECOND_HAND"}, nil
	}
	if info.art {
		return &VatRate{MarginSchemeIndicator: "ARTWORK"}, nil
	}
	if info.antique {
		return &VatRate{MarginSchemeIndicator: "ANTIQUE"}, nil
	}

	// Missing vat amount mismatch

	// If percent is nil
	if rate.Percent == nil {
		return &VatRate{NoVatCharge: true}, nil
	}

	return nil, ErrNoVatRateField

}

func newVatRateCombo(c *tax.Combo, info *taxInfo) (*VatRate, error) {
	// First if it is not exent or simplified invoice we can return the percentage
	if c.Percent != nil {
		return &VatRate{VatPercentage: c.Percent.Amount().Rescale(4).Float64()}, nil
	}

	// Check if in the rate extensions there is extkeyexemptioncode or extkeyvatoutofscopecode
	for k, v := range c.Ext {
		if k == "hu-exemption-code" { //hu.ExtKeyExemptionCode {
			return &VatRate{VatExemption: &DetailedReason{Case: v.String(), Reason: "Exempt"}}, nil
		}

		if k == "hu-vat-out-of-scope-code" { //hu.ExtKeyVatOutOfScopeCode {
			return &VatRate{VatOutOfScope: &DetailedReason{Case: v.String(), Reason: "Out of Scope"}}, nil
		}
	}

	// Check if it is a domestic reverse charge
	if info.domesticReverseCharge {
		return &VatRate{VatDomesticReverseCharge: true}, nil
	}

	// Check the margin scheme indicators
	if info.travelAgency {
		return &VatRate{MarginSchemeIndicator: "TRAVEL_AGENCY"}, nil
	}
	if info.secondHand {
		return &VatRate{MarginSchemeIndicator: "SECOND_HAND"}, nil
	}
	if info.art {
		return &VatRate{MarginSchemeIndicator: "ARTWORK"}, nil
	}
	if info.antique {
		return &VatRate{MarginSchemeIndicator: "ANTIQUE"}, nil
	}

	// Missing vat amount mismatch

	if c.Percent == nil {
		return &VatRate{NoVatCharge: true}, nil
	}

	return nil, ErrNoVatRateField
}

// Until PR approved in regimes this wont work
func newTaxInfo(inv *bill.Invoice) *taxInfo {
	info := &taxInfo{}
	if inv.Tax != nil {
		for _, scheme := range inv.Tax.Tags {
			switch scheme {
			case tax.TagSimplified:
				info.simplifiedInvoice = true
			case "domestic-reverse-charge": //case hu.TagDomesticReverseCharge:
				info.domesticReverseCharge = true
			case "travel-agency": //hu.TagTravelAgency:
				info.travelAgency = true
			case "second-hand": //hu.TagSecondHand:
				info.secondHand = true
			case "art": //hu.TagArt:
				info.art = true
			case "antiques": //hu.TagAntique:
				info.antique = true
			}
		}
	}
	return info
}
