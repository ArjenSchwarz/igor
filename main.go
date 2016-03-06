package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type message struct {
	Value string `json:"value"`
}

func main() {
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

type body struct {
	Body string `json:"body"`
}
