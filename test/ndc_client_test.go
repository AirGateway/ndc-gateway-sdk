package tests

import( "testing"
        "fmt"
        "ndc-go-sdk/ndc_client" )

var (
  options ndc.ClientOptions
  client *ndc.Client
)

func init() {
  options = ndc.ClientOptions{
    Endpoint: "http://127.0.0.1:8000/",
    ConfigPath: "../config/ndc-openndc.yml",
  }
}

func TestConfigLoad( t *testing.T ) {
  fmt.Println( "TestConfigLoad", options )

  var config map[interface{}]interface{}

  err := ndc.LoadConfig( options.ConfigPath, &config )

  fmt.Println(len(config))
  if err != nil {
    t.Fatalf("ndc.LoadConfig returned error: %v", err )
  }

  if len(config) == 0 {
    t.Errorf("Configuration map structure is empty!")
  }
}

func TestClientInit( t *testing.T ) {
  client, err := ndc.NewClient( &options )
  fmt.Println( "TestClientInit", client )

  if err != nil {
    t.Fatalf( "ndc.NewClient returned error: %v", err )
  }
}
