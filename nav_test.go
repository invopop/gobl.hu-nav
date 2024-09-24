package nav

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
	"github.com/stretchr/testify/require"
)

func TestReportInvoice(t *testing.T) {

	software := NewSoftware(
		tax.Identity{Country: l10n.ES.Tax(), Code: cbc.Code("B12345678")},
		"Invopop",
		"ONLINE_SERVICE",
		"1.0.0",
		"TestDev",
		"test@dev.com",
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

	invoice, err := os.ReadFile("test/data/out/invoice-standard.xml")
	if err != nil {
		t.Fatalf("Failed to read sample invoice file: %v", err)
	}

	transactionID, err := navClient.ReportInvoice(invoice, "CREATE")

	fmt.Println("Transaction ID: ", transactionID)

	// Assert the result
	require.NoError(t, err, "Expected no error")

	resultsList, err := navClient.GetTransactionStatus(transactionID)
	require.NoError(t, err, "Expected no error")

	// Print result in xml format for debugging
	xmlData, err := xml.MarshalIndent(resultsList, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling to XML: %v\n", err)
		return
	}

	fmt.Println(string(xmlData))
}
