# gobl.hu-nav
Go library to convert [GOBL](https://github.com/invopop/gobl) invoices into TicketBAI declarations and send them to the Hungarian web services.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License v2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). For contributions to this library to be accepted, we will require you to accept transferring your copyright to Invopop Ltd.

## Usage

### Go package

#### Conversion
Usage of the XInvoice conversion library is quite straight forward. You must first have a GOBL Envelope including an invoice ready to convert.

```go
package main

import (
    "os"

    "github.com/invopop/gobl"
    nav "github.com/invopop/gobl.hu-nav"
)

func main() {
    data, _ := os.ReadFile("./test/data/invoice-test.json")

    env := new(gobl.Envelope)
    if err := json.Unmarshal(data, env); err != nil {
        panic(err)
    }

    // Prepare the Nav document
    doc, err := nav.NewDocument(env)
    if err != nil {
        panic(err)
    }

    // Create the XML output
    out, err := nav.BytesIndent(doc)
    if err != nil {
        panic(err)
    }

    // TODO: do something with the output
}
```

#### Invoice Reporting

Once the invoice is generated, it can be reported to the Hungarian authoritites. You must first have a technical user created in the [Online Szamla](https://onlineszamla.nav.gov.hu/home).

```go
package main

import (
    "os"

    "github.com/invopop/gobl"
    nav "github.com/invopop/gobl.hu-nav"
)

func main() {

    // Software is the information regarding the system used to report the invoices
    software := NewSoftware(
		tax.Identity{Country: l10n.ES.Tax(), Code: cbc.Code("B12345678")},
		"Invopop",
		"ONLINE_SERVICE",
		"1.0.0",
		"TestDev",
		"test@dev.com",
	)

    // User is all the data obtained from the technical user that it is needed to report the invoices
    user := NewUser(
        "username",
        "password",
        "signature_key",
        "exchange_key",
        "taxID"
    )

    // Create a new client with the user and software data and choose if you want to issue the invoices in the testing or production environment
    navClient := NewNav(user, software, InTesting())

    //We load the invoice
    inv, err := os.ReadFile("test/data/out/output.xml")
	if err != nil {
		panic(err)
	}

    // Report the invoice
    transactionId, err := navClient.ReportInvoice(invoice)
    if err != nil {
        panic(err)
    }

    // Once the invoice is reported, you can check the status
    // If you check the status too early you would get a status of PROCESSING, which means that you should try again later to query the status
    resultsList, err := navClient.GetTransactionStatus(transactionId)

    //The output contains the status and a list of technical and business validation messages. To visualize the output, you can create a XML output:
    out, err := nav.BytesIndent(resultsList)
    if err != nil {
        panic(err)
    }

    // TODO: do something with the output
}
```

### Command Line
#### Conversion

The GOBL NAV package tool also includes a command line helper. You can install manually in your Go environment with:

```bash
go install ./cmd/gobl.nav
```

Usage is very straightforward:

```bash
gobl.nav convert ./test/data/invoice.json
```


## Limitations/Things to do

### Invoice Modification
- For invoice modification the only step left is to get the line number of the invoice that we want to modify and include it in the field `LineModificationReference`

### Doc Conversion
- Batch invoicing not supported
- Support fiscal representatives
- Aggregate invoices not supported
- Product refund charges not supported (Field Product Fee Summary in the Invoice)
- Nav supports 100 invoice creation/modification in the same request. For the moment, we only support 1 invoice at each request.


