package ndc

import (
  "fmt"
  //"log"
	"bytes"
  "bufio"
  "strconv"
  "strings"
  "net/http"
	"io/ioutil"
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

type Extras struct {
	Value map[string]string
}
type ClientOptions struct {
	ConfigPath string
}

type Client struct {
	Options         ClientOptions
	HasTemplateVars bool
	Extras					map[string]Extras
	Config          map[string]yaml.MapSlice
	RawConfig       []byte
	HttpClient      *http.Client
}

type postProcess func(string)

func NewClient(options *ClientOptions, extras map[string]Extras) (*Client, error) {
	client := &Client{Options: *options}
	client.Config = map[string]yaml.MapSlice{}
	client.Extras = extras
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
func (client *Client) RequestAsynch(message Message, callback postProcess) {
	Response := client.Request(message)

	message_aux := ""
  reader := bufio.NewReader(Response.Body)
	for {
		line, err := reader.ReadBytes('\n')
		//fmt.Println(line)
		if err==nil {
			message_aux = message_aux + string(line)
			if strings.Contains(message_aux, "<!-- AG-EOM -->"){
				callback(message_aux)
				message_aux = ""
			}
		}else{fmt.Println("ERROR", err);break;}
	}
}
func (client *Client) RequestSynch(message Message) (string) {
  fmt.Println( "-> Doing Request:\n---\n" )
	Response := client.Request(message)
	fmt.Println( "-> Receiving response:\n---\n" )
	//fmt.Println( Response , "\n---\n-> Response body:\n---\n")
	body_, _ := ioutil.ReadAll(Response.Body)
	return string(body_);
}

func (client *Client) Request(message Message) (*http.Response) {

	var Config, ServerConfig, RestConfig map[string]interface{}
	//var Config, RestConfig map[string]interface{}

	message.Client = client

	body, _ := message.Prepare()

	if client.HasTemplateVars {
		Config = client.PrepareConfig(message)
	} else {
		 //Config = client.Config
	}

	//fmt.Println(Config)
	RestConfig 		 = Config["rest"].(map[string]interface{})
	ServerConfig 	 = Config["server"].(map[string]interface{})
  var RequestUrl  interface {}
  if env, ok := client.Extras["enviroment"]; ok==true{
    RequestUrl 		= ServerConfig["url_"+env.Value["url"]]
  }else{
    RequestUrl    = ServerConfig["url_production"]
  }
	RequestReader := bytes.NewReader(body)


	Request, _ 		:= http.NewRequest("POST", RequestUrl.(string), RequestReader)
	client.AppendHeaders(Request, RestConfig["headers"])
	//elem, ok := client.Extras["headers"]
	if headers, ok := client.Extras["headers"]; ok==true{
		for Header, Value := range headers.Value {
			//fmt.Println(Header, Value)
			Request.Header.Del(Header)
			Request.Header.Add(Header, Value)
		}
	}
	Response, _ 	:= client.HttpClient.Do(Request)

	return Response;

	//fmt.Println(Response.Body);

}
func convert( b []byte ) string {
    s := make([]string,len(b))
    for i := range b {
        s[i] = strconv.Itoa(int(b[i]))
    }
    return strings.Join(s,",")
}
