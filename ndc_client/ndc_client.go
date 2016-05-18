package ndc

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
}

func NewClient( options *ClientOptions ) *Client {
	client := &Client{Options: *options}
	LoadConfig( client.Options.Config )
	return client
}

func LoadConfig( path string ) map[string]interface{} {
	return nil
}

func (r *Client) SendRequest( Method string, Data string ) {
  return
}
