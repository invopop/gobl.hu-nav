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
	exchangeKey := os.Getenv("EXCHANGE_KEY")
	taxID := os.Getenv("TAX_ID")

	user := NewUser(userID, userPWD, signKey, exchangeKey, taxID)

	client := New(user, software, Environment("testing"))

	result, err := client.GetStatus("4OYE2J5GEOGWKMYV")

	// Assert the results
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Print result in xml format for debugging
	xmlData, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling to XML: %v\n", err)
		return
	}

	fmt.Println(string(xmlData))

	for _, r := range result {
		assert.Equal(t, "1", r.Index)
		assert.Equal(t, "DONE", r.InvoiceStatus)
		assert.False(t, r.CompressedContentIndicator)
	}

}
