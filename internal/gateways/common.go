package gateways

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/xml"
	"hash"
	"math/rand"
	"strings"
	"time"

	"github.com/invopop/gobl/tax"
	"golang.org/x/crypto/sha3"
)

const (
	charset        = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	RequestVersion = "3.0"
	HeaderVersion  = "1.0"
)

// Header is the common header for all requests
// A new RequestId and Timestamp is generated for each request
type Header struct {
	RequestId      string `xml:"common:requestId"`
	Timestamp      string `xml:"common:timestamp"`
	RequestVersion string `xml:"common:requestVersion"`
	HeaderVersion  string `xml:"common:headerVersion"`
}

// UserRequest is the common user for all requests. Login, PasswordHash and TaxNumber remain the same,
// while RequestSignature is computed differently for manageInvoice.
type UserRequest struct {
	Login            string           `xml:"common:login"`
	PasswordHash     PasswordHash     `xml:"common:passwordHash"`
	TaxNumber        string           `xml:"common:taxNumber"`
	RequestSignature RequestSignature `xml:"common:requestSignature"`
}

// PasswordHash is the hash of the password
type PasswordHash struct {
	CryptoType string `xml:"cryptoType,attr"`
	Value      string `xml:",chardata"`
}

// RequestSignature is the signature of the request. It is computed differently for manageInvoice.
type RequestSignature struct {
	CryptoType string `xml:"cryptoType,attr"`
	Value      string `xml:",chardata"`
}

// Software is the information about the software used for issuing the invoices
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

// GeneralErrorResponse is the common error response for all requests
type GeneralErrorResponse struct {
	XMLName  xml.Name  `xml:"GeneralErrorResponse"`
	Header   *Header   `xml:"header"`
	Result   *Result   `xml:"result"`
	Software *Software `xml:"software"`
}

// Result is the common result for all requests
// If request is OK, only FuncCode is returned
type Result struct {
	FuncCode      string         `xml:"funcCode"`
	ErrorCode     string         `xml:"errorCode,omitempty"`
	Message       string         `xml:"message,omitempty"`
	Notifications *Notifications `xml:"notifications,omitempty"`
}

// Notifications is the list of notifications
type Notifications struct {
	Notification []*Notification `xml:"notification"`
}

// Notification includes a code and a text
type Notification struct {
	NotificationCode string `xml:"notificationCode"`
	NotificationText string `xml:"notificationText"`
}

// NewHeader creates a new Header with the given requestID and timestamp
func NewHeader(requestID string, timestamp time.Time) *Header {
	return &Header{
		RequestId:      requestID,
		Timestamp:      timestamp.Format("2006-01-02T15:04:05.00Z"),
		RequestVersion: RequestVersion,
		HeaderVersion:  HeaderVersion,
	}
}

// NewUser creates a new User
func (g *Client) NewUser(requestID string, timestamp time.Time, options ...string) *UserRequest {
	signature := ""
	if len(options) > 0 {
		base := options[0]
		signature = computeRequestSignature(requestID, timestamp, g.user.signKey, base)
	} else {
		signature = computeRequestSignature(requestID, timestamp, g.user.signKey)
	}
	return &UserRequest{
		Login:        g.user.login,
		PasswordHash: PasswordHash{CryptoType: "SHA-512", Value: hashPassword(g.user.password)},
		TaxNumber:    g.user.taxNumber,
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
	timeSignature := timestamp.Format("20060102150405")

	hashBase := requestID + timeSignature + signKey

	if len(options) == 0 {
		return hashInput(sha3.New512(), hashBase)
	}

	hashedInvoice := hashInput(sha3.New512(), options[0])

	hashBase += hashedInvoice

	return hashInput(sha3.New512(), hashBase)

}

// NewSoftware creates a new Software with the information about the software developer
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

func NewRequestID(timestamp time.Time) string {
	timeUnique := timestamp.Format("20060102150405")

	randomNumber := rand.Intn(17)

	return timeUnique + generateRandomString(randomNumber)
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
