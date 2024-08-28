# gobl.hu-nav
Convert GOBL into Hungarian NAV XML documents

The invoice data content of the data report must be embedded, encoded in BASE64 format, in the ManageInvoiceRequest/invoiceoperations/invoiceOperation/InvoiceData element.

## Limitations

- We don't support batch invoicing (It is used only for batch modifications)
- We don't support modification of invoices
- We don't support fiscal representatives
- We don't support aggregate invoices
- In the VAT rate we are missing the vat amount mismatch field (used when VAT has been charged under section 11 or 14)
- We don't support refund product charges (Field Product Fee Summary in the Invoice)

- For each requestID, it is needed to have a new id for each request (non-repeat). For the moment, this is done just at random but there is a small probability that the request IDs match.

- Nav supports 100 invoice creation/modification in the same request. For the moment, we only support 1 invoice at each request. If this is changed, it should also be changed in status.go, as now status only answers with the errors/warning of the first invoice.

## Invoice issuing order

1. Generate the doc with the format required by the NAV.
2. Obtain an exchange token with 5 mins validity.
3. Issue the invoice (An Ok message when issuing the invoice doesn't mean that it is correctly issued)
4. Check your transaction status (This would give us information about the status of the invoice issuing)

If you do step number 4 just after 3, you would get a status of processing, we should retry the operation until we get a status of DONE or ABORT.

