# NDC Go SDK

This is a Golang package that wraps any NDC-compliant API.

## Installation

```go get github.com/open-ndc/ndc-go-sdk```

```go build github.com/open-ndc/ndc-go-sdk```

## Usage

```
package main

import (
  "github.com/ndc-request/ndc-go-sdk"
	"io/ioutil"
  "fmt"
)

func main() {

  client, _ := ndc.NewClient(&ndc.ClientOptions{ConfigPath: "github.com/ndc-request/ndc-go-sdk/config/ndc-openndc.yml"})
  client.HasTemplateVars = true
  params := ndc.Params{
    ndc.Param{
      "Travelers",
      ndc.Params{
        ndc.Param{
          "Traveler",
          ndc.Params{
            ndc.Param{
              "AnonymousTraveler",
              ndc.Params{
                ndc.Param{
                  "PTC", "ADT",
                },
              },
            },
          },
        },
      },
    },
    ndc.Param{
      "CoreQuery",
      ndc.Params{
        ndc.Param{
          "OriginDestinations",
          ndc.Params{
            ndc.Param{
              "OriginDestination",
              ndc.Params{
                ndc.Param{
                  "Departure",
                  ndc.Params{
                    ndc.Param{
                      "AirportCode", "LHR",
                    },
                    ndc.Param{
                      "Date", "2016-05-27",
                    },
                  },
                },
                ndc.Param{
                  "Arrival",
                  ndc.Params{
                    ndc.Param{
                      "AirportCode", "JFK",
                    },
                    ndc.Param{
                      "Date", "2016-05-29",
                    },
                  },
                },
              },
            },
          },
        },
      },
    },
  }
  response := client.Request(ndc.Message{
    Method: "AirShopping",
    Params: params,
  })

  defer response.Body.Close()

  fmt.Println( "-> Receiving response:\n---\n" )
  fmt.Println( response , "\n---\n-> Response body:\n---\n")
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println( string(body) )
  fmt.Println( "\n--\n")
}

```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
