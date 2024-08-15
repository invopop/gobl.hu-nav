package nav

type InvoiceLines struct {
	MergedItemIndicator bool   `xml:"mergedItemIndicator"` // Indicates if the data report contains aggregated item data
	Lines               []Line `xml:"line"`
}

type Line struct {
	LineNumber                int                       `xml:"lineNumber"`
	LineModificationReference LineModificationReference `xml:"lineModificationReference,omitempty"`
	ReferencesToOtherLines    []ReferenceToOtherLine    `xml:"referencesToOtherLines,omitempty"` // References to related items
	AdvanceData               AdvanceData               `xml:"advanceData,omitempty"`            // Data related to advanced payment
	ProductCodes              []ProductCode             `xml:"productCodes,omitempty"`           // Product codes
	LineExpressionIndicator   bool                      `xml:"lineExpressionIndicator"`          // true if the quantity unit of the item can be expressed as a natural unit of measurement
	LineNatureIndicator       string                    `xml:"lineNatureIndicator,omitempty"`    // Denotes sale of product or service
	LineDescription           string                    `xml:"lineDescription,omitempty"`
	Quantity                  float64                   `xml:"quantity,omitempty"`
	UnitOfMeasure             string                    `xml:"unitOfMeasure,omitempty"`
	UnitOfMeasureOwn          string                    `xml:"unitOfMeasureOwn,omitempty"` // Own quantity unit
	UnitPrice                 float64                   `xml:"unitPrice,omitempty"`
	UnitPriceHUF              float64                   `xml:"unitPriceHUF,omitempty"`
	LineDiscountData          DiscountData              `xml:"lineDiscountData,omitempty"`
	LineAmountsNormal         LineAmountsNormal         `xml:"lineAmountsNormal,omitempty"`     // For normal or aggregate invoices
	LineAmountsSimplified     LineAmountsSimplified     `xml:"lineAmountsSimplified,omitempty"` // For simplified invoices
	IntermediatedService      bool                      `xml:"intermediatedService,omitempty"`  // true if indirect service
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
