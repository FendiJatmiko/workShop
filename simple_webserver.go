package main
import (
  "fmt"
  "net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

  fmt.Fprintln(w, "<h1>#############</h1>")
  fmt.Fprintln(w, "<p>this is paragraph</p>")
  fmt.Fprintln(w, "<h2>Hello, World!</h2>")
  fmt.Fprintln(w, "<h1>#############</h1>")
}

func main() {
  http.HandleFunc("/", RootHandler)
  http.ListenAndServe(":8000", nil)
}
