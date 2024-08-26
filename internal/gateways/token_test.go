package gateways

import (
	"log"
	"os"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTokenExchangeRequest(t *testing.T) {
	// Set up test data
	software := NewSoftware(
		tax.Identity{Country: l10n.ES.Tax(), Code: cbc.Code("B12345678")},
		"Invopop",
		"ONLINE_SERVICE",
		"1.0.0",
		"TestDev",
		"dev@test.com",
	)

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	userID := os.Getenv("USER_ID")
	userPWD := os.Getenv("USER_PWD")
	signToken := os.Getenv("SIGN_TOKEN")
	taxID := os.Getenv("TAX_ID")

	requestData := NewTokenExchangeData(userID, userPWD, signToken, taxID, software)

	// Execute the function
	token, err := NewTokenExchangeRequest(requestData)

	// Assert results
	require.NoError(t, err, "Expected no error from NewTokenExchangeRequest")
	assert.NotEmpty(t, token, "Expected non-empty token from NewTokenExchangeRequest")
}
