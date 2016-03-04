package plugins

import (
	"bytes"
	"errors"
	"strings"

	// "github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

type HelpPlugin struct {
	name        string
	version     string
	author      string
	description string
}

func Help() IgorPlugin {
	plugin := HelpPlugin{
		name:        "Help",
		version:     "0.1",
		author:      "Arjen Schwarz",
		description: "Basic help functionalities",
	}
	return plugin
}

func (HelpPlugin) Response(message string) (slack.SlackResponse, error) {
	response := slack.SlackResponse{}
	message = strings.ToLower(message)
	if strings.Compare(message, "help") == 0 {
		response = handleHelp(message, response)
	}
	if strings.Compare(message, "explain yourself") == 0 {
		response = handleExplain(message, response)
	}
	if strings.Compare(message, "who are you?") == 0 {
		response = handleWhoAreYou(message, response)
	}
	if response.Text == "" {
		return response, errors.New("No match")
	}
	return response, nil
}

func (HelpPlugin) Descriptions() map[string]string {
	descriptions := make(map[string]string)
	descriptions["help"] = "This help message, proving a list of available commands"
	descriptions["who are you?"] = "A public introduction of Igor"
	descriptions["explain yourself"] = "A public explanation of Igor"
	return descriptions
}

func handleHelp(message string, response slack.SlackResponse) slack.SlackResponse {
	response.Text = "I can see that you're trying to find an igor, would you like some help with that?"
	allPlugins := GetPlugins()
	var buffer bytes.Buffer
	for _, value := range allPlugins {
		for command, description := range value.Descriptions() {
			buffer.WriteString("- *" + command + "*: " + description + "\n")
		}
	}
	attach := slack.Attachment{}
	attach.Title = "Available igors"
	attach.Text = buffer.String()
	attach.EnableMarkdownFor("text")
	response.AddAttachment(attach)
	return response
}

func handleExplain(message string, response slack.SlackResponse) slack.SlackResponse {
	response.Text = "Igors are useful servants. We are legion and can do many things."
	response.SetPublic()
	attach := slack.Attachment{}
	attach.Text = "/igor help will show you everything I can currently do."
	response.AddAttachment(attach)
	return response
}

func handleWhoAreYou(message string, response slack.SlackResponse) slack.SlackResponse {
	response.Text = "I am a Slack slash command, written in Go, and running on Lambda."
	response.SetPublic()
	return response
}
