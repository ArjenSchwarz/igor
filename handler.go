package main

import (
	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
)

// handle is the main handling function. It parses the received message and
// ensures that a response is collected.
// It also ensures that the resulting response is properly escaped
func handle(body body) slack.Response {
	request := slack.LoadRequestFromQuery(body.Body)
	config := config.GeneralConfig()
	response := slack.Response{}
	if !request.Validate(config) {
		response = slack.ValidationErrorResponse()
	} else {
		response = determineResponse(request, config)
	}
	response.Escape()
	return response
}

// determineResponse parses the responses from a list of plugin triggers
func determineResponse(request slack.Request, config config.Config) slack.Response {
	pluginlist := plugins.GetPlugins(config)
	//TODO clean this up
	for _, value := range pluginlist {
		response, err := value.Work(request)
		// TODO differentiate between not found and something went wrong errors
		if err == nil {
			return response
		}
	}
	return slack.NothingFoundResponse(request)
}
