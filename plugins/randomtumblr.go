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

// RandomTumblr instantiates a RandomTumblrPlugin
func RandomTumblr() RandomTumblrPlugin {
	pluginName := "randomTumblr"
	pluginConfig := ParseRandomTumblrConfig()
	description := "Igor provides random entries from Tumblr blogs"
	plugin := RandomTumblrPlugin{
		name:        pluginName,
		description: description,
		Config:      pluginConfig,
	}
	return plugin
}

func (r RandomTumblrPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["tumblr"] = "Shows a completely random tumblr post"
	for name, details := range r.Config.Tumblrs {
		descriptions["tumblr "+name] = fmt.Sprintf("Shows a random post from the %s tumblr", details.Name)
	}
	return descriptions
}

func (r RandomTumblrPlugin) Work(request slack.SlackRequest) (slack.SlackResponse, error) {
	response := slack.SlackResponse{}
	var chosentumblr TumblrDetails
	if len(request.Text) == 6 && request.Text == "tumblr" {
		//Not the most efficient way of randomizing, but good enough for a small map
		rand.Seed(time.Now().UTC().UnixNano())
		list := []string{}
		for name := range r.Config.Tumblrs {
			list = append(list, name)
		}
		randnr := rand.Intn(len(list))
		chosentumblr = r.Config.Tumblrs[list[randnr]]
	} else if len(request.Text) > 6 && request.Text[:6] == "tumblr" {
		tumblr := request.Text[7:]
		for name, details := range r.Config.Tumblrs {
			if name == tumblr {
				chosentumblr = details
			}
		}
	}
	if chosentumblr.Url != "" {
		response, err := addTumblrAttachment(response, chosentumblr)
		if err == nil {
			response.SetPublic()
		}
		return response, err
	}
	return response, errors.New("No Match")
}

func addTumblrAttachment(response slack.SlackResponse, chosentumblr TumblrDetails) (slack.SlackResponse, error) {
	url := fmt.Sprintf("%s/random", chosentumblr.Url)
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
		ImageUrl: img,
	}
	response.AddAttachment(attach)
	return response, err
}

// ParseRandomTumblrConfig collects the config as defined in the config file for
// the random Tumblr plugin
func ParseRandomTumblrConfig() RandomTumblrConfig {
	configFile := config.GetConfigFile()

	config := RandomTumblrConfig{}

	err := yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func (p RandomTumblrPlugin) Description() string {
	return p.description
}
func (p RandomTumblrPlugin) Name() string {
	return p.name
}

type RandomTumblrPlugin struct {
	name        string
	description string
	Config      RandomTumblrConfig
}

type RandomTumblrConfig struct {
	Tumblrs map[string]TumblrDetails `yaml:"randomtumblr"`
}

type TumblrDetails struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	ImgSrc   string `yaml:"image_src"`
	TitleSrc string `yaml:"title_src"`
}
