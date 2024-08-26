package gateways

import (
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"math/rand"
	"strings"
	"time"

	"github.com/invopop/gobl/tax"
	"golang.org/x/crypto/sha3"
)

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

func NewHeader(requestID string, timestamp time.Time) *Header {
	return &Header{
		RequestId:      requestID,
		Timestamp:      timestamp.Format("2006-01-02T15:04:05.000Z"),
		RequestVersion: "3.0",
		HeaderVersion:  "1.0",
	}
}

func NewUser(userName string, password string, taxNumber string, signKey string, requestID string, timestamp time.Time, options ...string) *User {
	signature := ""
	if len(options) > 0 {
		base := options[0]
		signature = computeRequestSignature(requestID, timestamp, signKey, base)
	} else {
		signature = computeRequestSignature(requestID, timestamp, signKey)
	}
	return &User{
		Login:        userName,
		PasswordHash: PasswordHash{CryptoType: "SHA-512", Value: hashPassword(password)},
		TaxNumber:    taxNumber,
		RequestSignature: RequestSignature{
			CryptoType: "SHA3-512",
			Value:      signature,
		},
	}
}

func hashPassword(password string) string {
	hash := sha512.New()

	return hashInput(hash, password)
}

func computeRequestSignature(requestID string, timestamp time.Time, signKey string, options ...string) string {
	hash := sha3.New512()

	timeSignature := timestamp.Format("20060102150405")

	hashBase := requestID + timeSignature + signKey

	if len(options) == 0 {
		return hashInput(hash, hashBase)
	}

	hashedInvoice := hashInput(hash, options[0])

	hashBase += hashedInvoice

	return hashInput(hash, hashBase)

}

func NewSoftware(taxNumber tax.Identity, name string, operation string, version string, devName string, devContact string) *Software {

	if operation != "ONLINE_SERVICE" && operation != "LOCAL_SOFTWARE" {
		operation = "ONLINE_SERVICE"
	}

	return &Software{
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

func hashInput(hash hash.Hash, input string) string {
	hash.Write([]byte(input))

	hashSum := hash.Sum(nil)

	hashHex := hex.EncodeToString(hashSum)

	hashHexUpper := strings.ToUpper(hashHex)

	return hashHexUpper
}
