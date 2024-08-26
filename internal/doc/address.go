package doc

import (
	"fmt"

	"github.com/invopop/gobl/org"
)

type Address struct {
	SimpleAddress   *SimpleAddress   `xml:"base:simpleAddress,omitempty"`
	DetailedAddress *DetailedAddress `xml:"base:detailedAddress,omitempty"`
}

// DetailedAddressType represents detailed address data
type DetailedAddress struct {
	CountryCode         string `xml:"base:countryCode"`
	Region              string `xml:"base:region,omitempty"`
	PostalCode          string `xml:"base:postalCode"`
	City                string `xml:"base:city"`
	StreetName          string `xml:"base:streetName"`
	PublicPlaceCategory string `xml:"base:publicPlaceCategory"`
	Number              string `xml:"base:number,omitempty"`
	Building            string `xml:"base:building,omitempty"`
	Staircase           string `xml:"base:staircase,omitempty"`
	Floor               string `xml:"base:floor,omitempty"`
	Door                string `xml:"base:door,omitempty"`
	LotNumber           string `xml:"base:lotNumber,omitempty"`
}

// GOBL does not support dividing the address into public place category and street name
// For the moment we can use SimpleAddress

// SimpleAddressType represents a simple address
type SimpleAddress struct {
	CountryCode             string `xml:"countryCode"`
	Region                  string `xml:"region,omitempty"`
	PostalCode              string `xml:"base:postalCode"`
	City                    string `xml:"base:city"`
	AdditionalAddressDetail string `xml:"base:additionalAddressDetail"`
}

func NewAddress(address *org.Address) *Address {
	return &Address{
		DetailedAddress: NewDetailedAddress(address),
	}
	/*return &Address{
		SimpleAddress: &SimpleAddress{
			CountryCode:             address.Country.String(),
			PostalCode:              address.Code,
			City:                    address.Locality,
			AdditionalAddressDetail: formatAddress(address),
		},
	}*/
}

func NewDetailedAddress(address *org.Address) *DetailedAddress {
	return &DetailedAddress{
		CountryCode:         address.Country.String(),
		Region:              address.Region,
		PostalCode:          address.Code,
		City:                address.Locality,
		StreetName:          address.Street,
		Number:              address.Number,
		Building:            address.Block,
		Floor:               address.Floor,
		Door:                address.Door,
		PublicPlaceCategory: "utca", //address.StreetType, //Waiting for PR to be approved
	}
}

// This is used only for SimpleAddress
func formatAddress(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return "PO Box / Apdo " + address.PostOfficeBox
	}

	formattedAddress := fmt.Sprintf("%s, %s", address.Street, address.Number)

	if address.Block != "" {
		formattedAddress += (", " + address.Block)
	}

	if address.Floor != "" {
		formattedAddress += (", " + address.Floor)
	}

	if address.Door != "" {
		formattedAddress += (" " + address.Door)
	}

	if address.StreetExtra != "" {
		formattedAddress += ("\n" + address.StreetExtra)
	}

	return formattedAddress
}
