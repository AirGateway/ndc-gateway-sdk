package tests

import( "testing"
        "ndc-go-sdk/ndc_client" )

var (
  options ndc.ClientOptions
  client *ndc.Client
)

func init() {
  options = ndc.ClientOptions{
    //Endpoint: "http://127.0.0.1:8000/",
    Endpoint: "http://prxy.airgateway.net:8080/gtwy/ndc/",
    ConfigPath: "../config/ndc-openndc.yml",
  }
}

func TestConfigLoad( t *testing.T ) {
  t.Log( "TestConfigLoad", options )

  var config map[interface{}]interface{}

  err := ndc.LoadConfig( options.ConfigPath, &config )

  if err != nil {
    t.Fatalf("ndc.LoadConfig returned error: %v", err )
  }

  if len(config) == 0 {
    t.Errorf("Configuration map structure is empty!")
  }
}

func TestClientInit( t *testing.T ) {
  client, err := ndc.NewClient( &options )
  t.Log( "TestClientInit", client )

  if err != nil {
    t.Fatalf( "ndc.NewClient returned error: %v", err )
  }
}
