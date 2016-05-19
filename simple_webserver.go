package main
import (
  "fmt"
  "net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "<p>this is paragraph</p>")
  fmt.Fprintln(w, "<h2>Hello, World!</h2>")
}

func nameHandler(w http.ResponseWriter, r *http.Request) {
  hallo := r.FormValue("hallo")
  nama  := r.FormValue("nama")
  fmt.Fprintf(w, "<h1>%s, %s!</h1>", hallo, nama)
}

func main() {
  http.HandleFunc("/", RootHandler)
  http.HandleFunc("/nama", nameHandler)
  http.ListenAndServe(":8000", nil)
}
