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

  client := ndc.Client{
    Config: "config/ndc-iata-kronos.yml",
  }
}

```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
