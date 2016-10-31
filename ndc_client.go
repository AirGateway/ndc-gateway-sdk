package ndc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	"gopkg.in/yaml.v2"
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
	ConfigPath string
}

type Client struct {
	Options         ClientOptions
	HasTemplateVars bool
	Config          map[string]yaml.MapSlice
	RawConfig       []byte
	HttpClient      *http.Client
}

func NewClient(options *ClientOptions) (*Client, error) {
	client := &Client{Options: *options}
	client.Config = map[string]yaml.MapSlice{}
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

func MapSliceToMap(slice yaml.MapSlice, m map[string]interface{}) map[string]interface{} {
	if m == nil {
		m = make(map[string]interface{})
	}
	for _, entry := range slice {
		switch entry.Value.(type) {
		case yaml.MapSlice:
			m[fmt.Sprint(entry.Key)] = MapSliceToMap(entry.Value.(yaml.MapSlice), nil)
		default:
			m[fmt.Sprint(entry.Key)] = fmt.Sprint(entry.Value)
		}
	}

	return m
}

func (client *Client) LoadConfig() error {
	config, err := ioutil.ReadFile(client.Options.ConfigPath)

	client.RawConfig = config
	err = yaml.Unmarshal(client.RawConfig, &client.Config)

	if ConfigHasTemplateVars(&client.RawConfig) > 0 {
		client.HasTemplateVars = true
	}

	return err
}

func (client *Client) PrepareConfig(message Message) (Config map[string]interface{}) {

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

	ConfigMapSlice := yaml.MapSlice{}
	yaml.Unmarshal([]byte(ModifiedConfig), &ConfigMapSlice)

	Config = MapSliceToMap(ConfigMapSlice, nil)

	return
}

func (client *Client) AppendHeaders(r *http.Request, HeadersConfig interface{}) {
	headers := HeadersConfig.(map[string]interface{})
	for Header, Value := range headers {
		r.Header.Add(Header, Value.(string))
	}
}

func (client *Client) Request(message Message) *http.Response {

	var Config, ServerConfig, RestConfig map[string]interface{}
	//var Config, RestConfig map[string]interface{}

	message.Client = client

	body, _ := message.Prepare()

	//fmt.Println(body)

	if client.HasTemplateVars {
		Config = client.PrepareConfig(message)
	} else {
		 //Config = client.Config
	}

	//fmt.Println(Config)
	RestConfig = Config["rest"].(map[string]interface{})
	ServerConfig = Config["server"].(map[string]interface{})
	RequestUrl := ServerConfig["url"]
	//RequestUrl := RestConfig["url"]
	RequestReader := bytes.NewReader(body)
	Request, _ := http.NewRequest("POST", RequestUrl.(string), RequestReader)

	/*fmt.Println(RequestReader)
	fmt.Println("")
	fmt.Println("")

	buf := new(bytes.Buffer)
	buf.ReadFrom(RequestReader)
	s := buf.String()
	fmt.Println(s)
	fmt.Println("")
	fmt.Println("")*/


	client.AppendHeaders(Request, RestConfig["headers"])

	Response, _ := client.HttpClient.Do(Request)

	//fmt.Println(Response)

	return Response
}
func convert( b []byte ) string {
    s := make([]string,len(b))
    for i := range b {
        s[i] = strconv.Itoa(int(b[i]))
    }
    return strings.Join(s,",")
}
