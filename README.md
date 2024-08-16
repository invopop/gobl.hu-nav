# gobl.hu-nav
Convert GOBL into Hungarian NAV XML documents

The invoice data content of the data report must be embedded, encoded in BASE64 format, in the ManageInvoiceRequest/invoiceoperations/invoiceOperation/InvoiceData element.

## Things to include in validation:
- If the 9th digit of the tax id is 5, the group member tax id must exist and should be 4. If the Vat Status is Domestic (Hungarian) always a vat id
- The 9th digit of the Vat Ids must be 1,2,3 or 5 and of the member groups must be 4.

## Limitations

- We don't support batch invoicing (It is used only for batch modifications)
- We don't support modification of invoices
- We don't support fiscal representatives
- We don't support aggregate invoices
- In the VAT rate we are missing the vat amount mismatch field