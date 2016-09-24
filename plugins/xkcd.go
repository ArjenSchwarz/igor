package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// XkcdPlugin provides access to XKCD comics
type XkcdPlugin struct {
	name        string
	description string
	request     slack.Request
	config      xkcdConfig
}

// Config returns the plugin configuration
func (plugin XkcdPlugin) Config() IgorConfig {
	return plugin.config
}

type xkcdConfig struct {
	languages      map[string]config.LanguagePluginDetails
	chosenLanguage string
}

func (config xkcdConfig) Languages() map[string]config.LanguagePluginDetails {
	return config.languages
}

func (config xkcdConfig) ChosenLanguage() string {
	return config.chosenLanguage
}

// Xkcd is a plugin that returns XKCD comics
func Xkcd(request slack.Request) (IgorPlugin, error) {
	pluginName := "xkcd"
	pluginConfig := xkcdConfig{
		languages: getPluginLanguages(pluginName),
	}
	plugin := XkcdPlugin{
		name:    pluginName,
		request: request,
		config:  pluginConfig,
	}
	return plugin, nil
}

type xkcdEntry struct {
	Number int    `json:"num"`
	Title  string `json:"title"`
	Alt    string `json:"alt"`
	Image  string `json:"img"`
}

// Work parses the request and ensures a request comes through if any triggers
// are matched.
func (plugin XkcdPlugin) Work() (slack.Response, error) {
	response := slack.Response{}
	message, language := getCommandName(plugin)
	if message == "" {
		return response, CreateNoMatchError("Nothing found")
	}
	plugin.config.chosenLanguage = language
	response.SetPublic()
	baseurl := "http://xkcd.com/"
	jsoncall := "info.0.json"
	switch message {
	case "xkcd":
		url := fmt.Sprintf("%s%s", baseurl, jsoncall)
		return plugin.parseXkcdMessage(url, response)
	case "xkcd_random":
		url := fmt.Sprintf("%s%s", baseurl, jsoncall)
		entry, err := getXkcdMessage(url)
		if err != nil {
			return response, err
		}
		rand.Seed(time.Now().UTC().UnixNano())
		comicnr := rand.Intn(entry.Number)
		url = fmt.Sprintf("%s%v/%s", baseurl, comicnr, jsoncall)
		return plugin.parseXkcdMessage(url, response)
	case "xkcd_specific":
		parts := strings.Split(plugin.Message(), " ")
		comicnr := ""
		if len(parts) > 1 {
			comicnr = strings.TrimSpace(strings.Replace(plugin.Message(), parts[0], "", 1))
		}
		url := fmt.Sprintf("%s%v/%s", baseurl, comicnr, jsoncall)
		return plugin.parseXkcdMessage(url, response)
	}
	return response, CreateNoMatchError("Nothing found")
}

func getXkcdMessage(url string) (xkcdEntry, error) {
	parsedResult := xkcdEntry{}
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return parsedResult, err
	}
	if resp.StatusCode == 404 {
		return parsedResult, errors.New("Incorrect comic number")
	}
	err = json.NewDecoder(resp.Body).Decode(&parsedResult)
	return parsedResult, err
}

func (plugin XkcdPlugin) parseXkcdMessage(url string, response slack.Response) (slack.Response, error) {
	parsedResult, err := getXkcdMessage(url)
	if err != nil {
		return response, err
	}
	commandDetails := getCommandDetails(plugin, "xkcd")
	response.Text = fmt.Sprintf("%s%v", commandDetails.Texts["response_text"], parsedResult.Number)
	attach := slack.Attachment{
		ImageURL: parsedResult.Image,
		Text:     parsedResult.Alt,
		Title:    parsedResult.Title,
	}
	response.AddAttachment(attach)
	return response, nil
}

// Describe provides the triggers the plugin can handle
func (plugin XkcdPlugin) Describe(language string) map[string]string {
	descriptions := make(map[string]string)

	for _, values := range getAllCommands(plugin, language) {
		descriptions[values.Command] = values.Description
	}
	return descriptions
}

// Name returns the name of the plugin
func (plugin XkcdPlugin) Name() string {
	return plugin.name
}

// Description returns a global description of the plugin
func (plugin XkcdPlugin) Description(language string) string {
	return getDescriptionText(plugin, language)
}

// Message returns the original request message
func (plugin XkcdPlugin) Message() string {
	return strings.ToLower(plugin.request.Text)
}
