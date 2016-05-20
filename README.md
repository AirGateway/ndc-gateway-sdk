# NDC Go SDK

This is a Golang package that wraps any NDC-compliant API.

## Installation

```go get github.com/open-ndc/ndc-go-sdk```

```go build github.com/open-ndc/ndc-go-sdk```

## Usage

```
package main

import "github.com/open-ndc/ndc-go-sdk"

func main() {

  client, err := ndc.NewClient(ndc.ClientOptions{
    ConfigPath: "config/ndc-openndc.yml"
  })

  params := map[string]interface{}{
   "CoreQuery": map[string]interface{}{
      "OriginDestinations": map[string]interface{}{
        "OriginDestination": map[string]interface{} {
          "Departure": map[string]interface{} {
            "AirportCode": "LHR",
            "Date": "2016-05-20",
          },
          "Arrival": map[string]interface{} {
            "AirportCode": "JFK",
          },
        },
      },
    },
  }

  response := client.Request(ndc.Message{
    Method: "AirShopping",
    Params: params,
  })

}

```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
