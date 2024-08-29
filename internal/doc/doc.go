package doc

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
)

const (
	XMNLSDATA     = "http://schemas.nav.gov.hu/OSA/3.0/data"
	XMNLSCOMMON   = "http://schemas.nav.gov.hu/NTCA/1.0/common"
	XMNLBASE      = "http://schemas.nav.gov.hu/OSA/3.0/base"
	XMNLXSI       = "http://www.w3.org/2001/XMLSchema-instance"
	XSIDataSchema = "http://schemas.nav.gov.hu/OSA/3.0/data invoiceData.xsd"
)

// Standard error responses.
var (
	ErrNoExchangeRate = newValidationError("no exchange rate to HUF found")
	ErrNoVatRateField = newValidationError("no vat rate field found")
)

// ValidationError is a simple wrapper around validation errors (that should not be retried) as opposed
// to server-side errors (that should be retried).
type ValidationError struct {
	err error
}

// Error implements the error interface for ClientError.
func (e *ValidationError) Error() string {
	return e.err.Error()
}

func newValidationError(text string) error {
	return &ValidationError{errors.New(text)}
}

// Document is the root element of the XML document.
type Document struct {
	XMLName               xml.Name     `xml:"InvoiceData"`
	XMLNS                 string       `xml:"xmlns,attr"`
	XMLNSXsi              string       `xml:"xmlns:xsi,attr"`
	XSISchema             string       `xml:"xsi:schemaLocation,attr"`
	XMLNSCommon           string       `xml:"xmlns:common,attr"`
	XMLNSBase             string       `xml:"xmlns:base,attr"`
	InvoiceNumber         string       `xml:"invoiceNumber"`
	InvoiceIssueDate      string       `xml:"invoiceIssueDate"`
	CompletenessIndicator bool         `xml:"completenessIndicator"` // Indicates whether the data exchange is identical with the invoice (the invoice does not contain any more data)
	InvoiceMain           *InvoiceMain `xml:"invoiceMain"`
}

// NewDocument creates a new Document from an envelope.
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}
	d := new(Document)
	d.XMLNS = XMNLSDATA
	d.XMLNSXsi = XMNLXSI
	d.XSISchema = XSIDataSchema
	d.XMLNSCommon = XMNLSCOMMON
	d.XMLNSBase = XMNLBASE
	d.InvoiceNumber = inv.Code
	d.InvoiceIssueDate = inv.IssueDate.String()
	d.CompletenessIndicator = false
	main, err := newInvoiceMain(inv)
	if err != nil {
		return nil, err
	}
	d.InvoiceMain = main
	return d, nil
}
