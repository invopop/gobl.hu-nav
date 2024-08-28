package gateways

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const statusEndpoint = "https://api-test.onlineszamla.nav.gov.hu/invoiceService/v3/queryTransactionStatus"

type QueryTransactionStatusRequest struct {
	XMLName               xml.Name  `xml:"QueryTransactionStatusRequest"`
	Common                string    `xml:"xmlns:common,attr"`
	Xmlns                 string    `xml:"xmlns,attr"`
	Header                *Header   `xml:"common:header"`
	User                  *User     `xml:"common:user"`
	Software              *Software `xml:"software"`
	TransactionId         string    `xml:"transactionId"`
	ReturnOriginalRequest bool      `xml:"returnOriginalRequest,omitempty"`
}

type QueryTransactionStatusResponse struct {
	XMLName           xml.Name           `xml:"QueryTransactionStatusResponse"`
	Header            *Header            `xml:"header"`
	Result            *Result            `xml:"result"`
	Software          *Software          `xml:"software"`
	ProcessingResults *ProcessingResults `xml:"processingResults"`
}

type ProcessingResults struct {
	ProcessingResult       []*ProcessingResult `xml:"processingResult"`
	OriginalRequestVersion string              `xml:"originalRequestVersion"`
	//AnnulmentData          *AnnulmentData      `xml:"annulmentData,omitempty"`
}

type ProcessingResult struct {
	Index                       string                       `xml:"index"`
	BatchIndex                  string                       `xml:"batchIndex,omitempty"`
	InvoiceStatus               string                       `xml:"invoiceStatus"`
	TechnicalValidationMessages *TechnicalValidationMessages `xml:"technicalValidationMessages,omitempty"`
	BusinessValidationMessages  *BusinessValidationMessages  `xml:"businessValidationMessages,omitempty"`
	CompressedContentIndicator  bool                         `xml:"compressedContentIndicator"`
	OriginalRequest             string                       `xml:"originalRequest,omitempty"`
}

type TechnicalValidationMessages struct {
	ValidationResultCode string `xml:"validationResultCode"`
	ValidationErrorCode  string `xml:"validationErrorCode,omitempty"`
	Message              string `xml:"message,omitempty"`
}

type BusinessValidationMessages struct {
	ValidationResultCode string   `xml:"validationResultCode"`
	ValidationErrorCode  string   `xml:"validationErrorCode,omitempty"`
	Message              string   `xml:"message,omitempty"`
	Pointer              *Pointer `xml:"pointer,omitempty"`
}

type Pointer struct {
	Tag                   string `xml:"tag,omitempty"`
	Value                 string `xml:"value,omitempty"`
	Line                  string `xml:"line,omitempty"`
	OriginalInvoiceNumber string `xml:"originalInvoiceNumber,omitempty"`
}

func GetStatus(username string, password string, taxNumber string, signKey string, soft *Software, transactionID string) (*ProcessingResult, error) {
	requestData := NewQueryTransactionStatusRequest(username, password, taxNumber, signKey, soft, transactionID)
	return QueryTransactionStatus(requestData)
}

func NewQueryTransactionStatusRequest(username string, password string, taxNumber string, signKey string, soft *Software, transactionID string) QueryTransactionStatusRequest {
	timestamp := time.Now().UTC()
	requestID := generateRandomString(20)
	return QueryTransactionStatusRequest{
		Xmlns:         "http://schemas.nav.gov.hu/OSA/3.0/api",
		Common:        "http://schemas.nav.gov.hu/NTCA/1.0/common",
		Header:        NewHeader(requestID, timestamp),
		User:          NewUser(username, password, taxNumber, signKey, requestID, timestamp),
		Software:      soft,
		TransactionId: transactionID,
	}
}

func QueryTransactionStatus(requestData QueryTransactionStatusRequest) (*ProcessingResult, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/xml").
		SetBody(requestData).
		Post(statusEndpoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 200 {
		var queryTransactionStatusResponse QueryTransactionStatusResponse
		err = xml.Unmarshal(resp.Body(), &queryTransactionStatusResponse)
		if err != nil {
			return nil, err
		}

		return queryTransactionStatusResponse.ProcessingResults.ProcessingResult[0], nil
	}

	var generalErrorResponse GeneralErrorResponse
	err = xml.Unmarshal(resp.Body(), &generalErrorResponse)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("error code: %s, message: %s", resp.Status(), generalErrorResponse.Result.ErrorCode)
}
