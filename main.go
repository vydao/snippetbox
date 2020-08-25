package  main

import (
    "log"
    "net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World!"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Show Snippet"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Create Snippet"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
    mux.HandleFunc("/snippet/create", createSnippet)

    log.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", mux)
    log.Fatal(err)
}