package gateways

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	manageInvoiceEndpoint = "https://api-test.onlineszamla.nav.gov.hu/invoiceService/v3/manageInvoice"
)

// ManageInvoiceRequest represents the root element
type ManageInvoiceRequest struct {
	XMLName           xml.Name           `xml:"http://schemas.nav.gov.hu/OSA/3.0/api ManageInvoiceRequest"`
	Common            string             `xml:"xmlns:common,attr"`
	Xmlns             string             `xml:"xmlns,attr"`
	Header            *Header            `xml:"common:header"`
	User              *User              `xml:"common:user"`
	Software          *Software          `xml:"software"`
	ExchangeToken     string             `xml:"exchangeToken"`
	InvoiceOperations *InvoiceOperations `xml:"invoiceOperations"`
}

// InvoiceOperations represents the invoiceOperations element
type InvoiceOperations struct {
	CompressedContent bool                `xml:"compressedContent"`
	InvoiceOperation  []*InvoiceOperation `xml:"invoiceOperation"`
}

// InvoiceOperation represents the invoiceOperation element
type InvoiceOperation struct {
	Index                 int                    `xml:"index"`
	InvoiceOperationType  string                 `xml:"invoiceOperation"`
	InvoiceData           string                 `xml:"invoiceData"`
	ElectronicInvoiceHash *ElectronicInvoiceHash `xml:"electronicInvoiceHash,omitempty"`
}

// ElectronicInvoiceHash represents the electronicInvoiceHash element
type ElectronicInvoiceHash struct {
	CryptoType string `xml:"cryptoType,attr"`
	Value      string `xml:",chardata"`
}

func NewManageInvoiceRequest(username string, password string, taxNumber string, signKey string, exchangeToken string, soft *Software, invoice string) ManageInvoiceRequest {
	timestamp := time.Now().UTC()
	requestID := generateRandomString(20) //This cannnot be repeated in the time window of 5 mins
	operationType := "CREATE"
	return ManageInvoiceRequest{
		Common:        "http://schemas.nav.gov.hu/OSA/3.0/common",
		Xmlns:         "http://schemas.nav.gov.hu/OSA/3.0/api",
		Header:        NewHeader(requestID, timestamp),
		User:          NewUser(username, password, taxNumber, signKey, requestID, timestamp, operationType+invoice),
		Software:      soft,
		ExchangeToken: exchangeToken,
		InvoiceOperations: &InvoiceOperations{
			CompressedContent: false,
			InvoiceOperation: []*InvoiceOperation{
				{
					Index:                1,
					InvoiceOperationType: operationType,
					InvoiceData:          invoice,
				},
			},
		},
	}
}

func PostManageInvoiceRequest(requestData ManageInvoiceRequest) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/xml").
		SetBody(requestData).
		Post(manageInvoiceEndpoint)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == 200 {
		return resp.String(), nil
	}

	return "", fmt.Errorf("error code: %s", resp.Status())
}
