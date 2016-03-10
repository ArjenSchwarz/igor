package plugins

import (
	"bytes"
	"errors"
	"strings"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

type HelpPlugin struct {
	name        string
	description string
}

func Help() IgorPlugin {
	plugin := HelpPlugin{
		name:        "help",
		description: "I provide help with the following commands",
	}
	return plugin
}

func (HelpPlugin) Work(request slack.SlackRequest) (slack.SlackResponse, error) {
	response := slack.SlackResponse{}
	message := strings.ToLower(request.Text)
	switch message {
	case "help":
		response = handleHelp(response)
	case "introduce yourself":
		response = handleIntroduction(response)
	case "tell me about yourself":
		response = handleTellMe(response)
	}
	if response.Text == "" {
		return response, errors.New("No match")
	}
	return response, nil
}

func (HelpPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["help"] = "This help message, providing a list of available commands"
	descriptions["introduce yourself"] = "A public introduction of Igor"
	descriptions["tell me about yourself"] = "Information about Igor"
	return descriptions
}

func handleHelp(response slack.SlackResponse) slack.SlackResponse {
	response.Text = "I can see that you're trying to find an Igor, would you like some help with that?"
	allPlugins := GetPlugins(config.ReadConfig())
	var buffer bytes.Buffer
	for _, plugin := range allPlugins {
		for command, description := range plugin.Describe() {
			buffer.WriteString("- *" + command + "*: " + description + "\n")
		}
		attach := slack.Attachment{}
		attach.Title = plugin.Description()
		attach.Text = buffer.String()
		attach.EnableMarkdownFor("text")
		response.AddAttachment(attach)
		buffer.Reset()
	}
	return response
}

func handleIntroduction(response slack.SlackResponse) slack.SlackResponse {
	response.Text = "I am Igor, reprethenting We-R-Igors."
	response.SetPublic()
	attach := slack.Attachment{}
	attach.Title = "A Spare Hand When Needed"
	attach.Text = "We come from Überwald, but are alwayth where we are needed motht.\n"
	attach.Text += "Run */igor help* to see which Igors are currently available."
	attach.EnableMarkdownFor("text")
	response.AddAttachment(attach)
	return response
}

func handleTellMe(response slack.SlackResponse) slack.SlackResponse {
	response.Text = "Originally Igors come from Überwald, but in this world our home is on GitHub."
	attach := slack.Attachment{}
	attach.Title = "GitHub"
	attach.Text = "Main repo on https://github.com/ArjenSchwarz/igor. Feel free to contribute"
	response.AddAttachment(attach)
	// TODO actually write the article mentioned below
	attach = slack.Attachment{}
	attach.Title = "ig.nore.me"
	attach.Text = "An introductory article about Igor and its creation can be found on https://ig.nore.me"
	response.AddAttachment(attach)
	return response
}

func (p HelpPlugin) Description() string {
	return p.description
}
func (p HelpPlugin) Name() string {
	return p.name
}
