package nav

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/hu"
	"github.com/invopop/gobl/tax"
)

// Vat Rate may contain exactly one of the 8 possible fields
type VatRate struct {
	VatPercentage float64 `xml:"vatPercentage,omitempty"`
	//VatContent    float64 `xml:"vatContent"` //VatContent is only for simplified invoices
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
	//simplifiedRegime bool
	outOfScope            bool
	domesticReverseCharge bool
	travelAgency          bool
	secondHand            bool
	art                   bool
	antique               bool
}

// NewVatRate creates a new VatRate from a taxid
func NewVatRate(rate *tax.RateTotal, info *taxInfo) *VatRate {
	if rate.Key != tax.RateExempt && rate.Key != tax.RateZero {
		return &VatRate{VatPercentage: rate.Percent.Amount().Float64()}
	}
	if rate.Key == tax.RateExempt {
		if info.outOfScope {
			// Q: Is there a way in GOBL to access the extension names?
			// This could maybe be done accessing the regime and there the extensions. We can use the name as the reason.
			return &VatRate{VatOutOfScope: &DetailedReason{Case: rate.Ext[hu.ExtKeyExemptionCode].String(), Reason: "Out of Scope"}}
		}
		if info.domesticReverseCharge {
			return &VatRate{VatDomesticReverseCharge: true}
		}
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
		return &VatRate{VatExemption: &DetailedReason{Case: rate.Ext[hu.ExtKeyExemptionCode].String(), Reason: "Exempt"}}
	}
	return &VatRate{VatPercentage: rate.Percent.Amount().Float64()}

	//TODO: Missing the last 2 cases (VatAmountMismatch and NoVatCharge)
}

// Until PR approved in regimes this wont work
func newTaxInfo(inv *bill.Invoice) *taxInfo {
	info := &taxInfo{}
	if inv.Tax != nil {
		for _, scheme := range inv.Tax.Tags {
			switch scheme {
			case hu.TagOutOfScope:
				info.outOfScope = true
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
