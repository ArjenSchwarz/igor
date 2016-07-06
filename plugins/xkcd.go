package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ArjenSchwarz/igor/slack"
)

// XkcdPlugin provides access to XKCD comics
type XkcdPlugin struct {
	name        string
	description string
	request     slack.Request
}

// Xkcd is a plugin that returns XKCD comics
func Xkcd(request slack.Request) (IgorPlugin, error) {
	pluginName := "xkcd"
	description := "Igor provides random entries from Tumblr blogs"
	plugin := XkcdPlugin{
		name:        pluginName,
		description: description,
		request:     request,
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
	if len(plugin.Message()) >= 4 && plugin.Message()[:4] == "xkcd" {
		response.SetPublic()
		baseurl := "http://xkcd.com/"
		jsoncall := "info.0.json"
		if plugin.Message() == "xkcd" {
			url := fmt.Sprintf("%s%s", baseurl, jsoncall)
			return parseXkcdMessage(url, response)
		} else if plugin.Message() == "xkcd random" {
			url := fmt.Sprintf("%s%s", baseurl, jsoncall)
			entry, err := getXkcdMessage(url)
			if err != nil {
				return response, err
			}
			rand.Seed(time.Now().UTC().UnixNano())
			comicnr := rand.Intn(entry.Number)
			url = fmt.Sprintf("%s%v/%s", baseurl, comicnr, jsoncall)
			return parseXkcdMessage(url, response)
		} else {
			comicnr := plugin.Message()[5:]
			url := fmt.Sprintf("%s%v/%s", baseurl, comicnr, jsoncall)
			return parseXkcdMessage(url, response)
		}
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

func parseXkcdMessage(url string, response slack.Response) (slack.Response, error) {
	parsedResult, err := getXkcdMessage(url)
	if err != nil {
		return response, err
	}
	response.Text = fmt.Sprintf("XKCD #%v", parsedResult.Number)
	attach := slack.Attachment{
		ImageURL: parsedResult.Image,
		Text:     parsedResult.Alt,
		Title:    parsedResult.Title,
	}
	response.AddAttachment(attach)
	return response, nil
}

// Describe provides the triggers the plugin can handle
func (plugin XkcdPlugin) Describe() map[string]string {
	descriptions := make(map[string]string)
	descriptions["xkcd"] = "Get the latest XKCD comic"
	descriptions["xkcd random"] = "Get a random XKCD comic"
	descriptions["xkcd [nr]"] = "Get a specific XKCD comic"
	return descriptions
}

// Name returns the name of the plugin
func (plugin XkcdPlugin) Name() string {
	return plugin.name
}

// Description returns a global description of the plugin
func (plugin XkcdPlugin) Description() string {
	return plugin.description
}

// Message returns the original request message
func (plugin XkcdPlugin) Message() string {
	return strings.ToLower(plugin.request.Text)
}
