package nav

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

type CustomerInfo struct {
	CustomerVatStatus string   `xml:"customerVatStatus"`
	CustomerVatData   *VatData `xml:"customerVatData,omitempty"`
	CustomerName      string   `xml:"customerName,omitempty"`
	CustomerAddress   *Address `xml:"customerAddress,omitempty"`
	// CustomerBankAccount string   `xml:"customerBankAccountNumber,omitempty"`
}

type VatData struct {
	CustomerTaxNumber  *CustomerTaxNumber `xml:"customerTaxNumber,omitempty"`
	CommunityVATNumber string             `xml:"communityVATNumber,omitempty"`
	ThirdStateTaxId    string             `xml:"thirdStateTaxId,omitempty"`
}

type CustomerTaxNumber struct {
	TaxPayerID           string     `xml:"taxpayerId"`
	VatCode              string     `xml:"vatCode,omitempty"`
	CountyCode           string     `xml:"countyCode,omitempty"`
	GroupMemberTaxNumber *TaxNumber `xml:"groupMemberTaxNumber,omitempty"`
}

func NewCustomerInfo(inv *bill.Invoice) *CustomerInfo {

	customer := inv.Customer

	if inv.Tax.ContainsTag(hu.TagPrivatePerson) {
		return &CustomerInfo{
			CustomerVatStatus: "PRIVATE_PERSON",
			CustomerName:      customer.Name,
			CustomerAddress:   NewAddress(customer.Addresses[0]),
		}
	}

	// If the customer is not a taxable person
	if customer.TaxID == nil {
		return &CustomerInfo{
			CustomerVatStatus: "OTHER",
			CustomerName:      customer.Name,
			CustomerAddress:   NewAddress(customer.Addresses[0]),
		}
	}

	taxID := customer.TaxID
	group := false
	status := "OTHER"

	if taxID.Country == l10n.HU.Tax() {
		status = "DOMESTIC"
		// One case for group Id and other for simple (Group ID has the 9th character as 5)
		if taxID.Code.String()[8:9] == "5" {
			group = true
		}
	}
	return &CustomerInfo{
		CustomerVatStatus: status,
		CustomerVatData:   newVatData(customer, group, status),
		CustomerName:      customer.Name,
		CustomerAddress:   NewAddress(customer.Addresses[0]),
	}
}

func newVatData(customer *org.Party, group bool, status string) *VatData {
	if status == "OTHER" {
		return newOtherVatData(customer.TaxID)
	}
	return newDomesticVatData(customer, group)
}

func newOtherVatData(taxID *tax.Identity) *VatData {
	if isEuropeanCountryCode(taxID.Country.Code()) {
		return &VatData{
			CommunityVATNumber: taxID.String(),
		}
	}
	return &VatData{
		ThirdStateTaxId: taxID.String(),
	}
}

func newDomesticVatData(customer *org.Party, group bool) *VatData {
	taxID := customer.TaxID
	if group {
		groupMemberCode := customer.Identities[0].Code.String()
		return &VatData{
			CustomerTaxNumber: &CustomerTaxNumber{
				TaxPayerID:           taxID.Code.String()[0:8],
				VatCode:              taxID.Code.String()[8:9],
				CountyCode:           taxID.Code.String()[9:11],
				GroupMemberTaxNumber: NewHungarianTaxNumber(groupMemberCode),
			},
		}
	}
	return &VatData{
		CustomerTaxNumber: &CustomerTaxNumber{
			TaxPayerID: taxID.Code.String()[0:8],
			VatCode:    taxID.Code.String()[8:9],
			CountyCode: taxID.Code.String()[9:11],
		},
	}
}

var europeanCountryCodes = []l10n.Code{
	l10n.AT, l10n.BE, l10n.BG, l10n.CY, l10n.CZ, l10n.DE, l10n.DK, l10n.EE, l10n.EL, l10n.ES,
	l10n.FI, l10n.FR, l10n.HR, l10n.HU, l10n.IE, l10n.IT, l10n.LT, l10n.LU, l10n.LV, l10n.MT,
	l10n.NL, l10n.PL, l10n.PT, l10n.RO, l10n.SE, l10n.SI, l10n.SK, l10n.XI,
}

func isEuropeanCountryCode(code l10n.Code) bool {
	for _, c := range europeanCountryCodes {
		if c == code {
			return true
		}
	}
	return false
}
