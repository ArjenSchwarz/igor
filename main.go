package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	body := body{Body: request.Body}
	// json.Unmarshal([]byte(request.Body), &body)

	log.Printf("Original body: %s\n", request.Body)

	log.Printf("Parsed body: %s\n", body)

	response := handle(body)

	responseArray, _ := json.Marshal(response)
	// n := bytes.IndexByte(responseArray, 0)
	responseString := string(responseArray)
	return events.APIGatewayProxyResponse{
		Body:       responseString,
		StatusCode: 200,
	}, nil
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
			body := body{Body: v.Encode()}
			response := handle(body)
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
