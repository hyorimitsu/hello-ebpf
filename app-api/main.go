package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func send(url string) string {
	req, _ := http.NewRequest("GET", url, nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray)
}

func handler(w http.ResponseWriter, r *http.Request) {
	url, ok := r.URL.Query()["url"]
	if !ok || len(url[0]) < 1 {
		fmt.Fprintf(w, "Hello, eBPF!")
		return
	}

	resp := send(url[0])
	fmt.Fprintf(w, resp)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
