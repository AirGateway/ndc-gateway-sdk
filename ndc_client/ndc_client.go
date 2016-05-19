package ndc

import(
	"io/ioutil"

	"github.com/matiasinsaurralde/yaml"
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
	ConfigPath string
}

type Client struct {
	Options ClientOptions
	Config map[string]interface{}
}

func NewClient( options *ClientOptions ) ( *Client, error ) {
	client := &Client{Options: *options}
	client.Config = make(map[string]interface{})
	err := LoadConfig( client.Options.ConfigPath, &client.Config )
	return client, err
}

func LoadConfig( path string, Config *map[string]interface{} ) error {
	RawConfig, err := ioutil.ReadFile( path )
	err = yaml.Unmarshal( RawConfig, *Config )
	return err
}

func (client *Client) Request(message Message) string {
	message.Client = client
	output, _ := message.Prepare()
	return string(output)
}
