package nav

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
)

// Depends wether the invoice is simplified or not
type InvoiceSummary struct {
	SummaryNormal *SummaryNormal `xml:"summaryNormal"`
	// This is to differentiate between normal or simplified invoice, for the moment we are only doing normal
	//SummarySimplified SummarySimplified `xml:"summarySimplified,omitempty"`
	//SummaryGrossData  SummaryGrossData  `xml:"summaryGrossData,omitempty"`
}

type SummaryNormal struct {
	SummaryByVatRate    []*SummaryByVatRate `xml:"summaryByVatRate"`
	InvoiceNetAmount    float64             `xml:"invoiceNetAmount"`
	InvoiceNetAmountHUF float64             `xml:"invoiceNetAmountHUF"`
	InvoiceVatAmount    float64             `xml:"invoiceVatAmount"`
	InvoiceVatAmountHUF float64             `xml:"invoiceVatAmountHUF"`
}

type SummaryByVatRate struct {
	VatRate        *VatRate        `xml:"vatRate"`
	VatRateNetData *VatRateNetData `xml:"vatRateNetData"`
	VatRateVatData *VatRateVatData `xml:"vatRateVatData"`
	//VatRateGrossData VatRateGrossData `xml:"vatRateGrossData, omitempty"`
}

type VatRateNetData struct {
	VatRateNetAmount    float64 `xml:"vatRateNetAmount"`
	VatRateNetAmountHUF float64 `xml:"vatRateNetAmountHUF"`
}

type VatRateVatData struct {
	VatRateVatAmount    float64 `xml:"vatRateVatAmount"`
	VatRateVatAmountHUF float64 `xml:"vatRateVatAmountHUF"`
}

func newSummaryByVatRate(rate *tax.RateTotal, info *taxInfo, ex float64) *SummaryByVatRate {
	return &SummaryByVatRate{
		VatRate: NewVatRate(rate, info),
		VatRateNetData: &VatRateNetData{
			VatRateNetAmount:    rate.Base.Float64(),
			VatRateNetAmountHUF: rate.Base.Float64() * ex,
		},
		VatRateVatData: &VatRateVatData{
			VatRateVatAmount:    rate.Amount.Float64(),
			VatRateVatAmountHUF: rate.Amount.Float64() * ex,
		},
	}
}

func NewInvoiceSummary(inv *bill.Invoice) (*InvoiceSummary, error) {
	vat := inv.Totals.Taxes.Category(tax.CategoryVAT)
	totalVat := 0.0
	summaryVat := []*SummaryByVatRate{}
	taxInfo := newTaxInfo(inv)
	ex, err := getInvoiceRate(inv)
	if err != nil {
		return nil, err
	}
	for _, rate := range vat.Rates {
		summaryVat = append(summaryVat, newSummaryByVatRate(rate, taxInfo, ex))
		totalVat += rate.Amount.Float64()
	}

	return &InvoiceSummary{
		SummaryNormal: &SummaryNormal{
			SummaryByVatRate:    summaryVat,
			InvoiceNetAmount:    inv.Totals.Total.Float64(),
			InvoiceNetAmountHUF: inv.Totals.Total.Float64() * ex,
			InvoiceVatAmount:    totalVat,
			InvoiceVatAmountHUF: totalVat * ex,
		},
	}, nil

}
