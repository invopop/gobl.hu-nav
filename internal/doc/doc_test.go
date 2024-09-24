package doc

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/invopop/gobl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardInvoice(t *testing.T) {
	// Read the test invoice JSON file
	data, err := os.ReadFile("../../test/data/invoice-standard.json")
	require.NoError(t, err, "Failed to read test invoice file")

	// Unmarshal the JSON into a gobl.Envelope
	env := new(gobl.Envelope)
	err = json.Unmarshal(data, env)
	require.NoError(t, err, "Failed to unmarshal test invoice JSON")

	// Call the NewDocument function
	doc, err := NewDocument(env)
	require.NoError(t, err, "Failed to create new document")

	xmlData, err := doc.toByte()
	require.NoError(t, err, "Failed to marshal document to XML")

	err = saveOutput(xmlData, "invoice-standard.xml")
	require.NoError(t, err, "Failed to save XML output")

	err = schemaValidation(xmlData)
	require.NoError(t, err, "Failed to validate XML output")

	fmt.Println(string(xmlData))

}

func TestCreditNote(t *testing.T) {
	// Read the test invoice JSON file
	data, err := os.ReadFile("../../test/data/credit-note.json")
	require.NoError(t, err, "Failed to read test invoice file")

	// Unmarshal the JSON into a gobl.Envelope
	env := new(gobl.Envelope)
	err = json.Unmarshal(data, env)
	require.NoError(t, err, "Failed to unmarshal test invoice JSON")

	// Call the NewDocument function
	doc, err := NewDocument(env)
	require.NoError(t, err, "Failed to create new document")

	xmlData, err := doc.toByte()
	require.NoError(t, err, "Failed to marshal document to XML")

	err = saveOutput(xmlData, "credit-note.xml")
	require.NoError(t, err, "Failed to save XML output")

	err = schemaValidation(xmlData)
	require.NoError(t, err, "Failed to validate XML output")

	fmt.Println(string(xmlData))

	assert.Equal(t, doc.InvoiceMain.Invoice.InvoiceLines.Lines[0].LineAmountsNormal.LineNetAmountData.LineNetAmount, "-600000.00", "lineNetAmount should be negative")
	assert.Equal(t, doc.InvoiceMain.Invoice.InvoiceSummary.SummaryNormal.SummaryByVatRate[0].VatRateVatData.VatRateVatAmount, "-162000.00", "totalAmount should be negative")

}

func TestB2CInvoice(t *testing.T) {
	// Read the test invoice JSON file
	data, err := os.ReadFile("../../test/data/b2c.json")
	require.NoError(t, err, "Failed to read test invoice file")

	// Unmarshal the JSON into a gobl.Envelope
	env := new(gobl.Envelope)
	err = json.Unmarshal(data, env)
	require.NoError(t, err, "Failed to unmarshal test invoice JSON")

	// Call the NewDocument function
	doc, err := NewDocument(env)
	require.NoError(t, err, "Failed to create new document")

	xmlData, err := doc.toByte()
	require.NoError(t, err, "Failed to marshal document to XML")

	err = saveOutput(xmlData, "b2c.xml")
	require.NoError(t, err, "Failed to save XML output")

	err = schemaValidation(xmlData)
	require.NoError(t, err, "Failed to validate XML output")

	fmt.Println(string(xmlData))

	assert.Equal(t, doc.InvoiceMain.Invoice.InvoiceHead.CustomerInfo.CustomerVatStatus, "PRIVATE_PERSON", "customerVatStatus should be PRIVATE_PERSON")
}

func TestForeignInvoice(t *testing.T) {
	// Read the test invoice JSON file
	data, err := os.ReadFile("../../test/data/foreign.json")
	require.NoError(t, err, "Failed to read test invoice file")

	// Unmarshal the JSON into a gobl.Envelope
	env := new(gobl.Envelope)
	err = json.Unmarshal(data, env)
	require.NoError(t, err, "Failed to unmarshal test invoice JSON")

	// Call the NewDocument function
	doc, err := NewDocument(env)
	require.NoError(t, err, "Failed to create new document")

	xmlData, err := doc.toByte()
	require.NoError(t, err, "Failed to marshal document to XML")

	err = saveOutput(xmlData, "foreign.xml")
	require.NoError(t, err, "Failed to save XML output")

	err = schemaValidation(xmlData)
	require.NoError(t, err, "Failed to validate XML output")

	fmt.Println(string(xmlData))

	assert.Equal(t, doc.InvoiceMain.Invoice.InvoiceLines.Lines[0].LineAmountsNormal.LineNetAmountData.LineNetAmount, "2000.00", "lineNetAmount should be 2000.00")
	assert.NotNil(t, doc.InvoiceMain.Invoice.InvoiceHead.CustomerInfo.CustomerVatData.ThirdStateTaxID, "thirdStateTaxID should not be nil")
}
