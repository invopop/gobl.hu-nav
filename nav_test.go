package nav

import (
	"encoding/base64"
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
		"dev@test.com",
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

	client := NewNav(userID, userPWD, signKey, exchangeKey, taxID, software)

	// Read and encode the invoice XML file
	xmlContent, err := os.ReadFile("examples/example.xml")
	if err != nil {
		t.Fatalf("Failed to read sample invoice file: %v", err)
	}
	encodedInvoice := base64.StdEncoding.EncodeToString(xmlContent)

	// Call the function being tested
	err = client.ReportInvoice(encodedInvoice)

	// Assert the result
	if err != nil {
		t.Errorf("ReportInvoice returned an unexpected error: %v", err)
	}

	require.NoError(t, err, "Expected no error from NewTokenExchangeRequest")
}
