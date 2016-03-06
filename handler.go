package main

import (
	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
)

func handle(body body) slack.SlackResponse {
	request := slack.LoadRequestFromQuery(body.Body)
	config := config.ReadConfig()
	response := slack.SlackResponse{}
	if !request.Validate(config) {
		response = slack.ValidationErrorResponse()
	} else {
		response = determineResponse(request)
	}
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
