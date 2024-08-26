package gateways

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	charset               = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenExchangeEndpoint = "https://api-test.onlineszamla.nav.gov.hu/invoiceService/v3/tokenExchange"
)

type TokenExchangeRequest struct {
	XMLName  xml.Name  `xml:"TokenExchangeRequest"`
	Common   string    `xml:"xmlns:common,attr"`
	Xmlns    string    `xml:"xmlns,attr"`
	Header   *Header   `xml:"common:header"`
	User     *User     `xml:"common:user"`
	Software *Software `xml:"software"`
}

type TokenExchangeResponse struct {
	XMLName              xml.Name  `xml:"TokenExchangeResponse"`
	Header               *Header   `xml:"header"`
	Result               *Result   `xml:"result"`
	Software             *Software `xml:"software"`
	EncodedExchangeToken string    `xml:"encodedExchangeToken"`
	TokenValidityFrom    string    `xml:"tokenValidityFrom"`
	TokenValidityTo      string    `xml:"tokenValidityTo"`
}

type Result struct {
	FuncCode string `xml:"funcCode"`
}

type GeneralErrorResponse struct {
	XMLName  xml.Name     `xml:"GeneralErrorResponse"`
	Header   *Header      `xml:"header"`
	Result   *ErrorResult `xml:"result"`
	Software *Software    `xml:"software"`
}

type ErrorResult struct {
	FuncCode  string `xml:"funcCode"`
	ErrorCode string `xml:"errorCode"`
	Message   string `xml:"message"`
}

//type Client struct {
// }

func PostTokenExchangeRequest(requestData TokenExchangeRequest) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/xml").
		SetBody(requestData).
		Post(tokenExchangeEndpoint)

	if err != nil {
		return "", err
	}

	if resp.StatusCode() == 200 {
		var tokenExchangeResponse TokenExchangeResponse
		err = xml.Unmarshal(resp.Body(), &tokenExchangeResponse)
		if err != nil {
			return "", err
		}

		return tokenExchangeResponse.EncodedExchangeToken, nil
	}

	var generalErrorResponse GeneralErrorResponse
	err = xml.Unmarshal(resp.Body(), &generalErrorResponse)
	if err != nil {
		return "", err
	}

	return "", fmt.Errorf("error code: %s, message: %s", resp.Status(), generalErrorResponse.Result.ErrorCode)

}

func NewTokenExchangeRequest(userName string, password string, signKey string, taxNumber string, soft *Software) TokenExchangeRequest {
	timestamp := time.Now().UTC()
	requestID := generateRandomString(20) //This cannnot be repeated in the time window of 5 mins
	return TokenExchangeRequest{
		Xmlns:    "http://schemas.nav.gov.hu/OSA/3.0/api",
		Common:   "http://schemas.nav.gov.hu/NTCA/1.0/common",
		Header:   NewHeader(requestID, timestamp),
		User:     NewUser(userName, password, taxNumber, signKey, requestID, timestamp),
		Software: soft,
	}
}

func decryptToken(encodedToken string, keyString string) (string, error) {
	key := []byte(keyString)

	// Decode the base64 encoded encrypted key
	ciphertext, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return "", err
	}

	// Create a buffer for the decrypted text
	decrypted := make([]byte, len(ciphertext))

	// Decrypt the ciphertext
	for bs, be := 0, block.BlockSize(); bs < len(ciphertext); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(decrypted[bs:be], ciphertext[bs:be])
	}

	// Remove padding (if any)
	decrypted = unpad(decrypted)

	return string(decrypted), nil
}

// PKCS7 unpadding
func unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
