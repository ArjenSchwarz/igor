package plugins

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// RandomTumblrPlugin provides random entries from Tumblr blogs
type RandomTumblrPlugin struct {
	name        string
	description string
	Config      randomTumblrConfig
	request     slack.Request
}

// RandomTumblr instantiates a RandomTumblrPlugin
func RandomTumblr(request slack.Request) (IgorPlugin, error) {
	pluginName := "randomTumblr"
	pluginConfig := randomTumblrConfig{}
	err := config.ParseConfig(&pluginConfig)
	if err != nil {
		return RandomTumblrPlugin{}, err
	}
	description := "Igor provides random entries from Tumblr blogs"
	plugin := RandomTumblrPlugin{
		name:        pluginName,
		description: description,
		Config:      pluginConfig,
		request:     request,
	}
	return plugin, nil
}

// Describe provides the triggers RandomTumblrPlugin can handle
func (plugin RandomTumblrPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["tumblr"] = "Shows a completely random tumblr post"
	for name, details := range plugin.Config.Randomtumblr {
		descriptions["tumblr "+name] = fmt.Sprintf("Shows a random post from the %s tumblr", details.Name)
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
	if len(plugin.Message()) == 6 && plugin.Message() == "tumblr" {
		//Not the most efficient way of randomizing, but good enough for a small map
		rand.Seed(time.Now().UTC().UnixNano())
		list := []string{}
		for name := range plugin.Config.Randomtumblr {
			list = append(list, name)
		}
		randnr := rand.Intn(len(list))
		chosentumblr = plugin.Config.Randomtumblr[list[randnr]]
	} else if len(plugin.Message()) > 6 && plugin.Message()[:6] == "tumblr" {
		tumblr := plugin.Message()[7:]
		for name, details := range plugin.Config.Randomtumblr {
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
	response.AddAttachment(attach)
	return response, err
}

// Description returns a global description of the plugin
func (plugin RandomTumblrPlugin) Description() string {
	return plugin.description
}

// Name returns the name of the plugin
func (plugin RandomTumblrPlugin) Name() string {
	return plugin.name
}

func (plugin RandomTumblrPlugin) Message() string {
	return plugin.request.Text
}

type randomTumblrConfig struct {
	Randomtumblr map[string]tumblrDetails
}

type tumblrDetails struct {
	Name     string
	URL      string
	Imagesrc string
	Titlesrc string
}
