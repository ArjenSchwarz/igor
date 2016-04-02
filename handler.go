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
	config, err := config.GeneralConfig()
	if err != nil {
		response := slack.SomethingWrongResponse(request)
		response.Escape()
		return response
	}
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
	hasError := false
	forcePublic := false
	if request.Text[0] == '!' {
		forcePublic = true
		request.Text = request.Text[1:]
	}
	//TODO clean this up
	for _, value := range pluginlist {
		response, err := value.Work(request)
		if err == nil {
			if forcePublic {
				response.SetPublic()
			}
			return response
		}
		switch err.(type) {
		case *plugins.NoMatchError:
		default:
			// Something actually went wrong with one of the plugins,
			// return that something went wrong if nothing matches
			// Don't send the actual message though
			hasError = true
		}
	}
	if hasError {
		return slack.SomethingWrongResponse(request)
	}

	return slack.NothingFoundResponse(request)
}
