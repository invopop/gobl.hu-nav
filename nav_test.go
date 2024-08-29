package nav

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestReportInvoice(t *testing.T) {

	software := NewSoftware(
		tax.Identity{Country: l10n.ES.Tax(), Code: cbc.Code("B12345678")},
		"Invopop",
		"ONLINE_SERVICE",
		"1.0.0",
		"TestDev",
		"pablo.menendez@invopop.com",
	)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	userID := os.Getenv("USER_ID")
	userPWD := os.Getenv("USER_PWD")
	signKey := os.Getenv("SIGN_KEY")
	exchangeKey := os.Getenv("EXCHANGE_KEY")
	taxID := os.Getenv("TAX_ID")

	user := NewUser(userID, userPWD, signKey, exchangeKey, taxID)

	navClient := NewNav(user, software, InTesting())

	xmlContent, err := os.ReadFile("test/data/out/output.xml")
	if err != nil {
		t.Fatalf("Failed to read sample invoice file: %v", err)
	}
	encodedInvoice := base64.StdEncoding.EncodeToString(xmlContent)

	navClient.FetchToken()

	transactionId, err := navClient.ReportInvoice(encodedInvoice)

	fmt.Println("Transaction ID: ", transactionId)

	// Assert the result
	if err != nil {
		t.Errorf("ReportInvoice returned an unexpected error: %v", err)
	}
	require.NoError(t, err, "Expected no error")

	resultsList, err := navClient.GetTransactionStatus(transactionId)
	if err != nil {
		t.Errorf("GetTransactionStatus returned an unexpected error: %v", err)
	}
	require.NoError(t, err, "Expected no error")

	// Print result in xml format for debugging
	xmlData, err := xml.MarshalIndent(resultsList, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling to XML: %v\n", err)
		return
	}

	fmt.Println(string(xmlData))
}
