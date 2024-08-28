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
	XMLName           xml.Name           `xml:"ManageInvoiceRequest"`
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

type ManageInvoiceResponse struct {
	XMLName       xml.Name  `xml:"ManageInvoiceResponse"`
	Header        *Header   `xml:"header"`
	Result        *Result   `xml:"result"`
	Software      *Software `xml:"software"`
	TransactionId string    `xml:"transactionId"`
}

func ReportInvoice(username string, password string, taxNumber string, signKey string, exchangeToken string, soft *Software, invoice string) (string, error) {
	requestData := NewManageInvoiceRequest(username, password, taxNumber, signKey, exchangeToken, soft, invoice)
	return PostManageInvoiceRequest(requestData)
}

func NewManageInvoiceRequest(username string, password string, taxNumber string, signKey string, exchangeToken string, soft *Software, invoice string) ManageInvoiceRequest {
	timestamp := time.Now().UTC()
	requestID := generateRandomString(20) //This must be unique for each request
	operationType := "CREATE"
	return ManageInvoiceRequest{
		Common:        "http://schemas.nav.gov.hu/NTCA/1.0/common",
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
		var manageInvoiceResponse ManageInvoiceResponse
		err = xml.Unmarshal(resp.Body(), &manageInvoiceResponse)
		if err != nil {
			return "", err
		}
		return manageInvoiceResponse.TransactionId, nil
	}

	var generalErrorResponse GeneralErrorResponse
	err = xml.Unmarshal(resp.Body(), &generalErrorResponse)
	if err != nil {
		return "", err
	}

	return "", fmt.Errorf("error code: %s, message: %s", resp.Status(), generalErrorResponse.Result.ErrorCode)
}
