package main

import (
	// "encoding/json"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
	// "log"
	"net/url"
	// "strings"
)

func handle(body *body) *slack.SlackResponse {
	query, _ := url.ParseQuery(body.Body)
	message := query.Get("text")
	// config := ReadConfig()
	//

	response := determineResponse(message)
	return response
}

// Parse the responses from a list of plugin triggers
func determineResponse(message string) *slack.SlackResponse {
	// pluginmanager := plugins.GetPlugins
	// pluginlist := pluginmanager
	plugin := plugins.Help()
	response, _ := plugin.Response(message)
	return response
}
