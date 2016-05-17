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

type Client struct {
  Endpoint string
  Config string
}

func (r *Client) SendRequest( Method string, Data string ) {
  return
}
