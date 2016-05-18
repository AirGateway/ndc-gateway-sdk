package main

import( "net/http"
        "fmt" )

func ndcHandler( w http.ResponseWriter, r *http.Request ) {
  fmt.Fprintf(w, "Response %s", r.URL.Path[1:])
}

func main() {
  http.Handle( "/", http.FileServer(http.Dir("./app") ) )
  http.HandleFunc( "/ndc", ndcHandler )
  http.ListenAndServe( ":8080", nil )
}
