package main

import (
	// "encoding/json"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
	"log"
	"net/url"
	// "strings"
)

// TODO validate the key
func handle(body *body) slack.SlackResponse {
	query, _ := url.ParseQuery(body.Body)
	message := query.Get("text")
	// config := ReadConfig()
	//

	response := determineResponse(message)
	return response
}

// Parse the responses from a list of plugin triggers
func determineResponse(message string) slack.SlackResponse {
	pluginlist := plugins.GetPlugins()
	//TODO clean this up
	for _, value := range pluginlist {
		response, err := value.Response(message)
		if err != nil {
			log.Println("Something went wrong")
		}
		return response
	}
	// TODO return "nothing ofound" result if nothing is there
	return slack.SlackResponse{}
}
