package plugins

import (
	"bytes"
	"strings"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// HelpPlugin provides help functions
type HelpPlugin struct {
	name        string
	description string
	request     slack.Request
	config      helpConfig
}

// Config returns the plugin configuration
func (plugin HelpPlugin) Config() IgorConfig {
	return plugin.config
}

type helpConfig struct {
	languages      map[string]config.LanguagePluginDetails
	chosenLanguage string
}

// Languages returns the languages available for the plugin
func (config helpConfig) Languages() map[string]config.LanguagePluginDetails {
	return config.languages
}

// ChosenLanguage returns the language active for this plugin
func (config helpConfig) ChosenLanguage() string {
	return config.chosenLanguage
}

// Help instantiates the HelpPlugin
func Help(request slack.Request) IgorPlugin {
	pluginName := "help"
	pluginConfig := helpConfig{
		languages: getPluginLanguages(pluginName),
	}
	plugin := HelpPlugin{
		name:    pluginName,
		request: request,
		config:  pluginConfig,
	}

	return plugin
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
//
//  * help
//  * introduce yourself
//  * tell me about yourself
func (plugin HelpPlugin) Work() (slack.Response, error) {
	response := slack.Response{}
	message, language := getCommandName(plugin)
	plugin.config.chosenLanguage = language
	switch message {
	case "help":
		tmpresponse, err := plugin.handleHelp(response)
		if err != nil {
			return tmpresponse, err
		}
		response = tmpresponse
	case "intro":
		response = plugin.handleIntroduction(response)
	case "tellme":
		response = plugin.handleTellMe(response)
	case "whoami":
		response = plugin.handleWhoAmI(response)
	}
	if response.Text == "" {
		return response, CreateNoMatchError("Nothing found")
	}
	return response, nil
}

// Describe provides the triggers HelpPlugin can handle
func (plugin HelpPlugin) Describe(language string) map[string]string {
	descriptions := make(map[string]string)

	for _, values := range getAllCommands(plugin, language) {
		descriptions[values.Command] = values.Description
	}
	return descriptions
}

func (plugin HelpPlugin) handleHelp(response slack.Response) (slack.Response, error) {
	commandDetails := getCommandDetails(plugin, "help")
	response.Text = commandDetails.Texts["response_text"]
	config, err := config.GeneralConfig()
	if err != nil {
		return response, err
	}
	allPlugins := GetPlugins(plugin.request, config)
	c := make(chan slack.Attachment)
	for _, igor := range allPlugins {
		go func(igor IgorPlugin, language string) {
			var buffer bytes.Buffer
			for command, description := range igor.Describe(language) {
				buffer.WriteString("- *" + command + "*: " + description + "\n")
			}
			attach := slack.Attachment{}
			attach.Title = igor.Description(language)
			attach.Text = buffer.String()
			attach.EnableMarkdownFor("text")
			c <- attach
		}(igor, plugin.config.chosenLanguage)
	}
	for i := 0; i < len(allPlugins); i++ {
		response.AddAttachment(<-c)
	}
	return response, nil
}

func (plugin HelpPlugin) handleIntroduction(response slack.Response) slack.Response {
	commandDetails := getCommandDetails(plugin, "intro")
	response.Text = commandDetails.Texts["response_text"]
	response.SetPublic()
	attach := slack.Attachment{}
	attach.Title = commandDetails.Texts["attach_title"]
	attach.Text = commandDetails.Texts["attach_text"]
	attach.EnableMarkdownFor("text")
	response.AddAttachment(attach)
	return response
}

func (plugin HelpPlugin) handleTellMe(response slack.Response) slack.Response {
	commandDetails := getCommandDetails(plugin, "tellme")
	response.Text = commandDetails.Texts["response_text"]
	attach := slack.Attachment{}
	attach.Title = "GitHub"
	attach.Text = commandDetails.Texts["github_text"]
	response.AddAttachment(attach)
	attach = slack.Attachment{}
	attach.Title = "ig.nore.me"
	attach.Text = commandDetails.Texts["site_text"]
	response.AddAttachment(attach)
	return response
}

func (plugin HelpPlugin) handleWhoAmI(response slack.Response) slack.Response {
	commandDetails := getCommandDetails(plugin, "whoami")
	request := plugin.request
	response.Text = commandDetails.Texts["response_text"]
	attach := slack.Attachment{}
	attach.Title = commandDetails.Texts["attach_title"]
	attach.AddField(slack.Field{Title: "Name", Value: request.UserName, Short: true})
	attach.AddField(slack.Field{Title: "UserID", Value: request.UserID, Short: true})
	attach.AddField(slack.Field{Title: "Channel", Value: request.ChannelName, Short: true})
	attach.AddField(slack.Field{Title: "ChannelID", Value: request.ChannelID, Short: true})
	attach.AddField(slack.Field{Title: "Team", Value: request.TeamDomain, Short: true})
	attach.AddField(slack.Field{Title: "TeamID", Value: request.TeamID, Short: true})
	response.AddAttachment(attach)
	return response
}

// Description returns a global description of the plugin
func (plugin HelpPlugin) Description(language string) string {
	return getDescriptionText(plugin, language)
}

// Name returns the name of the plugin
func (plugin HelpPlugin) Name() string {
	return plugin.name
}

// Message returns a formatted version of the original message
func (plugin HelpPlugin) Message() string {
	return strings.ToLower(plugin.request.Text)
}
