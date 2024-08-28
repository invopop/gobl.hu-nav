package gateways

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestQueryTransactionStatus(t *testing.T) {
	// Set up test data
	software := NewSoftware(
		tax.Identity{Country: l10n.ES.Tax(), Code: cbc.Code("B12345678")},
		"Invopop",
		"ONLINE_SERVICE",
		"1.0.0",
		"TestDev",
		"pablo.menendez@invopop.com",
	)

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	userID := os.Getenv("USER_ID")
	userPWD := os.Getenv("USER_PWD")
	signKey := os.Getenv("SIGN_KEY")
	taxID := os.Getenv("TAX_ID")

	requestData := NewQueryTransactionStatusRequest(userID, userPWD, taxID, signKey, software, "4OYE2J5GEOGWKMYV")

	result, err := QueryTransactionStatus(requestData)

	if err != nil {
		fmt.Printf("Error querying transaction status: %v\n", err)
		return
	}

	// Print result in xml format for debugging
	xmlData, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling to XML: %v\n", err)
		return
	}

	fmt.Println(string(xmlData))

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Index)
	assert.Equal(t, "DONE", result.InvoiceStatus)
	assert.False(t, result.CompressedContentIndicator)
}
