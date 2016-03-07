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
	if strings.Compare(message, "help") == 0 {
		response = handleHelp(message, response)
	}
	if strings.Compare(message, "introduce yourself") == 0 {
		response = handleIntroduction(message, response)
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
	return descriptions
}

func handleHelp(message string, response slack.SlackResponse) slack.SlackResponse {
	response.Text = "I can see that you're trying to find an Igor, would you like some help with that?"
	allPlugins := GetPlugins()
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

func handleIntroduction(message string, response slack.SlackResponse) slack.SlackResponse {
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

func (p HelpPlugin) Description() string {
	return p.description
}
func (p HelpPlugin) Name() string {
	return p.name
}
