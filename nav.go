package nav

import (
	"github.com/invopop/gobl.hu-nav/internal/gateways"
	"github.com/invopop/gobl/tax"
)

type nav struct {
	login       string
	password    string
	signKey     string
	exchangeKey string
	taxNumber   string
	software    *gateways.Software
	token       *gateways.TokenInfo
}

func NewNav(login, password, signKey, exchangeKey, taxNumber string, software *gateways.Software) *nav {
	return &nav{
		login:       login,
		password:    password,
		signKey:     signKey,
		exchangeKey: exchangeKey,
		taxNumber:   taxNumber,
		software:    software,
	}
}

func NewSoftware(taxNumber tax.Identity, name string, operation string, version string, devName string, devContact string) *gateways.Software {
	return gateways.NewSoftware(taxNumber, name, operation, version, devName, devContact)
}

func (n *nav) ReportInvoice(invoice string) error {
	// First check if we have a token and it is valid
	if n.token == nil || n.token.Expired() {
		token, err := gateways.GetToken(n.login, n.password, n.signKey, n.exchangeKey, n.taxNumber, n.software)
		if err != nil {
			return err
		}
		n.token = token
	}

	// Now we can report the invoice
	err := gateways.ReportInvoice(n.login, n.password, n.taxNumber, n.signKey, n.token.Token, n.software, invoice)
	if err != nil {
		return err
	}
	return nil
}
