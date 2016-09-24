package plugins

import (
	"regexp"
	"strings"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// IgorPlugin is the interface that needs to be followed by all plugins
type IgorPlugin interface {
	Work() (slack.Response, error)
	Describe(string) map[string]string
	Name() string
	Description(string) string
	Message() string
	Config() IgorConfig
}

// IgorConfig is the interface for all plugin Configuration
type IgorConfig interface {
	Languages() map[string]config.LanguagePluginDetails
	ChosenLanguage() string
}

// GetPlugins retrieves all the plugins that are activated. It checks the
// config for a whitelist and blacklist as well.
func GetPlugins(request slack.Request, config config.Config) map[string]IgorPlugin {
	plugins := make(map[string]IgorPlugin)
	plugins["help"] = Help(request)
	//TODO should handle these errors somehow. Returning an error when the
	//plugin isn't called doesn't make a lot of sense though
	plugins["weather"], _ = Weather(request)
	plugins["tumblr"], _ = RandomTumblr(request)
	plugins["status"], _ = Status(request)
	plugins["xkcd"], _ = Xkcd(request)

	// Whitelist plugins
	if config.Whitelist != nil {
		whitelist := make(map[string]IgorPlugin)
		whitelist["help"] = Help(request) //Help is always required
		for _, allowedPlugin := range config.Whitelist {
			whitelist[allowedPlugin] = plugins[allowedPlugin]
		}
		plugins = whitelist
	}

	// Blacklist plugins
	if config.Blacklist != nil {
		for _, pluginname := range config.Blacklist {
			if pluginname != "help" { // Help is always required
				delete(plugins, pluginname)
			}
		}
	}
	return plugins
}

// NoMatchError is an error type to indicate a plugin didn't find a match
type NoMatchError struct {
	Message string
}

// Error returns a string interpretation of the NoMatchError
func (e *NoMatchError) Error() string {
	return "No match found:" + e.Message
}

// CreateNoMatchError creates a new NoMatchError instance
func CreateNoMatchError(message string) *NoMatchError {
	return &NoMatchError{Message: message}
}

func getCommandName(plugin IgorPlugin) (string, string) {
	// It's possible for a command to have substitutions
	// Therefore, this needs to be taken into account
	reMain := regexp.MustCompile("(.*) \\[")
	reCommand := regexp.MustCompile("^([^ ]*) ?")
	subCommandArray := reCommand.FindStringSubmatch(plugin.Message())
	subCommand := ""
	if subCommandArray != nil {
		subCommand = strings.ToLower(subCommandArray[1])
	}
	for language, details := range plugin.Config().Languages() {
		for name, value := range details.Commands {
			matchArray := reMain.FindStringSubmatch(value.Command)
			match := ""
			if matchArray != nil {
				match = strings.ToLower(matchArray[1])
			}
			if match != "" && match == subCommand {
				return name, language
			} else if strings.ToLower(plugin.Message()) == strings.ToLower(value.Command) {
				return name, language
			}
		}
	}
	return "", ""
}

func getCommandDetails(plugin IgorPlugin, commandName string) config.LanguagePluginCommandDetails {
	return getAllCommands(plugin, "")[commandName]
}

func getAllCommands(plugin IgorPlugin, language string) map[string]config.LanguagePluginCommandDetails {
	language = getPluginLanguage(plugin, language)
	return plugin.Config().Languages()[language].Commands
}

func getDescriptionText(plugin IgorPlugin, language string) string {
	language = getPluginLanguage(plugin, language)
	return plugin.Config().Languages()[language].Description
}

func getPluginLanguage(plugin IgorPlugin, language string) string {
	if language == "" {
		language = plugin.Config().ChosenLanguage()
	}
	if _, ok := plugin.Config().Languages()[language]; !ok {
		generalConfig, _ := config.GeneralConfig()
		language = generalConfig.DefaultLanguage
	}
	return language
}

func getPluginLanguages(pluginname string) map[string]config.LanguagePluginDetails {
	generalConfig, _ := config.GeneralConfig()
	details := make(map[string]config.LanguagePluginDetails)
	for language, langConfig := range generalConfig.Languages {
		if val, ok := langConfig.Plugins[pluginname]; ok {
			details[language] = val
		}
	}
	return details
}
