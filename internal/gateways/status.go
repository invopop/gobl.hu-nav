package gateways

import (
	"encoding/xml"
	"fmt"
	"time"
)

// QueryTransactionStatusRequest contains the needed data to query the status of a transaction
type QueryTransactionStatusRequest struct {
	XMLName               xml.Name     `xml:"QueryTransactionStatusRequest"`
	Common                string       `xml:"xmlns:common,attr"`
	Xmlns                 string       `xml:"xmlns,attr"`
	Header                *Header      `xml:"common:header"`
	User                  *UserRequest `xml:"common:user"`
	Software              *Software    `xml:"software"`
	TransactionID         string       `xml:"transactionId"`
	ReturnOriginalRequest bool         `xml:"returnOriginalRequest,omitempty"`
}

// QueryTransactionStatusResponse contains the response from the NAV API after querying the status of a transaction
type QueryTransactionStatusResponse struct {
	XMLName           xml.Name           `xml:"QueryTransactionStatusResponse"`
	Header            *Header            `xml:"header"`
	Result            *Result            `xml:"result"`
	Software          *Software          `xml:"software"`
	ProcessingResults *ProcessingResults `xml:"processingResults"`
}

// ProcessingResults contains the results of a transaction
// It contains a list of ProcessingResult, which contains the status of each invoice in the transaction
type ProcessingResults struct {
	ProcessingResult       []*ProcessingResult `xml:"processingResult"`
	OriginalRequestVersion string              `xml:"originalRequestVersion"`
	//AnnulmentData          *AnnulmentData      `xml:"annulmentData,omitempty"`
}

// ProcessingResult contains the status of an invoice in a transaction
// It also contains the messages from the technical and business validations
type ProcessingResult struct {
	Index                       string                       `xml:"index"`
	BatchIndex                  string                       `xml:"batchIndex,omitempty"`
	InvoiceStatus               string                       `xml:"invoiceStatus"`
	TechnicalValidationMessages *TechnicalValidationMessages `xml:"technicalValidationMessages,omitempty"`
	BusinessValidationMessages  *BusinessValidationMessages  `xml:"businessValidationMessages,omitempty"`
	CompressedContentIndicator  bool                         `xml:"compressedContentIndicator"`
	OriginalRequest             string                       `xml:"originalRequest,omitempty"`
}

// TechnicalValidationMessages are the result of the technical validation
type TechnicalValidationMessages struct {
	ValidationResultCode string `xml:"validationResultCode"`
	ValidationErrorCode  string `xml:"validationErrorCode,omitempty"`
	Message              string `xml:"message,omitempty"`
}

// BusinessValidationMessages are the result of the business validation
type BusinessValidationMessages struct {
	ValidationResultCode string   `xml:"validationResultCode"`
	ValidationErrorCode  string   `xml:"validationErrorCode,omitempty"`
	Message              string   `xml:"message,omitempty"`
	Pointer              *Pointer `xml:"pointer,omitempty"`
}

// Pointer points to the specific part of the invoice that caused the validation to fail
type Pointer struct {
	Tag                   string `xml:"tag,omitempty"`
	Value                 string `xml:"value,omitempty"`
	Line                  string `xml:"line,omitempty"`
	OriginalInvoiceNumber string `xml:"originalInvoiceNumber,omitempty"`
}

// GetStatus queries the status of a transaction
func (g *Client) GetStatus(transactionID string) ([]*ProcessingResult, error) {
	requestData := g.newQueryTransactionStatusRequest(transactionID)
	return g.queryTransactionStatus(requestData)
}

func (g *Client) newQueryTransactionStatusRequest(transactionID string) QueryTransactionStatusRequest {
	timestamp := time.Now().UTC()
	requestID := newRequestID(timestamp)
	return QueryTransactionStatusRequest{
		Xmlns:         "http://schemas.nav.gov.hu/OSA/3.0/api",
		Common:        "http://schemas.nav.gov.hu/NTCA/1.0/common",
		Header:        NewHeader(requestID, timestamp),
		User:          g.NewUser(requestID, timestamp),
		Software:      g.software,
		TransactionID: transactionID,
	}
}

func (g *Client) queryTransactionStatus(requestData QueryTransactionStatusRequest) ([]*ProcessingResult, error) {
	resp, err := g.rest.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/xml").
		SetBody(requestData).
		Post(StatusEndpoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 200 {
		var queryTransactionStatusResponse QueryTransactionStatusResponse
		err = xml.Unmarshal(resp.Body(), &queryTransactionStatusResponse)
		if err != nil {
			return nil, err
		}

		return queryTransactionStatusResponse.ProcessingResults.ProcessingResult, nil
	}

	var generalErrorResponse GeneralErrorResponse
	err = xml.Unmarshal(resp.Body(), &generalErrorResponse)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("error code: %s, message: %s", resp.Status(), generalErrorResponse.Result.ErrorCode)
}
