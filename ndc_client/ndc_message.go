package ndc

import(
  "encoding/xml"
  "github.com/clbanning/mxj"
  "bytes"
  "time"
)

type Message struct {
  XMLName xml.Name

  Method string `xml:"-"`
  Params map[string]interface{} `xml:"-"`

  EchoToken string  `xml:"EchoToken,attr"`
  TimeStamp string  `xml:"TimeStamp,attr"`
  Version string  `xml:"Version,attr"`
  TransactionIdentifier string  `xml:"TransactionIdentifier,attr"`

  XMLNS string `xml:"xmlns,attr"`
  XMLNSXSI string `xml:"xmlns:xsi,attr"`
}

func( message *Message ) ToXml() ( []byte, error ) {

  var XmlWriter = new( bytes.Buffer )

  message.XMLName.Local = message.Method
  message.XMLNS = "http://www.iata.org/IATA/EDIST"
  message.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"

  // Should we use? https://github.com/joeshaw/iso8601
  message.TimeStamp = time.Now().Format(time.RFC3339)

  Map := mxj.Map(message.Params)

  _, err := Map.XmlWriterRaw(XmlWriter)

  output, err := xml.MarshalIndent( message, "  ", "    ")

  return output, err
}
