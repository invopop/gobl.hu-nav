package doc

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/invopop/gobl/bill"
)

/* <?xml version="1.0" encoding="UTF-8"?>
<InvoiceData xmlns="http://schemas.nav.gov.hu/OSA/3.0/data" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://schemas.nav.gov.hu/OSA/3.0/data invoiceData.xsd"
xmlns:common="http://schemas.nav.gov.hu/NTCA/1.0/common" xmlns:base="http://schemas.nav.gov.hu/OSA/3.0/base" >*/

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

type Document struct {
	XMLName               xml.Name     `xml:"InvoiceData"`
	XMLNS                 string       `xml:"xmlns,attr"`
	XMLNSXsi              string       `xml:"xmlns:xsi,attr"`
	XSISchema             string       `xml:"xsi:schemaLocation,attr"`
	XMLNSCommon           string       `xml:"xmlns:common,attr"`
	XMLNSBase             string       `xml:"xmlns:base,attr"`
	InvoiceNumber         string       `xml:"invoiceNumber"`
	InvoiceIssueDate      string       `xml:"invoiceIssueDate"`
	CompletenessIndicator bool         `xml:"completenessIndicator"` // Indicates whether the data report is the invoice itself
	InvoiceMain           *InvoiceMain `xml:"invoiceMain"`
}

// Convert it to XML before returning
func NewDocument(inv *bill.Invoice) (*Document, error) {
	d := new(Document)
	d.XMLNS = XMNLSDATA
	d.XMLNSXsi = XMNLXSI
	d.XSISchema = XSIDataSchema
	d.XMLNSCommon = XMNLSCOMMON
	d.XMLNSBase = XMNLBASE
	d.InvoiceNumber = inv.Code
	d.InvoiceIssueDate = inv.IssueDate.String()
	d.CompletenessIndicator = false
	main, err := NewInvoiceMain(inv)
	if err != nil {
		return nil, err
	}
	d.InvoiceMain = main
	return d, nil
}

// BytesIndent returns the indented XML document bytes
func (doc *Document) BytesIndent() ([]byte, error) {
	return toBytesIndent(doc)
}

func toBytesIndent(doc any) ([]byte, error) {
	buf, err := buffer(doc, xml.Header, true)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buffer(doc any, base string, indent bool) (*bytes.Buffer, error) {
	buf := bytes.NewBufferString(base)

	enc := xml.NewEncoder(buf)
	if indent {
		enc.Indent("", "  ")
	}

	if err := enc.Encode(doc); err != nil {
		return nil, fmt.Errorf("encoding document: %w", err)
	}

	return buf, nil
}
