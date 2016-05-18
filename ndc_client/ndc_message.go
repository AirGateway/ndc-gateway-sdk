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

  XMLNS string `xml:"xmlns,attr"`
  XMLNSXSI string `xml:"xmlns:xsi,attr"`

  EchoToken string  `xml:"EchoToken,attr"`
  TimeStamp string  `xml:"TimeStamp,attr"`
  Version string  `xml:"Version,attr"`
  TransactionIdentifier string  `xml:"TransactionIdentifier,attr"`

  Body string `xml:",innerxml"`
}

func( message *Message ) ToXml() ( []byte, error ) {

  // Namespace, etc.

  var XmlWriter = new( bytes.Buffer )

  message.XMLName.Local = message.Method
  message.XMLNS = "http://www.iata.org/IATA/EDIST"
  message.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"

  TimeStamp := time.Now().Format(time.RFC3339)
  EchoToken := sha1.New()
  EchoToken.Write( []byte(TimeStamp) )

  // Should we use? https://github.com/joeshaw/iso8601
  message.EchoToken = hex.EncodeToString( EchoToken.Sum(nil) )
  message.TimeStamp = TimeStamp
  message.Version = "1.1.5"
  message.TransactionIdentifier = "TR-00000"

  // Template based body:

  bodyMap := mxj.Map( message.Client.Config["ndc"].(map[string]interface{}) )
  bodyWriter := new(bytes.Buffer)
  bodyRawBytes, _ := bodyMap.XmlWriterRaw( bodyWriter, "_ndc_body" )

  bodyString := string(bodyRawBytes)
  bodyString = strings.Replace( bodyString, "<_ndc_body>", "", 1 )
  bodyString = strings.Replace( bodyString, "</_ndc_body>", "", 1 )

  message.Body = bodyString


  Map := mxj.Map(message.Params)

  _, err := Map.XmlWriterRaw(XmlWriter)

  output, err := xml.MarshalIndent( message, "  ", "    ")

  return output, err
}
