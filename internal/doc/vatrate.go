package doc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/hu"
	"github.com/invopop/gobl/tax"
)

// VatRate contains the VAT rate information.
// VatRate may contain exactly one of the 8 possible fields
type VatRate struct {
	VatPercentage            string             `xml:"vatPercentage,omitempty"`
	VatContent               string             `xml:"vatContent,omitempty"` //VatContent is only for simplified invoices
	VatExemption             *DetailedReason    `xml:"vatExemption,omitempty"`
	VatOutOfScope            *DetailedReason    `xml:"vatOutOfScope,omitempty"`
	VatDomesticReverseCharge bool               `xml:"vatDomesticReverseCharge,omitempty"`
	MarginSchemeIndicator    string             `xml:"marginSchemeIndicator,omitempty"`
	VatAmountMismatch        *VatAmountMismatch `xml:"vatAmountMismatch,omitempty"`
	NoVatCharge              bool               `xml:"noVatCharge,omitempty"`
}

// DetailedReason contains the case and reason of a VAT exemption or out of scope
type DetailedReason struct {
	Case   string `xml:"case"`
	Reason string `xml:"reason"`
}

// VatAmountMismatch contains the vat rate and case of a vat amount mismatch
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

func newVatRate(obj any, info *taxInfo) (*VatRate, error) {
	switch obj := obj.(type) {
	case *tax.RateTotal:
		return newVatRateTotal(obj, info)
	case *tax.Combo:
		return newVatRateCombo(obj, info)
	}
	return nil, nil
}

func newVatRateTotal(rate *tax.RateTotal, info *taxInfo) (*VatRate, error) {
	// First if it is not exent or simplified invoice we can return the percentage
	if rate.Percent != nil {
		return &VatRate{VatPercentage: rate.Percent.Base().String()}, nil
	}

	// If it is a simplified invoice we can return the content
	if info.simplifiedInvoice {
		return &VatRate{VatContent: rate.Amount.Rescale(2).String()}, nil
	}

	// Check if in the rate extensions there is extkeyexemptioncode or extkeyvatoutofscopecode
	for k, v := range rate.Ext {
		if k == hu.ExtKeyExemptionCode {
			return &VatRate{VatExemption: &DetailedReason{Case: v.String(), Reason: "Exempt"}}, nil
		}

		if k == hu.ExtKeyVatOutOfScopeCode {
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

	//TODO: Missing vat amount mismatch

	// If percent is nil
	if rate.Percent == nil {
		return &VatRate{NoVatCharge: true}, nil
	}

	return nil, ErrNoVatRateField

}

func newVatRateCombo(c *tax.Combo, info *taxInfo) (*VatRate, error) {
	// First if it is not exent or simplified invoice we can return the percentage
	if c.Percent != nil {
		return &VatRate{VatPercentage: c.Percent.Base().String()}, nil
	}

	// Check if in the rate extensions there is extkeyexemptioncode or extkeyvatoutofscopecode
	for k, v := range c.Ext {
		if k == hu.ExtKeyExemptionCode {
			return &VatRate{VatExemption: &DetailedReason{Case: v.String(), Reason: "Exempt"}}, nil
		}

		if k == hu.ExtKeyVatOutOfScopeCode {
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

func newTaxInfo(inv *bill.Invoice) *taxInfo {
	info := &taxInfo{}
	if inv.Tax != nil {
		for _, scheme := range inv.Tax.Tags {
			switch scheme {
			case tax.TagSimplified:
				info.simplifiedInvoice = true
			case hu.TagDomesticReverseCharge:
				info.domesticReverseCharge = true
			case hu.TagTravelAgency:
				info.travelAgency = true
			case hu.TagSecondHand:
				info.secondHand = true
			case hu.TagArt:
				info.art = true
			case hu.TagAntiques:
				info.antique = true
			}
		}
	}
	return info
}
