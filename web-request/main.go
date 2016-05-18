package main

import( "net/http"
        "fmt"
        "ndc-go-sdk/ndc_client")

var options ndc.ClientOptions
var client ndc.Client

func ndcHandler( w http.ResponseWriter, r *http.Request ) {
  fmt.Fprintf(w, "Response %s", r.URL.Path[1:])

  fmt.Println( "ndcHandler" )

  options := ndc.ClientOptions{
    Endpoint: "http://127.0.0.1:8000/",
    ConfigPath: "config/ndc-openndc.yml",
  }

  client, _ := ndc.NewClient( &options )

  params := map[string]interface{}{
    "CoreQuery": map[string]interface{}{
      "OriginDestinations": map[string]interface{}{
        "OriginDestination": map[string]interface{} {
          "Departure": map[string]interface{} {
            "AirportCode": "MUC",
            "Date": "2016-04-01",
          },
          "Arrival": map[string]interface{} {
            "AirportCode": "LHR",
            "Date": "2016-04-10",
          },
        },
      },
    },
  }

  client.Request(ndc.Message{
    Method: "AirShoppingRQ",
    Params: params,
  })
}

func main() {
  http.Handle( "/", http.FileServer(http.Dir("./app") ) )
  http.HandleFunc( "/ndc", ndcHandler )
  http.ListenAndServe( ":8080", nil )
}
