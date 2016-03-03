package plugins

import (
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

func (HelpPlugin) Response(message string) (*slack.SlackResponse, error) {
	response := new(slack.SlackResponse)
	if strings.Compare(message, "help") == 0 {
		response.Text = "I can see that you're trying to find an igor, would you like some help with that?"
	}
	if strings.Compare(message, "explain yourself") == 0 {
		response.Text = "Igors are useful servants. I can do many things once they're written as a plugin."
		response.ResponseType = "in_channel"
		// attachments := new(Attachment[])
		var attachments []*slack.Attachment
		attach := new(slack.Attachment)
		attach.Text = "We live to serve"
		attach.PreText = "I have my grandfather's hands"
		attach.Title = "Many wonders have we built"
		attachments = append(attachments, attach)
		response.Attachments = attachments
	}
	if strings.Compare(message, "who are you?") == 0 {
		response.Text = "I am a Slack slash command, written in Go, and running on Lambda."
		response.ResponseType = "in_channel"
	} else {
		return nil, errors.New("No match")
	}
	return response, nil
}

func (HelpPlugin) HelpMessages() []string {
	var messages []string
	messages = append(messages, "help")
	messages = append(messages, "explain yourself")
	messages = append(messages, "who are you?")
	return messages
}
