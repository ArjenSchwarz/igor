package plugins

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// RandomTumblrPlugin provides random entries from Tumblr blogs
type RandomTumblrPlugin struct {
	name        string
	description string
	config      randomTumblrConfig
	request     slack.Request
}

// Config returns the plugin configuration
func (plugin RandomTumblrPlugin) Config() IgorConfig {
	return plugin.config
}

// RandomTumblr instantiates a RandomTumblrPlugin
func RandomTumblr(request slack.Request) (IgorPlugin, error) {
	pluginName := "randomTumblr"
	pluginConfig := randomTumblrConfig{
		languages: getPluginLanguages(pluginName),
	}
	err := config.ParseConfig(&pluginConfig)
	if err != nil {
		return RandomTumblrPlugin{}, err
	}
	plugin := RandomTumblrPlugin{
		name:    pluginName,
		config:  pluginConfig,
		request: request,
	}
	return plugin, nil
}

// Describe provides the triggers RandomTumblrPlugin can handle
func (plugin RandomTumblrPlugin) Describe(language string) map[string]string {
	descriptions := make(map[string]string)
	pluginCommands := getAllCommands(plugin, language)
	descriptions[pluginCommands["tumblr"].Command] = pluginCommands["tumblr"].Description
	for name, details := range plugin.config.Randomtumblr {
		key := strings.Replace(pluginCommands["specifictumblr"].Command, "[replace]", "%s", -1)
		value := strings.Replace(pluginCommands["specifictumblr"].Description, "[replace]", "%s", -1)
		descriptions[fmt.Sprintf(key, name)] = fmt.Sprintf(value, details.Name)
	}
	return descriptions
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
//
// * tumblr
// * tumblr [configured tumblr name]
func (plugin RandomTumblrPlugin) Work() (slack.Response, error) {
	response := slack.Response{}
	var chosentumblr tumblrDetails
	for _, details := range plugin.config.Languages() {
		maincommand := details.Commands["tumblr"].Command
		if len(plugin.Message()) == len(maincommand) && plugin.Message() == maincommand {
			//Not the most efficient way of randomizing, but good enough for a small map
			rand.Seed(time.Now().UTC().UnixNano())
			list := []string{}
			for name := range plugin.config.Randomtumblr {
				list = append(list, name)
			}
			randnr := rand.Intn(len(list))
			chosentumblr = plugin.config.Randomtumblr[list[randnr]]
		} else if len(plugin.Message()) > len(maincommand) && plugin.Message()[:len(maincommand)] == maincommand {
			tumblr := plugin.Message()[len(maincommand)+1:]
			for name, details := range plugin.config.Randomtumblr {
				if name == tumblr {
					chosentumblr = details
				}
			}
		}
		if chosentumblr.URL != "" {
			response, err := addTumblrAttachment(response, chosentumblr)
			if err == nil {
				response.SetPublic()
			}
			return response, err
		}
	}
	return response, CreateNoMatchError("Nothing found")
}

func addTumblrAttachment(response slack.Response, chosentumblr tumblrDetails) (slack.Response, error) {
	url := fmt.Sprintf("%s/random", chosentumblr.URL)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return response, err
	}
	title := ""
	if chosentumblr.Titlesrc != "" {
		title = doc.Find(chosentumblr.Titlesrc).Text()
	}
	img, exists := doc.Find(chosentumblr.Imagesrc).Attr("src")
	if !exists {
		return response, errors.New("No image found")
	}
	attach := slack.Attachment{
		Title:    title,
		ImageURL: img,
	}
	if title == "" {
		attach.Title = chosentumblr.Name
	}
	attach.TitleLink = doc.Url.String()
	response.AddAttachment(attach)
	return response, err
}

// Description returns a global description of the plugin
func (plugin RandomTumblrPlugin) Description(language string) string {
	return getDescriptionText(plugin, language)
}

// Name returns the name of the plugin
func (plugin RandomTumblrPlugin) Name() string {
	return plugin.name
}

// Message returns a formatted version of the original message
func (plugin RandomTumblrPlugin) Message() string {
	return plugin.request.Text
}

func (config randomTumblrConfig) Languages() map[string]config.LanguagePluginDetails {
	return config.languages
}

func (config randomTumblrConfig) ChosenLanguage() string {
	return config.chosenLanguage
}

type randomTumblrConfig struct {
	Randomtumblr   map[string]tumblrDetails
	languages      map[string]config.LanguagePluginDetails
	chosenLanguage string
}

type tumblrDetails struct {
	Name     string
	URL      string
	Imagesrc string
	Titlesrc string
}
