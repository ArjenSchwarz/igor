package main

import (
	// "encoding/json"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
	// "log"
	// "strings"
)

// TODO validate the key
func handle(body *body) slack.SlackResponse {
	request := slack.LoadRequestFromQuery(body.Body)
	// config := ReadConfig()
	//

	response := determineResponse(request)
	return response
}

// Parse the responses from a list of plugin triggers
func determineResponse(request slack.SlackRequest) slack.SlackResponse {
	pluginlist := plugins.GetPlugins()
	//TODO clean this up
	for _, value := range pluginlist {
		response, err := value.Work(request)
		if err == nil {
			return response
		}
	}
	return slack.NothingFoundResponse(request)
}
