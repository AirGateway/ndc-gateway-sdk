package ndc

import(
  "encoding/xml"
)

type Message struct {
  XMLName xml.Name
  Method string `xml:"-"`

  EchoToken string  `xml:"EchoToken,attr"`
  TimeStamp string  `xml:"TimeStamp,attr"`
  Version string  `xml:"Version,attr"`
  TransactionIdentifier string  `xml:"TransactionIdentifier,attr"`

  XMLNS string `xml:"xmlns,attr"`
  XMLNSXSI string `xml:"xmlns:xsi,attr"`
}

func( message *Message ) ToXml() []byte {

  message.XMLName.Local = message.Method
  message.XMLNS = "http://www.iata.org/IATA/EDIST"
  message.XMLNSXSI = "http://www.w3.org/2001/XMLSchema-instance"

  output, _ := xml.MarshalIndent( message, "  ", "    ")

  return output
}
