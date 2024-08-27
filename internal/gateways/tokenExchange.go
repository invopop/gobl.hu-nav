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
	tokenExchangeEndpoint = "https://api-test.onlineszamla.nav.gov.hu/invoiceService/v3/tokenExchange"
)

type TokenInfo struct {
	Token      string
	Expiration time.Time
}

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

func GetToken(userName string, password string, signKey string, exchangeKey string, taxNumber string, soft *Software) (*TokenInfo, error) {
	requestData := newTokenExchangeRequest(userName, password, signKey, taxNumber, soft)
	token, err := postTokenExchangeRequest(requestData)
	if err != nil {
		return nil, err
	}

	err = token.decrypt(exchangeKey)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func postTokenExchangeRequest(requestData TokenExchangeRequest) (*TokenInfo, error) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/xml").
		SetBody(requestData).
		Post(tokenExchangeEndpoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 200 {
		var tokenExchangeResponse TokenExchangeResponse
		err = xml.Unmarshal(resp.Body(), &tokenExchangeResponse)
		if err != nil {
			return nil, err
		}

		time, err := time.Parse("2006-01-02T15:04:05.000Z", tokenExchangeResponse.TokenValidityTo)
		if err != nil {
			return nil, err
		}

		return &TokenInfo{
			Token:      tokenExchangeResponse.EncodedExchangeToken,
			Expiration: time,
		}, nil
	}

	var generalErrorResponse GeneralErrorResponse
	err = xml.Unmarshal(resp.Body(), &generalErrorResponse)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("error code: %s, message: %s", resp.Status(), generalErrorResponse.Result.ErrorCode)

}

func newTokenExchangeRequest(userName string, password string, signKey string, taxNumber string, soft *Software) TokenExchangeRequest {
	timestamp := time.Now().UTC()
	requestID := generateRandomString(20) //This must be unique for each request
	return TokenExchangeRequest{
		Xmlns:    "http://schemas.nav.gov.hu/OSA/3.0/api",
		Common:   "http://schemas.nav.gov.hu/NTCA/1.0/common",
		Header:   NewHeader(requestID, timestamp),
		User:     NewUser(userName, password, taxNumber, signKey, requestID, timestamp),
		Software: soft,
	}
}

func (tok *TokenInfo) Expired() bool {
	return time.Now().After(tok.Expiration)
}

func (tokenInfo *TokenInfo) decrypt(keyString string) error {
	key := []byte(keyString)

	// Decode the base64 encoded encrypted key
	ciphertext, err := base64.StdEncoding.DecodeString(tokenInfo.Token)
	if err != nil {
		return err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create a buffer for the decrypted text
	decrypted := make([]byte, len(ciphertext))

	// Decrypt the ciphertext
	for bs, be := 0, block.BlockSize(); bs < len(ciphertext); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(decrypted[bs:be], ciphertext[bs:be])
	}

	// Remove padding (if any)
	decrypted = unpad(decrypted)

	tokenInfo.Token = string(decrypted)

	return nil
}

// PKCS7 unpadding
func unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
