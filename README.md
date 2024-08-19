# gobl.hu-nav
Convert GOBL into Hungarian NAV XML documents

The invoice data content of the data report must be embedded, encoded in BASE64 format, in the ManageInvoiceRequest/invoiceoperations/invoiceOperation/InvoiceData element.

## Limitations

- We don't support batch invoicing (It is used only for batch modifications)
- We don't support modification of invoices
- We don't support fiscal representatives
- We don't support aggregate invoices
- In the VAT rate we are missing the vat amount mismatch field (used when VAT has been charged under section 11 or 14)