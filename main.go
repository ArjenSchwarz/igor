package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/ArjenSchwarz/igor/slack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler handles incoming Lambda requests
func Handler(request events.APIGatewayProxyRequest) (slack.Response, error) {
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	response := handle(body{Body: request.Body})

	return response, nil
}

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
			response := handle(body{Body: v.Encode()})
			responseString, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(responseString)
		})
		http.ListenAndServe(":8080", nil)
	} else {
		lambda.Start(Handler)
	}
}

type body struct {
	Body string `json:"body"`
}
