package gateways

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"
)

// ManageInvoiceRequest represents the root element
type ManageInvoiceRequest struct {
	XMLName           xml.Name           `xml:"ManageInvoiceRequest"`
	Common            string             `xml:"xmlns:common,attr"`
	Xmlns             string             `xml:"xmlns,attr"`
	Header            *Header            `xml:"common:header"`
	User              *UserRequest       `xml:"common:user"`
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

func (g *Client) ReportInvoice(invoice []byte) (string, error) {
	// We first fetch the exchange token
	g.GetToken()
	encodedInvoice := base64.StdEncoding.EncodeToString(invoice)
	requestData := g.newManageInvoiceRequest(encodedInvoice)
	return g.postManageInvoiceRequest(requestData)
}

func (g *Client) newManageInvoiceRequest(invoice string) ManageInvoiceRequest {
	timestamp := time.Now().UTC()
	requestID := NewRequestID(timestamp)
	operationType := "CREATE" // For the moment, only CREATE is supported
	return ManageInvoiceRequest{
		Common:        "http://schemas.nav.gov.hu/NTCA/1.0/common",
		Xmlns:         "http://schemas.nav.gov.hu/OSA/3.0/api",
		Header:        NewHeader(requestID, timestamp),
		User:          g.NewUser(requestID, timestamp, operationType+invoice),
		Software:      g.software,
		ExchangeToken: g.token.Token,
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

func (g *Client) postManageInvoiceRequest(requestData ManageInvoiceRequest) (string, error) {
	resp, err := g.rest.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/xml").
		SetBody(requestData).
		Post(ManageInvoiceEndpoint)

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
