package nav

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/hu"
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

func NewVatRate(obj any, info *taxInfo) *VatRate {
	switch obj := obj.(type) {
	case *tax.RateTotal:
		return newVatRateTotal(obj, info)
	case *tax.Combo:
		return newVatRateCombo(obj, info)
	}
	return nil
}

// NewVatRate creates a new VatRate from a taxid
func newVatRateTotal(rate *tax.RateTotal, info *taxInfo) *VatRate {
	// First if it is not exent or simplified invoice we can return the percentage
	if rate.Percent != nil {
		return &VatRate{VatPercentage: rate.Percent.Amount().Rescale(4).Float64()}
	}

	// If it is a simplified invoice we can return the content
	if info.simplifiedInvoice {
		return &VatRate{VatContent: rate.Amount.Rescale(4).Float64()}
	}

	// Check if in the rate extensions there is extkeyexemptioncode or extkeyvatoutofscopecode
	for k, v := range rate.Ext {
		if k == hu.ExtKeyExemptionCode {
			return &VatRate{VatExemption: &DetailedReason{Case: v.String(), Reason: "Exempt"}}
		}

		if k == hu.ExtKeyVatOutOfScopeCode {
			return &VatRate{VatOutOfScope: &DetailedReason{Case: v.String(), Reason: "Out of Scope"}}
		}
	}

	// Check if it is a domestic reverse charge
	if info.domesticReverseCharge {
		return &VatRate{VatDomesticReverseCharge: true}
	}

	// Check the margin scheme indicators

	if info.travelAgency {
		return &VatRate{MarginSchemeIndicator: "TRAVEL_AGENCY"}
	}
	if info.secondHand {
		return &VatRate{MarginSchemeIndicator: "SECOND_HAND"}
	}
	if info.art {
		return &VatRate{MarginSchemeIndicator: "ARTWORK"}
	}
	if info.antique {
		return &VatRate{MarginSchemeIndicator: "ANTIQUE"}
	}

	// Missing vat amount mismatch

	return &VatRate{NoVatCharge: true}

}

func newVatRateCombo(c *tax.Combo, info *taxInfo) *VatRate {
	// First if it is not exent or simplified invoice we can return the percentage
	if c.Percent != nil {
		return &VatRate{VatPercentage: c.Percent.Amount().Rescale(4).Float64()}
	}

	// Check if in the rate extensions there is extkeyexemptioncode or extkeyvatoutofscopecode
	for k, v := range c.Ext {
		if k == hu.ExtKeyExemptionCode {
			return &VatRate{VatExemption: &DetailedReason{Case: v.String(), Reason: "Exempt"}}
		}

		if k == hu.ExtKeyVatOutOfScopeCode {
			return &VatRate{VatOutOfScope: &DetailedReason{Case: v.String(), Reason: "Out of Scope"}}
		}
	}

	// Check if it is a domestic reverse charge
	if info.domesticReverseCharge {
		return &VatRate{VatDomesticReverseCharge: true}
	}

	// Check the margin scheme indicators

	if info.travelAgency {
		return &VatRate{MarginSchemeIndicator: "TRAVEL_AGENCY"}
	}
	if info.secondHand {
		return &VatRate{MarginSchemeIndicator: "SECOND_HAND"}
	}
	if info.art {
		return &VatRate{MarginSchemeIndicator: "ARTWORK"}
	}
	if info.antique {
		return &VatRate{MarginSchemeIndicator: "ANTIQUE"}
	}

	// Missing vat amount mismatch

	return &VatRate{NoVatCharge: true}
}

// Until PR approved in regimes this wont work
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
			case hu.TagAntique:
				info.antique = true
			}
		}
	}
	return info
}
