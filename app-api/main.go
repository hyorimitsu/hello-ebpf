package main

import (
	"fmt"
	"net/http"
)

func handler(str string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, str)
	}
}

func main() {
	http.HandleFunc("/", handler("Hello, World!"))
	http.ListenAndServe(":8080", nil)
}
