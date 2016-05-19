package ndc

import( "fmt"
  "bytes"
  "io/ioutil"
	"github.com/matiasinsaurralde/yaml"
  "github.com/clbanning/mxj"
  "reflect"
  //"encoding/xml"
)

type UnknownMapString map[string]interface{}

func ConvertStuff( m interface{}) (o map[string]interface{}) {
  o = make( map[string]interface{})
  // s := reflect.ValueOf( &m ).Elem()
  // t := reflect.TypeOf( &m )
  fmt.Println("ConvertStuff")

  // fmt.Println(t.Field()())
  // o["xd"] = "xdvalue"

  // for k, _ := range m {
      // fmt.Println(1,k,v)
      // val :=
      // key := k.(string)
      //o[k] = "dummy"
  // }

  // t := reflect.TypeOf( m )
  // fmt.Println(t.NumField())
  fmt.Println()
  fmt.Println("original",m )
  fmt.Println("--")
  v := reflect.ValueOf(m)

  fmt.Println("v",v.MapKeys() )
  fmt.Println(v)
  for _, key := range v.MapKeys() {
    // fmt.Println(1,key)
    // fmt.Println( v.MapIndex(key) )
    // fmt.Println( "*", key.Elem().String( ))
    keyName := key.Elem().String()
    // values := v.MapIndex(key)
    // val := v.FieldByName(keyName).(interface{})
    // fmt.Println(val)

    // fmt.Println( keyName, key.Type(), v.MapIndex(key) )
    fmt.Println( keyName )

    // var mapmap interface{} = v.MapIndex(key)
    var mapmap interface{} = v.MapIndex(key)
    // fmt.Println( "mapmap", mapmap.Type(), mapmap.(type) )
    switch t := mapmap.(type) {
    default:
      fmt.Printf("*** default %T\n", t)
    case *string:
      fmt.Println("*** string")
    }
    // ConvertStuff(mapmap)
    // var x = v.MapIndex(key).Elem()
    // fmt.Println( "!!!", reflect.TypeOf(x).Kind() )
    // tt := reflect.TypeOf( mapmap )
    // fmt.Println( "!!! ", tt, reflect.ValueOf(v.MapIndex(key)).Kind() )
    // vv := reflect.ValueOf(mapmap)
    // fmt.Println( "vv",  )
    // fmt.Println( "**", values.NumField() )
    // o[keyName] = "abc"
    // o[key.Name()] = key.Name()
  }

  // o["xd"] = "xda"

  return
}

func YamlStuff() {
  config := make(map[string]interface{})
  // err := LoadConfig( client.Options.ConfigPath, &client.Config )
  RawConfig, _ := ioutil.ReadFile( "config/ndc-openndc.yml" )

  yaml.Unmarshal( RawConfig, config )

  // fmt.Println(config)

  // ndcConfig := cleanupInterfaceMap(config["ndc"])

  // fmt.Println( string(RawConfig) )

  // ndcConfig := config["ndc"]

  // fmt.Println( ndcConfig )

  // ndcConfigTwo := ndcConfig.(UnknownMapString)
  // var val2 UnknownMapString  = ndcConfig.(map[string]interface{})
  // fmt.Println(val2)


  someOtherMap := make(map[string]interface{})

  subMap := make(map[interface{}]interface{})
  subMap["subkey"] = "subvalue"

  convertedSubMap := ConvertStuff(subMap)

  fmt.Println( "subMap", subMap)
  fmt.Println( "convertedSubMap", convertedSubMap)

  someOtherMap["key"] = "value"
  someOtherMap["superkey"] = convertedSubMap


  var XmlWriter = new( bytes.Buffer )

  Map := mxj.Map( config["ndc"].(map[string]interface{}) )
  // Map := mxj.Map( ConvertStuff(config["ndc"]) )
  Map.XmlWriterRaw(XmlWriter)
  // output, err := xml.MarshalIndent( message, "  ", "    ")

  fmt.Println( XmlWriter )


}
