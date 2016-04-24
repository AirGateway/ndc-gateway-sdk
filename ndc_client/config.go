package ndc

import (
//        "fmt"
//        "gopkg.in/yaml.v2"
)

var data = `
label: ONDC

rest:
  url: http://ndc-sandbox.dev/api/v0/ndc/
  headers:
    Accept: application/xml

    Content-Type: application/xml

soap:
auth:

ndc:
  Document:
    Name: NDC Wrapper
    ReferenceVersion: "1.0"
  Party:
    Sender:
      ORA_Sender:
        AirlineID: FA
        Name: Fake Air
        AgentUser:
          Name: TravelWadus
          Type: TravelManagementCompany
          PseudoCity: A4A
          AgentUserID: travelwadus
          IATA_Number: "00002015"
  Participants:
    Participant:
      AggregatorParticipant:
        Name: Wadus NDC Gateway
        AggregatorID: WAD-00000
  Parameters:
    CurrCodes:
      CurrCode: EUR
  Preference:
    AirlinePreferences:
      Airline:
        AirlineID: FA
    FarePreferences:
      FarePreferences:
        Types:
          Type:
            Code: '759'
    CabinPreferences:
      CabinType:
        Code: M
        Definition: Economy/coach discounted`

/*
func main() {
  config := Config{}
  err := yaml.Unmarshal([]byte(data), &config)
  fmt.Println( err, data, config )
}
*/
