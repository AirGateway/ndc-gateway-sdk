package ndc

import(
  "encoding/xml"
  "bytes"
  "time"
  "crypto/sha1"
  "encoding/hex"
  "strings"

  "github.com/clbanning/mxj"
)

type Message struct {
  XMLName xml.Name

  Client *Client `xml:"-"`

  Method string `xml:"-"`
  Params map[string]interface{} `xml:"-"`

  IsSoap  bool `xml:"-"`

  XMLNS string `xml:"xmlns,attr,omitempty"`
  XMLNSXSI string `xml:"xmlns:xsi,attr,omitempty"`

  EchoToken string  `xml:"EchoToken,attr,omitempty"`
  TimeStamp string  `xml:"TimeStamp,attr,omitempty"`
  Version string  `xml:"Version,attr,omitempty"`
  TransactionIdentifier string  `xml:"TransactionIdentifier,attr,omitempty"`

  Body string `xml:",innerxml"`
  ParamsBody string `xml:",innerxml"`
}

func( message *Message ) Prepare() ( []byte, error ) {

  message.IsSoap = message.Client.Config["soap"] != nil

  // Namespace, etc.

  message.XMLName.Local = message.Method

  if message.IsSoap {

  } else {
    TimeStamp := time.Now().Format(time.RFC3339)
    EchoToken := sha1.New()
    EchoToken.Write( []byte(TimeStamp) )

    message.XMLNS = "http://www.iata.org/IATA/EDIST"
    message.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"

    // Should we use? https://github.com/joeshaw/iso8601
    message.EchoToken = hex.EncodeToString( EchoToken.Sum(nil) )
    message.TimeStamp = TimeStamp
    message.Version = "1.1.5"
    message.TransactionIdentifier = "TR-00000"
  }

  // Template based body:

  bodyMap := mxj.Map( message.Client.Config["ndc"].(map[string]interface{}) )
  bodyWriter := new(bytes.Buffer)
  bodyRawBytes, _ := bodyMap.XmlWriterRaw( bodyWriter, "_ndc_body" )

  bodyString := string(bodyRawBytes)
  bodyString = strings.Replace( bodyString, "<_ndc_body>", "", 1 )
  bodyString = strings.Replace( bodyString, "</_ndc_body>", "", 1 )

  message.Body = bodyString

  // Params:

  paramsWriter := new( bytes.Buffer )

  paramsMap := mxj.Map(message.Params)

  paramsString, err := paramsMap.XmlWriterRaw(paramsWriter)

  message.ParamsBody = string(paramsString)

  output, err := xml.MarshalIndent( message, "  ", "    ")

  return output, err
}
