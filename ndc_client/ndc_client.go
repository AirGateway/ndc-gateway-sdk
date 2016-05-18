package ndc

import(
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	AcceptedContentType string = "application/xml"
)

var NDCSupportedMethods = map[string] struct{} {
	"AirShoppingRQ": {},
	"FlightPriceRQ": {},
	"SeatAvailabilityRQ": {},
	"ServiceListRQ": {},
	"ServicePriceRQ": {},
	"OrderCreateRQ": {},
	"OrderRetrieveRQ": {},
	"OrderListRQ": {},
	"OrderCancelRQ": {},
	"ItinReshopRQ": {},
}

type ClientOptions struct {
	Endpoint string
	Config string
}

type Client struct {
	Options ClientOptions
	Config map[interface{}]interface{}
}

func NewClient( options *ClientOptions ) ( *Client, error ) {
	client := &Client{Options: *options}
	err := LoadConfig( client.Options.Config, &client.Config )
	return client, err
}

func LoadConfig( path string, Config *map[interface{}]interface{} ) error {
	RawConfig, err := ioutil.ReadFile( path )
	err = yaml.Unmarshal( RawConfig, Config )
	return err
}

func (r *Client) SendRequest( Method string, Data string ) {
  return
}
