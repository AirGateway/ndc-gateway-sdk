package ndc

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	// "os"
	// "sort"
	"strings"
	"time"

	// "github.com/matiasinsaurralde/mxj"
	"gopkg.in/yaml.v2"
)

type Message struct {
	XMLName xml.Name

	Client *Client `xml:"-"`

	SoapConfig SoapConfig

	Method string `xml:"-"`
	Params Params `xml:"-"`

	IsSoap bool `xml:"-"`

	XMLNS    string `xml:"xmlns,attr,omitempty"`
	XMLNSXSI string `xml:"xmlns:xsi,attr,omitempty"`

	EchoToken             string `xml:"EchoToken,attr,omitempty"`
	TimeStamp             string `xml:"TimeStamp,attr,omitempty"`
	Version               string `xml:"Version,attr,omitempty"`
	TransactionIdentifier string `xml:"TransactionIdentifier,attr,omitempty"`

	Body       string `xml:",innerxml"`
	ParamsBody string `xml:",innerxml"`
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"s:Envelope"`
	Body    SOAPBody
}

type SOAPBody struct {
	XMLName xml.Name `xml:"s:Body"`
	// Body string `xml:",innerxml"`
	Message *Message
}

type SoapConfig struct {
	RequestNamespace  string
	ResponseNamespace string
	EnvelopeTagName   string
	EnvelopeAttrs     xml.Attr
	BodyTagName       string
	BodyAttrs         xml.Attr
}

type Params []Param
type Param struct {
	Key, Value interface{}
}

func (message *Message) GetSoapConfig() (config SoapConfig) {

	config = SoapConfig{
		RequestNamespace:  "",
		ResponseNamespace: "",
		EnvelopeTagName:   "Envelope",
		BodyTagName:       "Body",
		EnvelopeAttrs:     xml.Attr{},
		BodyAttrs:         xml.Attr{},
	}

	var attributes yaml.MapSlice

	for _, v := range message.Client.Config["soap"] {
		switch v.Key {
		case "request_namespace":
			config.RequestNamespace = v.Value.(string)
		case "response_namespace":
			config.ResponseNamespace = v.Value.(string)
		case "attributes":
			attributes = v.Value.(yaml.MapSlice)
		}
	}

	for _, attr := range attributes {
		attrValue := attr.Value.(string)
		attrValue = strings.Replace(attrValue, "\"", "", -1)
		switch attr.Key {
		case "envelope":
			kv := make([]string, 2)
			kv = strings.Split(attrValue, "=")
			config.EnvelopeAttrs = xml.Attr{Name: xml.Name{"", kv[0]}, Value: kv[1]}
		case "body":
			kv := make([]string, 2)
			kv = strings.Split(attrValue, "=")
			config.BodyAttrs = xml.Attr{Name: xml.Name{"", kv[0]}, Value: kv[1]}
		}
	}

	config.EnvelopeTagName = strings.Join([]string{config.RequestNamespace, ":", config.EnvelopeTagName}, "")
	config.BodyTagName = strings.Join([]string{config.RequestNamespace, ":", config.BodyTagName}, "")

	return
}

func (message *Message) RenderNDCParams(enc *xml.Encoder, item interface{}, key string, index int, length int, parentElements []string) {
	t := fmt.Sprintf("%T", item)

	if parentElements == nil {
		parentElements = make([]string, 0)
	}

	// enc.EncodeToken(element)

	element := xml.StartElement{
		Name: xml.Name{"", key},
		Attr: []xml.Attr{},
	}

	enc.EncodeToken(element)

	switch t {
	case "ndc.Params":
		Items := item.(Params)

		ItemLength := len(Items)

		parentElements = append(parentElements, key)

		for k, v := range Items {

			Item := Items[k].Value
			ItemType := fmt.Sprintf("%T", Item)

			if ItemType == "ndc.Param" {
				Item := v.Value.(Param)

				message.RenderNDCParams(enc, Item.Value, Item.Key.(string), k, ItemLength, parentElements)
			} else {
				message.RenderNDCParams(enc, v.Value, v.Key.(string), k, ItemLength, parentElements)
			}

		}
	default:

		var data string

		switch t {
		case "float64":
			data = fmt.Sprintf("%.1f", item)
		case "int":
			data = fmt.Sprintf("%d", item)
		default:
			data = fmt.Sprintf("%s", item)
		}

		enc.EncodeToken(xml.CharData(data))
		enc.EncodeToken(element.End())

	}

	if index >= length-1 && len(parentElements) > 0 {
		reverseElements := make([]string, 0)
		for i := len(parentElements) - 1; len(reverseElements) != len(parentElements) && i >= 0; i-- {
			element := parentElements[i]
			if element != "" {
				reverseElements = append(reverseElements, element)
			}
		}

		for i := 0; i < len(reverseElements); i++ {
			var e = reverseElements[i]
			if e != "" {
				enc.EncodeToken(xml.EndElement{xml.Name{"", e}})
			}
			// parentElements[i] = ""
			reverseElements[i] = ""
		}
	}

	return
}

func (message *Message) RenderNDCWrapper(enc *xml.Encoder, buf *bytes.Buffer, item interface{}, key string, root bool, index int, length int, parentElements []string) {

	if message.IsSoap && root {

		soapEnvelope := xml.StartElement{
			Name: xml.Name{"", message.SoapConfig.EnvelopeTagName},
			Attr: []xml.Attr{message.SoapConfig.EnvelopeAttrs},
		}

		soapBody := xml.StartElement{
			Name: xml.Name{"", message.SoapConfig.BodyTagName},
			Attr: []xml.Attr{message.SoapConfig.BodyAttrs},
		}

		enc.EncodeToken(soapEnvelope)
		enc.EncodeToken(soapBody)
	}

	if root {
		item = item.(yaml.MapSlice)
		requestWrapper := xml.StartElement{
			Name: xml.Name{"", message.Method + "RQ"},
		}
		enc.EncodeToken(requestWrapper)
	}

	if parentElements == nil {
		parentElements = make([]string, 0)
	}

	t := fmt.Sprintf("%T", item)

	if t == "yaml.MapItem" || t == "yaml.MapSlice" {

		mapItem := item.(yaml.MapSlice)

		var kItemLen int

		for k, v := range mapItem {

			kItem := mapItem[k].Value
			kItemT := fmt.Sprintf("%T", kItem)

			if kItemT == "yaml.MapSlice" {
				kItemLen = len(kItem.(yaml.MapSlice))
				element := xml.StartElement{
					Name: xml.Name{"", mapItem[k].Key.(string)},
					Attr: []xml.Attr{},
				}
				parentElements = append(parentElements, mapItem[k].Key.(string))
				enc.EncodeToken(element)

			} else {
				kItemLen = length - 1
			}

			i := v.Value

			message.RenderNDCWrapper(enc, buf, i, v.Key.(string), false, k, kItemLen, parentElements)
		}

	} else {
		element := xml.StartElement{
			Name: xml.Name{"", key},
			Attr: []xml.Attr{},
		}

		var data string

		switch t {
		case "float64":
			data = fmt.Sprintf("%.1f", item)
		case "int":
			data = fmt.Sprintf("%d", item)
		case "ndc.Param":
			data = "ndc_param"
		case "ndc.Params":
			// data = "ndc_params"
			var paramsBuffer = new(bytes.Buffer)
			paramsEncoder := xml.NewEncoder(paramsBuffer)
			paramsEncoder.Indent(" ", "    ")

			message.RenderNDCParams(paramsEncoder, item, "", -1, -1, nil)
			paramsEncoder.Flush()
			enc.Flush()
			data = ""
			buf.Write([]byte("\n"))

			buf.Write(paramsBuffer.Bytes())

			// enc.Encode( paramsBuffer.String() )
		default:
			data = fmt.Sprintf("%s", item)
		}

		enc.EncodeToken(element)
		enc.EncodeToken(xml.CharData(data))
		enc.EncodeToken(element.End())

		if index >= length && len(parentElements) > 0 {
			reverseElements := make([]string, 0)

			for i := len(parentElements) - 1; len(reverseElements) != len(parentElements) && i >= 0; i-- {
				element := parentElements[i]
				if element != "" {
					reverseElements = append(reverseElements, element)
				}
			}

			for i := 0; i < len(reverseElements); i++ {
				var e = reverseElements[i]
				if e != "" {
					enc.EncodeToken(xml.EndElement{xml.Name{"", e}})
				}
				//parentElements[i] = ""
				reverseElements[i] = ""
			}
		}
	}
}
func (message *Message) Prepare() ([]byte, error) {

	// SOAP

	message.IsSoap = message.Client.Config["soap"] != nil

	if message.IsSoap {
		message.SoapConfig = message.GetSoapConfig()
	}

	// Namespace, etc.

	message.XMLName.Local = message.Method + "RQ"

	TimeStamp := time.Now().Format(time.RFC3339)
	EchoToken := sha1.New()
	EchoToken.Write([]byte(TimeStamp))

	message.XMLNS = "http://www.iata.org/IATA/EDIST"
	message.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"

	// Should we use? https://github.com/joeshaw/iso8601
	message.EchoToken = hex.EncodeToString(EchoToken.Sum(nil))
	message.TimeStamp = TimeStamp
	message.Version = "1.1.5"
	message.TransactionIdentifier = "TR-00000"

	// Template based body:

	ndc := message.Client.Config["ndc"]

	encBuffer := new(bytes.Buffer)

	encBuffer.Write([]byte(xml.Header))

	enc := xml.NewEncoder(encBuffer)
	enc.Indent(" ", "    ")
	message.RenderNDCWrapper(enc, encBuffer, ndc, "", true, -1, -1, nil)

	message.RenderNDCWrapper(enc, encBuffer, message.Params, "", false, -1, -1, nil)

	requestWrapperEnd := xml.EndElement{
		Name: xml.Name{"", message.Method + "RQ"},
	}

	enc.EncodeToken(requestWrapperEnd)

	if message.IsSoap {

		soapEnvelopeEnd := xml.EndElement{
			Name: xml.Name{"", message.SoapConfig.EnvelopeTagName},
		}

		soapBodyEnd := xml.EndElement{
			Name: xml.Name{"", message.SoapConfig.BodyTagName},
		}

		enc.EncodeToken(soapBodyEnd)
		enc.EncodeToken(soapEnvelopeEnd)
	}

	enc.Flush()

	// Final output

	var err error
	output := encBuffer.Bytes()

	return output, err
}
