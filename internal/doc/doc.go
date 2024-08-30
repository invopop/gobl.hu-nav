// Package doc includes the conversion from GOBL to Nav XML format.
package doc

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
)

// XML schemas
const (
	XMNLSDATA     = "http://schemas.nav.gov.hu/OSA/3.0/data"
	XMNLSCOMMON   = "http://schemas.nav.gov.hu/NTCA/1.0/common"
	XMNLBASE      = "http://schemas.nav.gov.hu/OSA/3.0/base"
	XMNLXSI       = "http://www.w3.org/2001/XMLSchema-instance"
	XSIDataSchema = "http://schemas.nav.gov.hu/OSA/3.0/data invoiceData.xsd"
)

const baseDirectory = "../../test/data/out/"

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

	// Invert if we're dealing with a credit note
	if inv.Type == bill.InvoiceTypeCreditNote {
		if err := inv.Invert(); err != nil {
			return nil, fmt.Errorf("inverting invoice: %w", err)
		}
	}

	d := &Document{
		XMLNS:                 XMNLSDATA,
		XMLNSXsi:              XMNLXSI,
		XSISchema:             XSIDataSchema,
		XMLNSCommon:           XMNLSCOMMON,
		XMLNSBase:             XMNLBASE,
		InvoiceNumber:         inv.Code,
		InvoiceIssueDate:      inv.IssueDate.String(),
		CompletenessIndicator: false,
	}
	main, err := newInvoiceMain(inv)
	if err != nil {
		return nil, err
	}
	d.InvoiceMain = main
	return d, nil
}

func schemaValidation(xmlData []byte) error {

	schema, err := loadSchema()
	if err != nil {
		return fmt.Errorf("error loading schema: %w", err)
	}

	docXML, err := libxml2.ParseString(string(xmlData))
	if err != nil {
		return fmt.Errorf("error parsing XML: %w", err)
	}
	defer docXML.Free()

	if err := schema.Validate(docXML); err != nil {
		validationErrors := err.(xsd.SchemaValidationError)
		fmt.Println("XML Validation Errors:")
		for _, verr := range validationErrors.Errors() {
			fmt.Printf("- %s\n", verr.Error())
		}
		return fmt.Errorf("XML validation failed: %w", err)
	}

	return nil
}

func (d *Document) toByte() ([]byte, error) {
	xmlData, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("Error marshalling to XML: %w", err)
	}

	return xmlData, nil
}

func saveOutput(xmlData []byte, fileName string) error {
	err := os.WriteFile(baseDirectory+fileName, xmlData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing XML to file: %w", err)
	}

	return nil
}

func loadSchema() (*xsd.Schema, error) {
	xsdContent, err := os.ReadFile("../../test/schemas/invoiceData.xsd")
	if err != nil {
		return nil, fmt.Errorf("Error reading XSD file: %w", err)
	}

	schema, err := xsd.Parse(xsdContent)
	if err != nil {
		return nil, fmt.Errorf("Error parsing XSD: %w", err)
	}

	return schema, nil
}
