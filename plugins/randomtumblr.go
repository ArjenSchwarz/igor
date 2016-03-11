package plugins

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/yaml.v2"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// RandomTumblrPlugin provides random entries from Tumblr blogs
type RandomTumblrPlugin struct {
	name        string
	description string
	Config      randomTumblrConfig
}

// RandomTumblr instantiates a RandomTumblrPlugin
func RandomTumblr() IgorPlugin {
	pluginName := "randomTumblr"
	pluginConfig := parseRandomTumblrConfig()
	description := "Igor provides random entries from Tumblr blogs"
	plugin := RandomTumblrPlugin{
		name:        pluginName,
		description: description,
		Config:      pluginConfig,
	}
	return plugin
}

// Describe provides the triggers RandomTumblrPlugin can handle
func (plugin RandomTumblrPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["tumblr"] = "Shows a completely random tumblr post"
	for name, details := range plugin.Config.Tumblrs {
		descriptions["tumblr "+name] = fmt.Sprintf("Shows a random post from the %s tumblr", details.Name)
	}
	return descriptions
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
//
// * tumblr
// * tumblr [configured tumblr name]
func (plugin RandomTumblrPlugin) Work(request slack.Request) (slack.Response, error) {
	response := slack.Response{}
	var chosentumblr tumblrDetails
	if len(request.Text) == 6 && request.Text == "tumblr" {
		//Not the most efficient way of randomizing, but good enough for a small map
		rand.Seed(time.Now().UTC().UnixNano())
		list := []string{}
		for name := range plugin.Config.Tumblrs {
			list = append(list, name)
		}
		randnr := rand.Intn(len(list))
		chosentumblr = plugin.Config.Tumblrs[list[randnr]]
	} else if len(request.Text) > 6 && request.Text[:6] == "tumblr" {
		tumblr := request.Text[7:]
		for name, details := range plugin.Config.Tumblrs {
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
	return response, errors.New("No Match")
}

func addTumblrAttachment(response slack.Response, chosentumblr tumblrDetails) (slack.Response, error) {
	url := fmt.Sprintf("%s/random", chosentumblr.URL)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return response, err
	}

	title := doc.Find(chosentumblr.TitleSrc).Text()
	img, exists := doc.Find(chosentumblr.ImgSrc).Attr("src")
	if !exists {
		return response, errors.New("No image found")
	}
	attach := slack.Attachment{
		Title:    title,
		ImageURL: img,
	}
	response.AddAttachment(attach)
	return response, err
}

// parseRandomTumblrConfig collects the config as defined in the config file for
// the random Tumblr plugin
func parseRandomTumblrConfig() randomTumblrConfig {
	configFile := config.GetConfigFile()

	config := randomTumblrConfig{}

	err := yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// Description returns a global description of the plugin
func (plugin RandomTumblrPlugin) Description() string {
	return plugin.description
}

// Name returns the name of the plugin
func (plugin RandomTumblrPlugin) Name() string {
	return plugin.name
}

type randomTumblrConfig struct {
	Tumblrs map[string]tumblrDetails `yaml:"randomtumblr"`
}

type tumblrDetails struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	ImgSrc   string `yaml:"image_src"`
	TitleSrc string `yaml:"title_src"`
}
