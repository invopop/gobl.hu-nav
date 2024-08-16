package nav

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
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
	ProductCodes            *ProductCodes `xml:"productCodes,omitempty"`        // Product codes
	LineExpressionIndicator bool          `xml:"lineExpressionIndicator"`       // true if the quantity unit of the item can be expressed as a natural unit of measurement
	LineNatureIndicator     string        `xml:"lineNatureIndicator,omitempty"` // Denotes sale of product or service
	LineDescription         string        `xml:"lineDescription,omitempty"`
	Quantity                float64       `xml:"quantity,omitempty"`
	UnitOfMeasure           string        `xml:"unitOfMeasure,omitempty"`
	UnitOfMeasureOwn        string        `xml:"unitOfMeasureOwn,omitempty"` // Own quantity unit
	UnitPrice               float64       `xml:"unitPrice,omitempty"`
	UnitPriceHUF            float64       `xml:"unitPriceHUF,omitempty"`
	//LineDiscountData        DiscountData          `xml:"lineDiscountData,omitempty"`
	LineAmountsNormal     *LineAmountsNormal     `xml:"lineAmountsNormal,omitempty"`     // For normal or aggregate invoices
	LineAmountsSimplified *LineAmountsSimplified `xml:"lineAmountsSimplified,omitempty"` // For simplified invoices
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

type LineAmountsNormal struct {
	LineNetAmountData   LineNetAmountData   `xml:"lineNetAmountData"`
	LineVatRate         VatRate             `xml:"lineVatRate"`
	LineVatData         LineVatData         `xml:"lineVatData,omitempty"`
	LineGrossAmountData LineGrossAmountData `xml:"lineGrossAmountData,omitempty"`
}

type LineNetAmountData struct {
	LineNetAmount    float64 `xml:"lineNetAmount"`
	LineNetAmountHUF float64 `xml:"lineNetAmountHUF"`
}

type LineVatData struct {
	LineVatAmount    float64 `xml:"lineVatAmount"`
	LineVatAmountHUF float64 `xml:"lineVatAmountHUF"`
}

type LineGrossAmountData struct {
	LineGrossAmount    float64 `xml:"lineGrossAmount"`
	LineGrossAmountHUF float64 `xml:"lineGrossAmountHUF"`
}

type LineAmountsSimplified struct {
	LineVatRate                  VatRate `xml:"lineVatRate"`
	LineGrossAmountSimplified    float64 `xml:"lineGrossAmountSimplified"`
	LineGrossAmountSimplifiedHUF float64 `xml:"lineGrossAmountSimplifiedHUF"`
}

var codeCategories = []string{
	"VTSZ", "SZJ", "KN", "AHK", "CSK", "KT", "EJ", "TESZOR",
}

var validUnits = map[org.Unit]string{
	org.UnitPiece: "PIECE", org.UnitKilogram: "KILOGRAM", org.UnitMetricTon: "TON", org.UnitKilowattHour: "KWH",
	org.UnitDay: "DAY", org.UnitHour: "HOUR", org.UnitMinute: "MINUTE", org.UnitMonth: "MONTH",
	org.UnitLitre: "LITRE", org.UnitKilometre: "KILOMETER", org.UnitCubicMetre: "CUBIC_METER",
	org.UnitMetre: "METER", org.UnitCarton: "CARTON", org.UnitPackage: "PACK"}

func NewInvoiceLines(inv *bill.Invoice) *InvoiceLines {

	return &InvoiceLines{}
}

func NewLine(line *bill.Line) *Line {
	productCodes := &ProductCodes{}
	if line.Item.Identities != nil {
		productCodes = NewProductCodes(line.Item.Identities)
	}

	lineExpression := false
	lineUnit := ""
	lineUnitOwn := ""
	if line.Item.Unit != "" {
		for unit, value := range validUnits {
			if line.Item.Unit == unit {
				lineExpression = true
				lineUnit = value
				break
			}
		}
		if lineExpression == false {
			lineUnitOwn = string(line.Item.Unit)
		}
	}

	lineNature := ""
	if line.Item.Key != "" {
		if line.Item.Key == "PRODUCT" {
			lineNature = "PRODUCT"
		} else if line.Item.Key == "SERVICE" {
			lineNature = "SERVICE"
		} else {
			lineNature = "OTHER"
		}
	}

	return &Line{
		LineNumber:              line.Index,
		ProductCodes:            productCodes,
		LineExpressionIndicator: lineExpression,
		LineNatureIndicator:     lineNature,
		LineDescription:         line.Item.Name,
		Quantity:                line.Quantity.Float64(),
		UnitOfMeasure:           lineUnit,
		UnitOfMeasureOwn:        lineUnitOwn,
		UnitPrice:               line.Item.Price.Float64(),
	}
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
