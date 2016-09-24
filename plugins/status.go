package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// StatusPlugin provides status reports for various services
type StatusPlugin struct {
	name        string
	description string
	config      statusConfig
	Checks      map[string]func() (slack.Attachment, error)
	MainChecks  map[string]func() (slack.Attachment, error)
	request     slack.Request
}

type statusConfig struct {
	Main           []string
	languages      map[string]config.LanguagePluginDetails
	chosenLanguage string
}

// Config returns the plugin configuration
func (plugin StatusPlugin) Config() IgorConfig {
	return plugin.config
}

func (config statusConfig) Languages() map[string]config.LanguagePluginDetails {
	return config.languages
}

func (config statusConfig) ChosenLanguage() string {
	return config.chosenLanguage
}

// Status instantiates the StatusPlugin
func Status(request slack.Request) (IgorPlugin, error) {
	pluginName := "status"
	pluginConfig, err := parseStatusConfig()
	if err != nil {
		return StatusPlugin{}, err
	}
	pluginConfig.languages = getPluginLanguages(pluginName)
	plugin := StatusPlugin{
		name:        pluginName,
		description: "",
		config:      pluginConfig,
		request:     request,
	}
	statuschecks := make(map[string]func() (slack.Attachment, error))
	statuschecks["github"] = plugin.handleGitHubStatus
	statuschecks["bitbucket"] = plugin.handleBitbucketStatus
	statuschecks["npmjs"] = plugin.handleNpmjsStatus
	statuschecks["disqus"] = plugin.handleDisqusStatus
	statuschecks["cloudflare"] = plugin.handleCloudflareStatus
	statuschecks["aws"] = plugin.handleShortAWSStatus
	statuschecks["travis"] = plugin.handleTravisCIStatus
	statuschecks["docker"] = plugin.handleDockerStatus
	plugin.Checks = statuschecks

	if len(pluginConfig.Main) == 0 {
		plugin.MainChecks = statuschecks
	} else {
		mainchecks := make(map[string]func() (slack.Attachment, error))
		for _, check := range pluginConfig.Main {
			if val, ok := statuschecks[check]; ok {
				mainchecks[check] = val
			}
		}
		plugin.MainChecks = mainchecks
	}

	return plugin, nil
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
func (plugin StatusPlugin) Work() (slack.Response, error) {
	statuschecks := plugin.Checks
	response := slack.Response{}
	message, language := getCommandName(plugin)
	plugin.config.chosenLanguage = language
	if message == "status" {
		c := make(chan slack.Attachment)
		for _, function := range plugin.MainChecks {
			go func(function func() (slack.Attachment, error)) {
				attachment, err := function()
				if err != nil {
					// return response, err
				}
				c <- attachment
			}(function)
		}
		for i := 0; i < len(plugin.MainChecks); i++ {
			response.AddAttachment(<-c)
		}
		commandDetails := getCommandDetails(plugin, "status")
		response.Text = commandDetails.Texts["response_text"]
		response.SetPublic()
	} else if message == "status_aws" {
		attachments, _ := plugin.handleAWSStatus()
		for _, attachment := range attachments {
			response.AddAttachment(attachment)
		}
		commandDetails := getCommandDetails(plugin, "status_aws")
		response.Text = commandDetails.Texts["response_text"]
		response.SetPublic()
	} else if message == "status_service" || message == "status_url" {
		parts := strings.Split(plugin.Message(), " ")
		tocheck := ""
		if len(parts) > 1 {
			tocheck = strings.TrimSpace(strings.Replace(plugin.Message(), parts[0], "", 1))
		}
		// Check if this is a predefined service
		if function, ok := statuschecks[tocheck]; ok {
			// Treat it as a predefined service
			attachment, err := function()
			if err != nil {
				fmt.Println(err)
				return response, err
			}
			response.AddAttachment(attachment)

			commandDetails := getCommandDetails(plugin, "status_service")
			response.Text = commandDetails.Texts["response_text"]
			response.SetPublic()
		} else {
			// Treat it as a website
			attachment, err := plugin.handleDomain(plugin.Message()[7:])
			if err != nil {
				return response, err
			}
			response.AddAttachment(attachment)
			commandDetails := getCommandDetails(plugin, "status_url")
			response.Text = commandDetails.Texts["response_text"]
			response.SetPublic()
		}
	}
	if response.Text == "" {
		return response, CreateNoMatchError("Nothing found")
	}
	return response, nil
}

// Describe provides the triggers StatusPlugin can handle
func (plugin StatusPlugin) Describe(language string) map[string]string {
	// Get a list of all services
	var servicelist []string
	for service := range plugin.Checks {
		servicelist = append(servicelist, service)
	}
	services := strings.Join(servicelist, ", ")

	descriptions := make(map[string]string)
	for name, values := range getAllCommands(plugin, language) {
		if name == "status_service" {
			cleanedDescription := strings.Replace(values.Description, "[replace]", "%s", -1)
			descriptions[values.Command] = fmt.Sprintf(cleanedDescription, services)
		} else {
			descriptions[values.Command] = values.Description
		}
	}
	return descriptions
}

// Description returns a global description of the plugin
func (plugin StatusPlugin) Description(language string) string {
	return getDescriptionText(plugin, language)
}

// Name returns the name of the plugin
func (plugin StatusPlugin) Name() string {
	return plugin.name
}

// Message returns a formatted version of the original message
func (plugin StatusPlugin) Message() string {
	return plugin.request.Text
}

func (plugin StatusPlugin) handleDomain(domain string) (slack.Attachment, error) {
	attachment := slack.Attachment{Title: domain}
	commandDetails := getCommandDetails(plugin, "status_url")
	resp, err := http.Get(fmt.Sprintf("https://isitup.org/%s.json", domain))
	defer resp.Body.Close()
	if err != nil {
		return attachment, err
	}
	var result struct {
		StatusCode int64 `json:"status_code"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return attachment, err
	}
	switch result.StatusCode {
	case 1:
		attachment.Color = slack.ResponseGood
		attachment.Text = commandDetails.Texts["good"]
	case 2:
		attachment.Color = slack.ResponseBad
		attachment.Text = commandDetails.Texts["bad"]
	default:
		return attachment, errors.New("Not a valid domain")
	}
	return attachment, nil
}

func (StatusPlugin) handleGitHubStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "GitHub", PreText: "http://status.github.com"}
	resp, err := http.Get("https://status.github.com/api/last-message.json")
	defer resp.Body.Close()
	if err != nil {
		return attachment, err
	}
	var result struct {
		Status string `json:"status"`
		Body   string `json:"body"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return attachment, err
	}
	attachment.Text = result.Body
	switch result.Status {
	case "good":
		attachment.Color = slack.ResponseGood
	case "minor":
		attachment.Color = slack.ResponseWarning
	case "major":
		attachment.Color = slack.ResponseBad
	}
	return attachment, nil
}

func (plugin StatusPlugin) handleBitbucketStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "Bitbucket", PreText: "http://status.bitbucket.org"}
	return plugin.handleStatusPageIo(attachment)
}

func (plugin StatusPlugin) handleNpmjsStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "NPM", PreText: "http://status.npmjs.org"}
	return plugin.handleStatusPageIo(attachment)
}

func (plugin StatusPlugin) handleDisqusStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "Disqus", PreText: "http://status.disqus.com"}
	return plugin.handleStatusPageIo(attachment)
}

func (plugin StatusPlugin) handleCloudflareStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "Cloudflare", PreText: "http://cloudflarestatus.com"}
	return plugin.handleStatusPageIo(attachment)
}

func (plugin StatusPlugin) handleTravisCIStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "Travis CI", PreText: "https://www.traviscistatus.com"}
	return plugin.handleStatusPageIo(attachment)
}

func (plugin StatusPlugin) handleAWSStatus() ([]slack.Attachment, error) {
	attachments := []slack.Attachment{}
	mainAttachment := slack.Attachment{Title: "AWS", PreText: "http://status.aws.amazon.com"}
	attachments = append(attachments, mainAttachment)
	nrResolved := 0
	nrProblems := 0

	commandDetails := getCommandDetails(plugin, "status_aws")

	doc, err := goquery.NewDocument(mainAttachment.PreText)
	if err != nil {
		return attachments, err
	}
	doc.Find("div#current_events_block table tr").Each(func(i int, s *goquery.Selection) {
		message := strings.Trim(s.Find("td").Eq(2).Text(), " \n")
		if message != "Service is operating normally" && message != "" {
			message = strings.Replace(message, "\n            more \n        \n        \n      ", "\n", 1)
			message = strings.Replace(message, ".", ".\n", -1)
			service := s.Find("td").Eq(1).Text()
			attachment := slack.Attachment{Title: service, Text: message}
			if message[:10] == "[RESOLVED]" {
				attachment.Color = slack.ResponseWarning
				nrResolved++
			} else {
				attachment.Color = slack.ResponseBad
				nrProblems++
			}
			attachments = append(attachments, attachment)
		}
	})
	if nrProblems != 0 {
		mainAttachment.Color = slack.ResponseBad
		mainAttachment.Text = fmt.Sprintf("%s: %s", commandDetails.Texts["nr_issues"], strconv.Itoa(nrProblems))
		if nrResolved != 0 {
			mainAttachment.Text += fmt.Sprintf("\n%s: %s", commandDetails.Texts["nr_resolved_issues"], strconv.Itoa(nrResolved))
		}
	} else if nrResolved != 0 {
		mainAttachment.Color = slack.ResponseWarning
		mainAttachment.Text = fmt.Sprintf("%s: %s", commandDetails.Texts["nr_resolved_issues"], strconv.Itoa(nrResolved))
	} else {
		mainAttachment.Color = slack.ResponseGood
		mainAttachment.Text = commandDetails.Texts["ok"]
	}
	attachments[0] = mainAttachment

	return attachments, nil
}

func (plugin StatusPlugin) handleShortAWSStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "AWS", PreText: "http://status.aws.amazon.com"}
	nrResolved := 0
	nrProblems := 0
	commandDetails := getCommandDetails(plugin, "status_aws")

	doc, err := goquery.NewDocument(attachment.PreText)
	if err != nil {
		return attachment, err
	}
	doc.Find("div#current_events_block table tr").Each(func(i int, s *goquery.Selection) {
		message := strings.Trim(s.Find("td").Eq(2).Text(), " \n")
		if message != "Service is operating normally" && message != "" {
			message = strings.Replace(message, "\n            more \n        \n        \n      ", "\n", 1)
			message = strings.Replace(message, ".", ".\n", -1)
			if message[:10] == "[RESOLVED]" {
				nrResolved++
			} else {
				nrProblems++
			}
		}
	})
	if nrProblems != 0 {
		attachment.Color = slack.ResponseBad
		attachment.Text = fmt.Sprintf("%s: %s\n", commandDetails.Texts["nr_issues"], strconv.Itoa(nrProblems))
		if nrResolved != 0 {
			attachment.Text += fmt.Sprintf("%s: %s\n", commandDetails.Texts["nr_resolved_issues"], strconv.Itoa(nrResolved))
		}
		attachment.Text += commandDetails.Texts["more_details"]
	} else if nrResolved != 0 {
		attachment.Color = slack.ResponseWarning
		attachment.Text = fmt.Sprintf("%s: %s\n", commandDetails.Texts["nr_resolved_issues"], strconv.Itoa(nrResolved))
		attachment.Text += commandDetails.Texts["more_details"]
	} else {
		attachment.Color = slack.ResponseGood
		attachment.Text = commandDetails.Texts["ok"]
	}

	return attachment, nil
}

func (plugin StatusPlugin) handleDockerStatus() (slack.Attachment, error) {
	attachment := slack.Attachment{Title: "Docker", PreText: "http://status.status.io"}
	return plugin.handleStatusIo(attachment)
}

func (StatusPlugin) handleStatusPageIo(attachment slack.Attachment) (slack.Attachment, error) {
	doc, err := goquery.NewDocument(attachment.PreText)
	if err != nil {
		return attachment, err
	}
	attachment.Text = doc.Find("div.page-status span.status").Text()
	pageStatus := doc.Find("div.page-status")
	if pageStatus.HasClass("status-none") {
		attachment.Color = slack.ResponseGood
	} else if pageStatus.HasClass("status-yellow") {
		attachment.Color = slack.ResponseWarning
	} else {
		attachment.Color = slack.ResponseBad
	}
	return attachment, nil
}

func (StatusPlugin) handleStatusIo(attachment slack.Attachment) (slack.Attachment, error) {
	doc, err := goquery.NewDocument(attachment.PreText)
	if err != nil {
		return attachment, err
	}
	attachment.Text = doc.Find("#statusbar_text").Text()
	if attachment.Text == "All Systems Operational" {
		attachment.Color = slack.ResponseGood
	} else {
		attachment.Color = slack.ResponseBad
	}
	return attachment, nil
}

func parseStatusConfig() (statusConfig, error) {
	pluginConfig := struct {
		Status statusConfig
	}{}

	err := config.ParseConfig(&pluginConfig)
	return pluginConfig.Status, err
}
