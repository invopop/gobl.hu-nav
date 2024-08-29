package doc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

type InvoiceLines struct {
	MergedItemIndicator bool   `xml:"mergedItemIndicator"` // Indicates if the data report contains aggregated item data
	Lines               []Line `xml:"line"`
}

type Line struct {
	LineNumber int `xml:"lineNumber"`
	//LineModificationReference LineModificationReference `xml:"lineModificationReference,omitempty"`
	//ReferencesToOtherLines    []ReferenceToOtherLine    `xml:"referencesToOtherLines,omitempty"` // References to related items
	//AdvanceData               AdvanceData               `xml:"advanceData,omitempty"`            // Data related to advanced payment
	ProductCodes            *ProductCodes          `xml:"productCodes,omitempty"`        // Product codes
	LineExpressionIndicator bool                   `xml:"lineExpressionIndicator"`       // true if the quantity unit of the item can be expressed as a natural unit of measurement
	LineNatureIndicator     string                 `xml:"lineNatureIndicator,omitempty"` // Denotes sale of product or service
	LineDescription         string                 `xml:"lineDescription,omitempty"`
	Quantity                string                 `xml:"quantity,omitempty"`
	UnitOfMeasure           string                 `xml:"unitOfMeasure,omitempty"`
	UnitOfMeasureOwn        string                 `xml:"unitOfMeasureOwn,omitempty"` // Own quantity unit
	UnitPrice               string                 `xml:"unitPrice,omitempty"`
	UnitPriceHUF            string                 `xml:"unitPriceHUF,omitempty"`
	LineDiscountData        *LineDiscountData      `xml:"lineDiscountData,omitempty"`
	LineAmountsNormal       *LineAmountsNormal     `xml:"lineAmountsNormal,omitempty"`     // For normal or aggregate invoices
	LineAmountsSimplified   *LineAmountsSimplified `xml:"lineAmountsSimplified,omitempty"` // For simplified invoices
	//IntermediatedService  bool                  `xml:"intermediatedService,omitempty"`  // true if indirect service
	//AggregateInvoiceLineData  AggregateInvoiceLineData  `xml:"aggregateInvoiceLineData,omitempty"` // Aggregate invoice data
	//NewTransportMean          NewTransportMean          `xml:"newTransportMean,omitempty"`         // Sale of new means of transport
	//DepositIndicator          bool                      `xml:"depositIndicator,omitempty"`         // true if the item is a deposit
	//ObligatedForProductFee    bool                      `xml:"obligatedForProductFee,omitempty"`   // true if a product fee obligation applies to the line item
	//GPCExcise                 float64                   `xml:",omitempty"`                         // Excise tax on natural gas
	//DieselOilPurchase         DieselOilPurchase         `xml:"dieselOilPurchase,omitempty"`        // Data on post-tax purchase diesel oil
	//NetaDeclaration           bool                      `xml:"netaDeclaration,omitempty"`          // true if the tax liability determined by the Public Health Product Tax falls on the taxpayer
	//ProductFeeClause          ProductFeeClause          `xml:"productFeeClause,omitempty"`         // Clauses on environmental product charges
	//LineProductFeeContent     LineProductFeeContent     `xml:"lineProductFeeContent,omitempty"`    // Data on product fee content
	//ConventionalLineInfo      ConventionalLineInfo      `xml:"conventionalLineInfo,omitempty"`     // Other conventional named data
	//AdditionalLineData        AdditionalData            `xml:"additionalLineData,omitempty"`       // Additional data
}

type ProductCodes struct {
	ProductCode []*ProductCode `xml:"productCode"`
}

// One of code value or codeownvalue must be present
type ProductCode struct {
	ProductCodeCategory string `xml:"productCodeCategory"` // Product code value for non-own product codes
	ProductCodeValue    string `xml:"productCodeValue,omitempty"`
	ProductCodeOwnValue string `xml:"productCodeOwnValue,omitempty"` // Own product code value
}

type LineDiscountData struct {
	DiscountDescription string `xml:"discountDescription"`
	DiscountValue       string `xml:"discountValue"`
	DiscountRate        string `xml:"discountRate"`
}

type LineAmountsNormal struct {
	LineNetAmountData   *LineNetAmountData   `xml:"lineNetAmountData"`
	LineVatRate         *VatRate             `xml:"lineVatRate"`
	LineVatData         *LineVatData         `xml:"lineVatData,omitempty"`
	LineGrossAmountData *LineGrossAmountData `xml:"lineGrossAmountData,omitempty"`
}

type LineNetAmountData struct {
	LineNetAmount    string `xml:"lineNetAmount"`
	LineNetAmountHUF string `xml:"lineNetAmountHUF"`
}

type LineVatData struct {
	LineVatAmount    string `xml:"lineVatAmount"`
	LineVatAmountHUF string `xml:"lineVatAmountHUF"`
}

// LineGrossAmountData is the Net amount + VAT amount (Not mandatory)
type LineGrossAmountData struct {
	LineGrossAmount    string `xml:"lineGrossAmount"`
	LineGrossAmountHUF string `xml:"lineGrossAmountHUF"`
}

type LineAmountsSimplified struct {
	LineVatRate                  *VatRate `xml:"lineVatRate"`
	LineGrossAmountSimplified    string   `xml:"lineGrossAmountSimplified"` //This amount is the net amount of the normal line
	LineGrossAmountSimplifiedHUF string   `xml:"lineGrossAmountSimplifiedHUF"`
}

var codeCategories = []string{
	"VTSZ", "SZJ", "KN", "AHK", "CSK", "KT", "EJ", "TESZOR",
}

var validUnits = map[org.Unit]string{
	org.UnitPiece: "PIECE", org.UnitKilogram: "KILOGRAM", org.UnitMetricTon: "TON", org.UnitKilowattHour: "KWH",
	org.UnitDay: "DAY", org.UnitHour: "HOUR", org.UnitMinute: "MINUTE", org.UnitMonth: "MONTH",
	org.UnitLitre: "LITRE", org.UnitKilometre: "KILOMETER", org.UnitCubicMetre: "CUBIC_METER",
	org.UnitMetre: "METER", org.UnitCarton: "CARTON", org.UnitPackage: "PACK"}

func NewInvoiceLines(inv *bill.Invoice) (*InvoiceLines, error) {

	invoiceLines := &InvoiceLines{}
	taxinfo := newTaxInfo(inv)
	rate, err := getInvoiceRate(inv)
	if err != nil {
		return nil, err
	}
	for _, line := range inv.Lines {
		invoiceLine, err := NewLine(line, taxinfo, rate)
		if err != nil {
			return nil, err
		}
		invoiceLines.Lines = append(invoiceLines.Lines, *invoiceLine)
	}
	invoiceLines.MergedItemIndicator = false

	return invoiceLines, nil
}

func NewLine(line *bill.Line, info *taxInfo, rate float64) (*Line, error) {
	lineNav := &Line{
		LineNumber:              line.Index,
		LineExpressionIndicator: false,
		LineDescription:         line.Item.Name,
		UnitPrice:               line.Item.Price.String(),
		UnitPriceHUF:            amountToHUF(line.Item.Price, rate).String(),
		Quantity:                line.Quantity.String(),
	}

	if line.Item.Identities != nil {
		lineNav.ProductCodes = NewProductCodes(line.Item.Identities)
	}

	if line.Item.Unit != "" {
		for unit, value := range validUnits {
			if line.Item.Unit == unit {
				lineNav.LineExpressionIndicator = true
				lineNav.UnitOfMeasure = value
				break
			}
		}
		if !lineNav.LineExpressionIndicator {
			lineNav.UnitOfMeasureOwn = string(line.Item.Unit)
		}
	}

	if line.Item.Key != "" {
		if line.Item.Key == "PRODUCT" {
			lineNav.LineNatureIndicator = "PRODUCT"
		} else if line.Item.Key == "SERVICE" {
			lineNav.LineNatureIndicator = "SERVICE"
		} else {
			lineNav.LineNatureIndicator = "OTHER"
		}
	}

	if line.Discounts != nil {
		discount := &LineDiscountData{}
		discount.DiscountDescription = ""
		discountValue := 0.0
		for _, dis := range line.Discounts {
			discount.DiscountDescription += dis.Reason + ". "
			discountValue += dis.Amount.Float64()
		}
		discount.DiscountValue = num.AmountFromFloat64(discountValue, 2).String()
		lineNav.LineDiscountData = discount
	}

	vatCombo := line.Taxes.Get(tax.CategoryVAT)
	if vatCombo != nil {
		if info.simplifiedInvoice {
			vatAmount := line.Total.Multiply(vatCombo.Percent.Amount())
			lineNav.LineAmountsSimplified = &LineAmountsSimplified{
				LineVatRate:                  &VatRate{VatContent: vatAmount.Rescale(2).String()},
				LineGrossAmountSimplified:    line.Total.Rescale(2).String(),
				LineGrossAmountSimplifiedHUF: amountToHUF(line.Total, rate).String(),
			}
		} else {
			vatRate, err := NewVatRate(vatCombo, info)
			if err != nil {
				return nil, err
			}
			lineNav.LineAmountsNormal = &LineAmountsNormal{
				LineNetAmountData: &LineNetAmountData{
					LineNetAmount:    line.Total.Rescale(2).String(),
					LineNetAmountHUF: amountToHUF(line.Total, rate).String(),
				},
				LineVatRate: vatRate,
			}
		}
	}
	return lineNav, nil
}

func NewProductCodes(identities []*org.Identity) *ProductCodes {
	if len(identities) == 0 {
		return nil
	}
	productCodes := &ProductCodes{}
	for _, identity := range identities {
		productCode := NewProductCode(identity)
		productCodes.ProductCode = append(productCodes.ProductCode, productCode)
	}
	return productCodes
}

func NewProductCode(identity *org.Identity) *ProductCode {
	if identity.Type == "OWN" {
		return &ProductCode{
			ProductCodeCategory: "OWN",
			ProductCodeOwnValue: identity.Code.String(),
		}
	}
	for _, category := range codeCategories {
		if identity.Type == cbc.Code(category) {
			return &ProductCode{
				ProductCodeCategory: category,
				ProductCodeValue:    identity.Code.String(),
			}
		}
	}
	return &ProductCode{
		ProductCodeCategory: "OTHER",
		ProductCodeValue:    identity.Code.String(),
	}
}

func amountToHUF(amount num.Amount, ex float64) num.Amount {
	result := amount.Multiply(num.AmountFromFloat64(ex, 5))
	return result.Rescale(2)
}
