package doc

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/invopop/gobl"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocument(t *testing.T) {
	// Read the test invoice JSON file
	data, err := os.ReadFile("../../test/data/invoice_test.json")
	require.NoError(t, err, "Failed to read test invoice file")

	// Unmarshal the JSON into a gobl.Envelope
	env := new(gobl.Envelope)
	err = json.Unmarshal(data, env)
	require.NoError(t, err, "Failed to unmarshal test invoice JSON")

	// Call the NewDocument function
	doc, err := NewDocument(env)
	require.NoError(t, err, "Failed to create new document")

	xmlData, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling to XML: %v\n", err)
		return
	}

	err = os.WriteFile("../../test/data/out/output.xml", xmlData, 0644)
	if err != nil {
		fmt.Println("Error writing XML to file:", err)
		return
	}

	// 2. Load XSD schema
	xsdContent, err := os.ReadFile("../../test/schemas/invoiceData.xsd")
	if err != nil {
		fmt.Println("Error reading XSD file:", err)
		return
	}

	schema, err := xsd.Parse(xsdContent)
	if err != nil {
		fmt.Println("Error parsing XSD:", err)
		return
	}
	defer schema.Free()

	// 3. Parse XML
	docXML, err := libxml2.ParseString(string(xmlData))
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}
	defer docXML.Free()

	// 4. Validate XML against schema
	if err := schema.Validate(docXML); err != nil {
		fmt.Println("Validation error:", err)
	} else {
		fmt.Println("XML is valid according to the schema")
	}

	fmt.Println(string(xmlData))

	// Assert the expected values
	assert.Equal(t, XMNLSDATA, doc.XMLNS, "Unexpected XMLNS value")
	assert.Equal(t, XMNLXSI, doc.XMLNSXsi, "Unexpected XMLNSXsi value")
	assert.Equal(t, XSIDataSchema, doc.XSISchema, "Unexpected XSISchema value")
	assert.Equal(t, XMNLSCOMMON, doc.XMLNSCommon, "Unexpected XMLNSCommon value")
	assert.Equal(t, XMNLBASE, doc.XMLNSBase, "Unexpected XMLNSBase value")
	assert.False(t, doc.CompletenessIndicator, "Unexpected CompletenessIndicator value")

	// Assert that InvoiceMain is not nil
	assert.NotNil(t, doc.InvoiceMain, "InvoiceMain should not be nil")

}
