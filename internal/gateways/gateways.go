package gateways

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	user     *User
	software *Software
	token    *TokenInfo

	rest *resty.Client
}

type User struct {
	login       string
	password    string
	signKey     string
	exchangeKey string
	taxNumber   string
}

// Environment defines the environment to use for connections
type Environment string

// Environment to use for connections
const (
	EnvironmentProduction Environment = "production"
	EnvironmentTesting    Environment = "testing"
	NavProductionURL                  = "https://api.onlineszamla.nav.gov.hu/"
	NavTestingURL                     = "https://api-test.onlineszamla.nav.gov.hu/"

	APIXMNLS  = "http://schemas.nav.gov.hu/OSA/3.0/api"
	APICommon = "http://schemas.nav.gov.hu/NTCA/1.0/common"

	TokenExchangeEndpoint = "invoiceService/v3/tokenExchange"
	ManageInvoiceEndpoint = "invoiceService/v3/manageInvoice"
	StatusEndpoint        = "invoiceService/v3/queryTransactionStatus"
)

// NewGateways creates a new gateways instance
func New(user *User, software *Software, environment Environment) *Client {
	c := &Client{
		user:     user,
		software: software,
	}

	c.rest = resty.New()

	switch environment {
	case EnvironmentProduction:
		c.rest = c.rest.SetBaseURL(NavProductionURL)
	default:
		c.rest = c.rest.SetBaseURL(NavTestingURL)
	}

	return c
}

// NewUser creates a new User instance
func NewUser(login, password, signKey, exchangeKey, taxNumber string) *User {
	return &User{
		login:       login,
		password:    password,
		signKey:     signKey,
		exchangeKey: exchangeKey,
		taxNumber:   taxNumber,
	}
}
