package nav

import (
	"encoding/xml"
	"errors"

	"github.com/invopop/gobl/bill"
)

/* <?xml version="1.0" encoding="UTF-8"?>
<InvoiceData xmlns="http://schemas.nav.gov.hu/OSA/3.0/data" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://schemas.nav.gov.hu/OSA/3.0/data invoiceData.xsd"
xmlns:common="http://schemas.nav.gov.hu/NTCA/1.0/common" xmlns:base="http://schemas.nav.gov.hu/OSA/3.0/base" >*/

// Standard error responses.
var (
	ErrNotHungarian           = newValidationError("only hungarian invoices are supported")
	ErrNoExchangeRate         = newValidationError("no exchange rate to HUF found")
	ErrInvalidGroupMemberCode = newValidationError("invalid group member code")
	ErrNoVatRateField         = newValidationError("no vat rate field found")
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
	XMLName xml.Name `xml:"InvoiceData"`
	XMLNS   string   `xml:"xmlns,attr"`
	//XMLNSXsi              string       `xml:"xmlns:xsi,attr"`
	//XSISchema             string       `xml:"xsi:schemaLocation,attr"`
	//XMLNSCommon           string       `xml:"xmlns:common,attr"`
	//XMLNSBase             string       `xml:"xmlns:base,attr"`
	InvoiceNumber         string       `xml:"invoiceNumber"`
	InvoiceIssueDate      string       `xml:"invoiceIssueDate"`
	CompletenessIndicator bool         `xml:"completenessIndicator"` // Indicates whether the data report is the invoice itself
	InvoiceMain           *InvoiceMain `xml:"invoiceMain"`
}

// Convert it to XML before returning
func NewDocument(inv *bill.Invoice) *Document {
	d := new(Document)
	d.XMLNS = "http://schemas.nav.gov.hu/OSA/3.0/data"
	//d.XMLNSXsi = "http://www.w3.org/2001/XMLSchema-instance"
	//d.XSISchema = "http://schemas.nav.gov.hu/OSA/3.0/data invoiceData.xsd"
	//d.XMLNSCommon = "http://schemas.nav.gov.hu/NTCA/1.0/common"
	//d.XMLNSBase = "http://schemas.nav.gov.hu/OSA/3.0/base"
	d.InvoiceNumber = inv.Code
	d.InvoiceIssueDate = inv.IssueDate.String()
	d.CompletenessIndicator = false
	main, err := NewInvoiceMain(inv)
	if err != nil {
		panic(err)
	}
	d.InvoiceMain = main
	return d
}

/*func main() {
	data, _ := os.ReadFile("invoice-valid.json")
	fmt.Println(string(data))
	env := new(gobl.Envelope)
	if err := json.Unmarshal(data, env); err != nil {
		panic(err)
	}

	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		fmt.Errorf("invalid type %T", env.Document)
	}

	doc := NewDocument(inv)
	// Print the XML
	output, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}*/
