package doc

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSupplierInfo(t *testing.T) {
	// Test case 1: Valid Hungarian supplier without group number
	supplier := &org.Party{
		TaxID: &tax.Identity{
			Country: l10n.HU.Tax(),
			Code:    "98109858",
		},
		Name: "Test Supplier",
		Addresses: []*org.Address{
			{
				Country:  l10n.HU.ISO(),
				Code:     "1234",
				Locality: "Budapest",
				Street:   "Main Street",
			},
		},
	}
	taxNumber := &TaxNumber{TaxPayerID: "98109858"}
	groupNumber := (*TaxNumber)(nil)

	supplierInfo, err := NewSupplierInfo(supplier)
	require.NoError(t, err)
	assert.NotNil(t, supplierInfo)
	assert.Equal(t, taxNumber, supplierInfo.SupplierTaxNumber)
	assert.Nil(t, supplierInfo.GroupMemberTaxNumber)
	assert.Equal(t, "Test Supplier", supplierInfo.SupplierName)
	assert.NotNil(t, supplierInfo.SupplierAddress)

	// Test case 2: Non-Hungarian supplier
	nonHUSupplier := &org.Party{
		TaxID: &tax.Identity{
			Country: l10n.GB.Tax(), // GB for Great Britain
			Code:    "87654321",
		},
		Name: "Non-Hungarian Supplier",
	}

	supplierInfo, err = NewSupplierInfo(nonHUSupplier)
	require.Error(t, err)
	assert.Nil(t, supplierInfo)
	assert.Equal(t, ErrNotHungarian, err)

	supplier = &org.Party{
		TaxID: &tax.Identity{
			Country: l10n.HU.Tax(),
			Code:    "88212131503",
		},
		Name: "Test Supplier",
		Addresses: []*org.Address{
			{
				Country:  l10n.HU.ISO(),
				Code:     "1234",
				Locality: "Budapest",
				Street:   "Main Street",
			},
		},
		Identities: []*org.Identity{
			{
				Code: "21114445423",
			},
		},
	}
	taxNumber = &TaxNumber{TaxPayerID: "88212131", VatCode: "5", CountyCode: "03"}
	groupNumber = &TaxNumber{TaxPayerID: "21114445", VatCode: "4", CountyCode: "23"}

	supplierInfo, err = NewSupplierInfo(supplier)

	require.NoError(t, err)
	assert.NotNil(t, supplierInfo)
	assert.Equal(t, taxNumber, supplierInfo.SupplierTaxNumber)
	assert.Equal(t, groupNumber, supplierInfo.GroupMemberTaxNumber)
	assert.Equal(t, "Test Supplier", supplierInfo.SupplierName)
	assert.NotNil(t, supplierInfo.SupplierAddress)
}
