package gateways

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"hash"
	"math/rand"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/invopop/gobl/tax"
	"golang.org/x/crypto/sha3"
)

const (
	charset               = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	tokenExchangeEndpoint = "https://api-test.onlineszamla.nav.gov.hu/invoiceService/v3/tokenExchange"
)

type TokenExchangeRequest struct {
	XMLName  xml.Name `xml:"TokenExchangeRequest"`
	Common   string   `xml:"xmlns:common,attr"`
	Xmlns    string   `xml:"xmlns,attr"`
	Header   Header   `xml:"common:header"`
	User     User     `xml:"common:user"`
	Software Software `xml:"software"`
}

type Header struct {
	RequestId      string `xml:"common:requestId"`
	Timestamp      string `xml:"common:timestamp"`
	RequestVersion string `xml:"common:requestVersion"`
	HeaderVersion  string `xml:"common:headerVersion"`
}

type User struct {
	Login            string           `xml:"common:login"`
	PasswordHash     PasswordHash     `xml:"common:passwordHash"`
	TaxNumber        string           `xml:"common:taxNumber"`
	RequestSignature RequestSignature `xml:"common:requestSignature"`
}

type PasswordHash struct {
	CryptoType string `xml:"cryptoType,attr"`
	Value      string `xml:",chardata"`
}

type RequestSignature struct {
	CryptoType string `xml:"cryptoType,attr"`
	Value      string `xml:",chardata"`
}

type Software struct {
	SoftwareId             string `xml:"softwareId"`
	SoftwareName           string `xml:"softwareName"`
	SoftwareOperation      string `xml:"softwareOperation"`
	SoftwareMainVersion    string `xml:"softwareMainVersion"`
	SoftwareDevName        string `xml:"softwareDevName"`
	SoftwareDevContact     string `xml:"softwareDevContact"`
	SoftwareDevCountryCode string `xml:"softwareDevCountryCode"`
	SoftwareDevTaxNumber   string `xml:"softwareDevTaxNumber"`
}

type TokenExchangeResponse struct {
	XMLName              xml.Name `xml:"TokenExchangeResponse"`
	Header               Header   `xml:"header"`
	Result               Result   `xml:"result"`
	Software             Software `xml:"software"`
	EncodedExchangeToken string   `xml:"encodedExchangeToken"`
	TokenValidityFrom    string   `xml:"tokenValidityFrom"`
	TokenValidityTo      string   `xml:"tokenValidityTo"`
}

type Result struct {
	FuncCode string `xml:"funcCode"`
}

type GeneralErrorResponse struct {
	XMLName  xml.Name    `xml:"GeneralErrorResponse"`
	Header   Header      `xml:"header"`
	Result   ErrorResult `xml:"result"`
	Software Software    `xml:"software"`
}

type ErrorResult struct {
	FuncCode  string `xml:"funcCode"`
	ErrorCode string `xml:"errorCode"`
	Message   string `xml:"message"`
}

//type Client struct {
// }

func NewTokenExchangeRequest(requestData TokenExchangeRequest) (string, error) {
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

func NewTokenExchangeData(userName string, password string, signKey, taxNumber string, soft Software) TokenExchangeRequest {
	timestamp := time.Now().UTC()
	requestID := generateRandomString(20) //This cannnot be repeated in the time window of 5 mins
	return TokenExchangeRequest{
		Xmlns:  "http://schemas.nav.gov.hu/OSA/3.0/api",
		Common: "http://schemas.nav.gov.hu/NTCA/1.0/common",
		Header: Header{
			RequestId:      requestID,
			Timestamp:      timestamp.Format("2006-01-02T15:04:05.000Z"),
			RequestVersion: "3.0",
			HeaderVersion:  "1.0",
		},
		User: User{
			Login:        userName,
			PasswordHash: PasswordHash{CryptoType: "SHA-512", Value: hashPassword(password)},
			TaxNumber:    taxNumber,
			RequestSignature: RequestSignature{
				CryptoType: "SHA3-512",
				Value:      computeRequestSignature(requestID, timestamp, signKey),
			},
		},
		Software: soft,
	}
}

func NewSoftware(taxNumber tax.Identity, name string, operation string, version string, devName string, devContact string) Software {

	if operation != "ONLINE_SERVICE" && operation != "LOCAL_SOFTWARE" {
		operation = "ONLINE_SERVICE"
	}

	return Software{
		SoftwareId:             NewSoftwareID(taxNumber),
		SoftwareName:           name,
		SoftwareOperation:      operation,
		SoftwareMainVersion:    version,
		SoftwareDevName:        devName,
		SoftwareDevContact:     devContact,
		SoftwareDevCountryCode: taxNumber.Country.String(),
		SoftwareDevTaxNumber:   taxNumber.Code.String(),
	}
}

func NewSoftwareID(taxNumber tax.Identity) string {
	// 18-length string:
	//first characters are the country code and tax id
	//the rest is random

	lenRandom := 18 - len(taxNumber.String())

	return taxNumber.String() + generateRandomString(lenRandom)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func hashPassword(password string) string {
	hash := sha512.New()

	return hashInput(hash, password)
}

func computeRequestSignature(requestID string, timestamp time.Time, signKey string) string {
	hash := sha3.New512()

	timeSignature := timestamp.Format("20060102150405")

	return hashInput(hash, requestID+timeSignature+signKey)
}

func hashInput(hash hash.Hash, input string) string {
	hash.Write([]byte(input))

	hashSum := hash.Sum(nil)

	hashHex := hex.EncodeToString(hashSum)

	hashHexUpper := strings.ToUpper(hashHex)

	return hashHexUpper
}
