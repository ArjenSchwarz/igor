package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type message struct {
	Value string `json:"value"`
}

var servervar bool

func init() {
	flag.BoolVar(&servervar, "server", false, "Run Igor as a server")
	flag.Parse()
}

func main() {
	if servervar {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			r.ParseForm()
			v := url.Values{}
			for field, values := range r.Form {
				for _, value := range values {
					v.Add(field, value)
				}
			}
			body := body{Body: v.Encode()}
			response := handle(body)
			responseString, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(responseString)
		})
		http.ListenAndServe(":8080", nil)
	} else {
		buf := new(bytes.Buffer)
		args := os.Args
		event := []byte(args[1])

		body := body{}
		json.Unmarshal(event, &body)

		response := handle(body)

		responseString, _ := json.Marshal(response)
		fmt.Fprintf(buf, "%s", responseString)
		buf.WriteTo(os.Stdout)
	}
}

type body struct {
	Body string `json:"body"`
}
