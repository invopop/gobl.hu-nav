package doc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// InvoiceSummary is the summary of the invoice
// It contains the totals of the invoice
type InvoiceSummary struct {
	SummaryNormal *SummaryNormal `xml:"summaryNormal,omitempty"`
	//SummarySimplified *SummarySimplified `xml:"summarySimplified,omitempty"`
	//SummaryGrossData  *SummaryGrossData  `xml:"summaryGrossData,omitempty"`
}

// SummaryNormal is the summary of a normal invoice
type SummaryNormal struct {
	SummaryByVatRate    []*SummaryByVatRate `xml:"summaryByVatRate"`
	InvoiceNetAmount    string              `xml:"invoiceNetAmount"`
	InvoiceNetAmountHUF string              `xml:"invoiceNetAmountHUF"`
	InvoiceVatAmount    string              `xml:"invoiceVatAmount"`
	InvoiceVatAmountHUF string              `xml:"invoiceVatAmountHUF"`
}

// SummaryByVatRate is the summary of a vat rate
type SummaryByVatRate struct {
	VatRate        *VatRate        `xml:"vatRate"`
	VatRateNetData *VatRateNetData `xml:"vatRateNetData"`
	VatRateVatData *VatRateVatData `xml:"vatRateVatData"`
	//VatRateGrossData VatRateGrossData `xml:"vatRateGrossData, omitempty"`
}

// VatRateNetData is the net data of a vat rate
type VatRateNetData struct {
	VatRateNetAmount    string `xml:"vatRateNetAmount"`
	VatRateNetAmountHUF string `xml:"vatRateNetAmountHUF"`
}

// VatRateVatData is the vat data of a vat rate
type VatRateVatData struct {
	VatRateVatAmount    string `xml:"vatRateVatAmount"`
	VatRateVatAmountHUF string `xml:"vatRateVatAmountHUF"`
}

func newSummaryByVatRate(rate *tax.RateTotal, info *taxInfo, ex num.Amount) (*SummaryByVatRate, error) {
	vatRate, err := newVatRate(rate, info)
	if err != nil {
		return nil, err
	}
	return &SummaryByVatRate{
		VatRate: vatRate,
		VatRateNetData: &VatRateNetData{
			VatRateNetAmount:    rate.Base.Rescale(2).String(),
			VatRateNetAmountHUF: amountToHUF(rate.Base, ex).String(),
		},
		VatRateVatData: &VatRateVatData{
			VatRateVatAmount:    rate.Amount.Rescale(2).String(),
			VatRateVatAmountHUF: amountToHUF(rate.Amount, ex).String(),
		},
	}, nil
}

func newInvoiceSummary(inv *bill.Invoice) (*InvoiceSummary, error) {
	vat := inv.Totals.Taxes.Category(tax.CategoryVAT)
	totalVat := num.MakeAmount(0, 5)
	summaryVat := []*SummaryByVatRate{}
	taxInfo := newTaxInfo(inv)
	ex, err := getExchangeRate(inv)

	if err != nil {
		return nil, err
	}
	for _, rate := range vat.Rates {
		summary, err := newSummaryByVatRate(rate, taxInfo, ex)
		if err != nil {
			return nil, err
		}
		summaryVat = append(summaryVat, summary)
		totalVat = totalVat.Add(rate.Amount)
	}

	return &InvoiceSummary{
		SummaryNormal: &SummaryNormal{
			SummaryByVatRate:    summaryVat,
			InvoiceNetAmount:    inv.Totals.Total.Rescale(2).String(),
			InvoiceNetAmountHUF: amountToHUF(inv.Totals.Total, ex).String(),
			InvoiceVatAmount:    totalVat.Rescale(2).String(),
			InvoiceVatAmountHUF: amountToHUF(totalVat, ex).String(),
		},
	}, nil

}
