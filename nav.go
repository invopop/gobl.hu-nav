// Package nav is the package used to interact with the NAV API
package nav

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl.hu-nav/internal/doc"
	"github.com/invopop/gobl.hu-nav/internal/gateways"
	"github.com/invopop/gobl/tax"
)

// Client is the main struct for interacting with the NAV API
type Client struct {
	gw  *gateways.Client
	env gateways.Environment
}

// Option is a function used for the different options of the Nav client
// For the moment, the only option is the environment (production or testing)
type Option func(*Client)

// NewNav creates a new Nav client
func NewNav(user *gateways.User, software *gateways.Software, opts ...Option) *Client {

	c := new(Client)

	for _, opt := range opts {
		opt(c)
	}

	c.gw = gateways.New(user, software, c.env)

	return c
}

// InProduction defines the connection to use the production environment.
func InProduction() Option {
	return func(c *Client) {
		c.env = gateways.EnvironmentProduction
	}
}

// InTesting defines the connection to use the testing environment.
func InTesting() Option {
	return func(c *Client) {
		c.env = gateways.EnvironmentTesting
	}
}

// FetchToken fetches the token from the NAV API
func (c *Client) FetchToken() error {
	return c.gw.GetToken()
}

// ReportInvoice reports an invoice to the NAV API
func (c *Client) ReportInvoice(invoice []byte, operationType string) (string, error) {
	return c.gw.ReportInvoice(invoice, operationType)
}

// GetTransactionStatus gets the status of an invoice reporting transaction
func (c *Client) GetTransactionStatus(transactionID string) ([]*gateways.ProcessingResult, error) {
	return c.gw.GetStatus(transactionID)
}

// NewSoftware creates a new Software with the information about the software developer
func NewSoftware(taxNumber tax.Identity, name string, operation string, version string, devName string, devContact string) *gateways.Software {
	return gateways.NewSoftware(taxNumber, name, operation, version, devName, devContact)
}

// NewUser creates a new User
func NewUser(login string, password string, signKey string, exchangeKey string, taxNumber string) *gateways.User {
	return gateways.NewUser(login, password, signKey, exchangeKey, taxNumber)
}

// NewDocument creates a new Nav Document from a GOBL envelope
func NewDocument(env *gobl.Envelope) (*doc.Document, error) {
	return doc.NewDocument(env)
}

// BytesIndent returns the indented XML document bytes
func BytesIndent(doc any) ([]byte, error) {
	buf, err := buffer(doc, xml.Header, true)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buffer(doc any, base string, indent bool) (*bytes.Buffer, error) {
	buf := bytes.NewBufferString(base)

	enc := xml.NewEncoder(buf)
	if indent {
		enc.Indent("", "  ")
	}

	if err := enc.Encode(doc); err != nil {
		return nil, fmt.Errorf("encoding document: %w", err)
	}

	return buf, nil
}
