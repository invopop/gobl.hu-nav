package doc

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const (
	vatStatusPrivatePerson = "PRIVATE_PERSON"
	vatStatusDomestic      = "DOMESTIC"
	vatStatusOther         = "OTHER"
)

// CustomerInfo contains the customer data.
type CustomerInfo struct {
	CustomerVatStatus string   `xml:"customerVatStatus"` // PRIVATE_PERSON, DOMESTIC, OTHER
	CustomerVatData   *VatData `xml:"customerVatData,omitempty"`
	CustomerName      string   `xml:"customerName,omitempty"`
	CustomerAddress   *Address `xml:"customerAddress,omitempty"`
	// CustomerBankAccount string   `xml:"customerBankAccountNumber,omitempty"`
}

// VatData contains the VAT subjectivity data of the customer.
type VatData struct {
	CustomerTaxNumber  *CustomerTaxNumber `xml:"customerTaxNumber,omitempty"`
	CommunityVATNumber string             `xml:"communityVATNumber,omitempty"`
	ThirdStateTaxID    string             `xml:"thirdStateTaxId,omitempty"`
}

// CustomerTaxNumber contains the domestic tax number or
// the group identification number, under which the purchase of goods
// or services is done
type CustomerTaxNumber struct {
	TaxPayerID           string     `xml:"base:taxpayerId"`
	VatCode              string     `xml:"base:vatCode,omitempty"`
	CountyCode           string     `xml:"base:countyCode,omitempty"`
	GroupMemberTaxNumber *TaxNumber `xml:"groupMemberTaxNumber,omitempty"`
}

func newCustomerInfo(customer *org.Party) *CustomerInfo {
	taxID := customer.TaxID
	if taxID == nil {
		return &CustomerInfo{
			CustomerVatStatus: vatStatusPrivatePerson,
			CustomerName:      customer.Name,
			CustomerAddress:   newAddress(customer.Addresses[0]),
		}
	}
	status := vatStatusOther

	if taxID.Country == l10n.HU.Tax() {
		if taxID.Code.String() == "" || (taxID.Code.String()[0:1] == "8" && len(taxID.Code) == 10) {
			return &CustomerInfo{
				CustomerVatStatus: vatStatusPrivatePerson,
				CustomerName:      customer.Name,
				CustomerAddress:   newAddress(customer.Addresses[0]),
			}
		}
		status = vatStatusDomestic

	}

	vatData := newVatData(customer, status)

	return &CustomerInfo{
		CustomerVatStatus: status,
		CustomerVatData:   vatData,
		CustomerName:      customer.Name,
		CustomerAddress:   newAddress(customer.Addresses[0]),
	}
}

func newVatData(customer *org.Party, status string) *VatData {
	if status == vatStatusOther {
		return newOtherVatData(customer.TaxID)
	}
	return newDomesticVatData(customer)
}

func newOtherVatData(taxID *tax.Identity) *VatData {
	if isEuropeanCountryCode(taxID.Country.Code()) {
		return &VatData{
			CommunityVATNumber: taxID.String(),
		}
	}
	return &VatData{
		ThirdStateTaxID: taxID.String(),
	}
}

func newDomesticVatData(customer *org.Party) *VatData {
	taxNumber, groupNumber := newTaxNumber(customer)

	if groupNumber != nil {
		return &VatData{
			CustomerTaxNumber: &CustomerTaxNumber{
				TaxPayerID:           taxNumber.TaxPayerID,
				VatCode:              taxNumber.VatCode,
				CountyCode:           taxNumber.CountyCode,
				GroupMemberTaxNumber: groupNumber,
			},
		}
	}

	return &VatData{
		CustomerTaxNumber: &CustomerTaxNumber{
			TaxPayerID: taxNumber.TaxPayerID,
			VatCode:    taxNumber.VatCode,
			CountyCode: taxNumber.CountyCode,
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
