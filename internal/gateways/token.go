package gateways

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"
)

type TokenInfo struct {
	Token      string
	Expiration time.Time
}

type TokenExchangeRequest struct {
	XMLName  xml.Name     `xml:"TokenExchangeRequest"`
	Common   string       `xml:"xmlns:common,attr"`
	Xmlns    string       `xml:"xmlns,attr"`
	Header   *Header      `xml:"common:header"`
	User     *UserRequest `xml:"common:user"`
	Software *Software    `xml:"software"`
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

func (g *Client) GetToken() error {
	requestData := g.newTokenExchangeRequest()

	token, err := g.postTokenExchangeRequest(requestData)
	if err != nil {
		return err
	}

	err = token.decrypt(g.user.exchangeKey)
	if err != nil {
		return err
	}

	g.token = token

	return nil
}

func (g *Client) postTokenExchangeRequest(requestData TokenExchangeRequest) (*TokenInfo, error) {
	resp, err := g.rest.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/xml").
		SetBody(requestData).
		Post(TokenExchangeEndpoint)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == 200 {
		var tokenExchangeResponse TokenExchangeResponse
		err = xml.Unmarshal(resp.Body(), &tokenExchangeResponse)
		if err != nil {
			return nil, err
		}

		var expirationTime time.Time
		expirationTime, err = time.Parse("2006-01-02T15:04:05.000Z", tokenExchangeResponse.TokenValidityTo)
		if err != nil {
			expirationTime, err = time.Parse("2006-01-02T15:04:05.00Z", tokenExchangeResponse.TokenValidityTo)
			if err != nil {
				return nil, err
			}
		}

		return &TokenInfo{
			Token:      tokenExchangeResponse.EncodedExchangeToken,
			Expiration: expirationTime,
		}, nil
	}

	var generalErrorResponse GeneralErrorResponse
	err = xml.Unmarshal(resp.Body(), &generalErrorResponse)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("error code: %s, message: %s", resp.Status(), generalErrorResponse.Result.ErrorCode)

}

func (g *Client) newTokenExchangeRequest() TokenExchangeRequest {
	timestamp := time.Now().UTC()
	requestID := NewRequestID(timestamp) //This must be unique for each request
	return TokenExchangeRequest{
		Xmlns:    APIXMNLS,
		Common:   APICommon,
		Header:   NewHeader(requestID, timestamp),
		User:     g.NewUser(requestID, timestamp),
		Software: g.software,
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
