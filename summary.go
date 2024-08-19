package nav

// Depends wether the invoice is simplified or not
type InvoiceSummary struct {
	SummaryNormal *SummaryNormal `xml:"summaryNormal,omitempty"`
	//SummarySimplified *SummarySimplified `xml:"summarySimplified,omitempty"`
	//SummaryGrossData  *SummaryGrossData  `xml:"summaryGrossData,omitempty"`
}

type SummaryNormal struct {
	SummaryByVatRate    []*SummaryByVatRate `xml:"summaryByVatRate"`
	InvoiceNetAmount    string              `xml:"invoiceNetAmount"`
	InvoiceNetAmountHUF string              `xml:"invoiceNetAmountHUF"`
	InvoiceVatAmount    string              `xml:"invoiceVatAmount"`
	InvoiceVatAmountHUF string              `xml:"invoiceVatAmountHUF"`
}

type SummaryByVatRate struct {
	VatRate        *VatRate        `xml:"vatRate"`
	VatRateNetData *VatRateNetData `xml:"vatRateNetData"`
	VatRateVatData *VatRateVatData `xml:"vatRateVatData"`
	//VatRateGrossData VatRateGrossData `xml:"vatRateGrossData, omitempty"`
}

type VatRateNetData struct {
	VatRateNetAmount    string `xml:"vatRateNetAmount"`
	VatRateNetAmountHUF string `xml:"vatRateNetAmountHUF"`
}

type VatRateVatData struct {
	VatRateVatAmount    string `xml:"vatRateVatAmount"`
	VatRateVatAmountHUF string `xml:"vatRateVatAmountHUF"`
}

/*func newSummaryByVatRate(rate *tax.RateTotal, info *taxInfo, ex float64) *SummaryByVatRate {
	return &SummaryByVatRate{
		VatRate: NewVatRate(rate, info),
		VatRateNetData: &VatRateNetData{
			VatRateNetAmount:    rate.Base.Rescale(2).String(),
			VatRateNetAmountHUF: amountToHUF(rate.Base, ex),
		},
		VatRateVatData: &VatRateVatData{
			VatRateVatAmount:    rate.Amount.Rescale(2).String(),
			VatRateVatAmountHUF: amountToHUF(rate.Amount, ex),
		},
	}
}*/

/*func NewInvoiceSummary(inv *bill.Invoice) (*InvoiceSummary, error) {
	vat := inv.Totals.Taxes.Category(tax.CategoryVAT)
	totalVat := num.MakeAmount(0, 5)
	summaryVat := []*SummaryByVatRate{}
	taxInfo := newTaxInfo(inv)
	ex, err := getInvoiceRate(inv)

	if err != nil {
		return nil, err
	}
	for _, rate := range vat.Rates {
		summaryVat = append(summaryVat, newSummaryByVatRate(rate, taxInfo, ex))
		totalVat = totalVat.Add(rate.Amount)
	}

	return &InvoiceSummary{
		SummaryNormal: &SummaryNormal{
			SummaryByVatRate:    summaryVat,
			InvoiceNetAmount:    inv.Totals.Total.Rescale(2).String(),
			InvoiceNetAmountHUF: amountToHUF(inv.Totals.Total, ex),
			InvoiceVatAmount:    totalVat.Rescale(2).String(),
			InvoiceVatAmountHUF: amountToHUF(totalVat, ex),
		},
	}, nil

}*/
