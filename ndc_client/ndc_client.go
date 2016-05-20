package ndc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"bytes"

	"github.com/matiasinsaurralde/yaml"
)

var NDCSupportedMethods = map[string]struct{}{
	"AirShoppingRQ":      {},
	"FlightPriceRQ":      {},
	"SeatAvailabilityRQ": {},
	"ServiceListRQ":      {},
	"ServicePriceRQ":     {},
	"OrderCreateRQ":      {},
	"OrderRetrieveRQ":    {},
	"OrderListRQ":        {},
	"OrderCancelRQ":      {},
	"ItinReshopRQ":       {},
}

var TemplateVars = []string{"request_name"}

type ClientOptions struct {
	Endpoint   string
	ConfigPath string
}

type Client struct {
	Options         ClientOptions
	HasTemplateVars bool
	Config          map[string]interface{}
	RawConfig       []byte
	HttpClient      *http.Client
}

func NewClient(options *ClientOptions) (*Client, error) {
	client := &Client{Options: *options}
	client.Config = make(map[string]interface{})
	client.HttpClient = &http.Client{}
	client.HasTemplateVars = false
	err := client.LoadConfig()
	return client, err
}

func ConfigHasTemplateVars(RawConfig *[]byte) int {
	config := string(*RawConfig)
	var matches = 0
	for _, VarName := range TemplateVars {
		VarIndex := strings.Index(config, VarName)
		if VarIndex > 0 {
			matches++
		}
	}
	return matches
}

func (client *Client) LoadConfig() error {
	config, err := ioutil.ReadFile(client.Options.ConfigPath)
	client.RawConfig = config
	err = yaml.Unmarshal(client.RawConfig, client.Config)

	if ConfigHasTemplateVars(&client.RawConfig) > 0 {
		client.HasTemplateVars = true
	}

	return err
}

func (client *Client) PrepareConfig(message Message) (Config map[string]interface{}) {
	Config = make(map[string]interface{})
	ModifiedConfig := string(client.RawConfig)

	for _, VarName := range TemplateVars {

		var VarValue = ""

		switch VarName {
		case "request_name":
			VarValue = message.Method

		}

		VarName = fmt.Sprintf("{{%s}}", VarName)

		ModifiedConfig = strings.Replace(ModifiedConfig, VarName, VarValue, -1)
	}
	yaml.Unmarshal([]byte(ModifiedConfig), Config)
	return
}

func (client *Client) AppendHeaders(r *http.Request, HeadersConfig interface{}) {
	headers := HeadersConfig.(map[string]interface{})
	for Header, Value := range headers {
		r.Header.Add(Header, Value.(string))
	}
}

func (client *Client) Request(message Message) *http.Response {

	var Config map[string]interface{}

	message.Client = client

	body, _ := message.Prepare()

	if client.HasTemplateVars {
		Config = client.PrepareConfig(message)
	} else {
		Config = client.Config
	}

	RequestUrl := Config["server"].(map[string]interface{})["url"]
	RequestReader := bytes.NewReader(body)
	Request, _ := http.NewRequest("POST", RequestUrl.(string), RequestReader)

	client.AppendHeaders(Request, Config["rest"].(map[string]interface{})["headers"])

	Response, _ := client.HttpClient.Do( Request )
	defer Response.Body.Close()

	return nil
}
