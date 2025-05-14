# Heleket API Go Wrapper

This repository contains an **unofficial Go wrapper** for the Heleket API, a crypto payment gateway. This wrapper simplifies the process of integrating Heleket functionality into your Go projects.

## Features

- Easy-to-use Go interface for Heleket API
- Supports payment and payout operations
- Handles static wallet functionalities
- Supports refund operations
- Supports resending and testing webhook requests
- Supports verifying signature
- Provides strongly typed responses

## Installation

To install the Heleket API Go wrapper, use `go get`:

```
go get github.com/rmilansky/go-heleket
```

## Usage

Here's a quick example of how to use the wrapper:

```go

import (
    "fmt"
    "github.com/rmilansky/go-heleket"
)

func main() {
    httpClient := http.DefaultClient
    client := heleket.New(httpClient, "your-merchant-id", "your-payment-api-key", "your-payout-api-key")
    
    // Create an invoice
    invoiceReq := &heleket.InvoiceRequest{
        Amount: "10",
        Currency: "USD",
        OrderId: "your-order-id",
        InvoiceRequestOptions: &heleket.invoiceRequestOptions{
            Network: "tron",
            UrlCallback: "https://yourdomain.com/callback"
        },
    }
    invoice, err := heleket.CreateInvoice(invoiceReq)
    if err != nil {
        // Handle error
    }
    
    fmt.Printf("Invoice created: %+v\n", invoice)
}
```

## API Coverage

This wrapper currently supports the following Heleket API functionalities:

- Payment operations
- Static wallet operations
- Refund operations
- Resending webhook requests

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
