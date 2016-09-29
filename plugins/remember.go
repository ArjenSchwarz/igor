package plugins

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// RememberPlugin provides remember functions
type RememberPlugin struct {
	name        string
	description string
	request     slack.Request
	config      rememberConfig
}

// Remember instantiates the RememberPlugin
func Remember(request slack.Request) (IgorPlugin, error) {
	pluginName := "remember"
	pluginConfig, err := parseRememberConfig()
	if err != nil {
		return RememberPlugin{}, err
	}
	pluginConfig.languages = getPluginLanguages(pluginName)
	plugin := RememberPlugin{
		name:    pluginName,
		request: request,
		config:  pluginConfig,
	}

	return plugin, nil
}

// Work parses the request and ensures a request comes through if any triggers
// are matched. Handled triggers:
//
//  * remember
//  * remember2
func (plugin RememberPlugin) Work() (slack.Response, error) {
	response := slack.Response{}
	if plugin.config.Dynamodb == "" {
		return response, CreateNoMatchError("No DynamoDB configured")
	}
	message, language := getCommandName(plugin)
	plugin.config.chosenLanguage = language
	switch message {
	case "remember":
		tmpresponse, err := plugin.handleRemember(response)
		if err != nil {
			return tmpresponse, err
		}
		response = tmpresponse
	case "show":
		return plugin.handleShow(response)
	case "forget":
		return plugin.handleForget(response)
	case "showall":
		return plugin.handleShowAll(response)
	}
	if response.Text == "" {
		return response, CreateNoMatchError("Nothing found")
	}
	return response, nil
}

// Describe provides the triggers RememberPlugin can handle
func (plugin RememberPlugin) Describe(language string) map[string]string {
	descriptions := make(map[string]string)
	if plugin.config.Dynamodb == "" {
		return descriptions
	}
	for commandName, values := range getAllCommands(plugin, language) {
		if commandName != "forget" || plugin.request.UserInList(plugin.config.Admins) {
			descriptions[values.Command] = values.Description
		}
	}
	return descriptions
}

func (plugin RememberPlugin) handleRemember(response slack.Response) (slack.Response, error) {
	commandDetails := getCommandDetails(plugin, "remember")
	if plugin.request.UserInList(plugin.config.Blacklist) {
		response.Text = commandDetails.Texts["forbidden"]
		return response, nil
	}
	parts := strings.Split(plugin.Message(), " ")
	name := strings.TrimSpace(parts[1])
	url := strings.TrimSpace(parts[2])

	response.Text = strings.Replace(commandDetails.Texts["response_text"], "[replace]", name, 1)
	sess, err := session.NewSession()
	if err != nil {
		return response, err
	}

	svc := dynamodb.New(sess)

	params := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {S: aws.String(name)},
			"url":  {S: aws.String(url)},
			"user": {S: aws.String(plugin.request.UserName)},
		},
		TableName: aws.String(plugin.config.Dynamodb),
	}
	_, err = svc.PutItem(params)

	return response, err
}

func (plugin RememberPlugin) handleForget(response slack.Response) (slack.Response, error) {
	commandDetails := getCommandDetails(plugin, "forget")
	if !plugin.request.UserInList(plugin.config.Admins) {
		response.Text = commandDetails.Texts["forbidden"]
		return response, nil
	}
	parts := strings.Split(plugin.Message(), " ")
	name := strings.TrimSpace(parts[1])

	response.Text = strings.Replace(commandDetails.Texts["response_text"], "[replace]", name, 1)
	sess, err := session.NewSession()
	if err != nil {
		return response, err
	}

	svc := dynamodb.New(sess)

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{ // Required
			"name": {S: aws.String(name)},
		},
		TableName: aws.String(plugin.config.Dynamodb), // Required
	}
	_, err = svc.DeleteItem(params)
	return response, err
}

func (plugin RememberPlugin) handleShow(response slack.Response) (slack.Response, error) {
	var subject string
	commandDetails := getCommandDetails(plugin, "show")
	parts := strings.Split(plugin.Message(), " ")
	subject = strings.TrimSpace(strings.Replace(plugin.Message(), parts[0], "", 1))

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return response, err
	}

	svc := dynamodb.New(sess)

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"name": {S: aws.String(subject)},
		},
		TableName: aws.String(plugin.config.Dynamodb),
	}
	resp, err := svc.GetItem(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				response.Text = commandDetails.Texts["no_result"]
				return response, nil
			}
		}
		response.UnfurlMedia = true
		response.UnfurlLinks = true
		return response, err
	}
	response.Text = aws.StringValue(resp.Item["url"].S)
	// attach := slack.Attachment{
	// 	Title:     aws.StringValue(resp.Item["name"].S),
	// 	TitleLink: ,
	// 	ImageURL:  aws.StringValue(resp.Item["url"].S),
	// }
	// response.AddAttachment(attach)
	response.SetPublic()
	return response, nil
}

func (plugin RememberPlugin) handleShowAll(response slack.Response) (slack.Response, error) {
	commandDetails := getCommandDetails(plugin, "showall")
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return response, err
	}

	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(plugin.config.Dynamodb),
	}
	resp, err := svc.Scan(params)
	if err != nil {
		return response, err
	}
	if aws.Int64Value(resp.Count) == int64(0) {
		response.Text = commandDetails.Texts["no_result"]
	} else {
		response.Text = commandDetails.Texts["response_text"]
		for _, item := range resp.Items {
			response.Text += fmt.Sprintf("\n * %s (%s)",
				aws.StringValue(item["name"].S),
				aws.StringValue(item["user"].S))
		}
	}

	return response, nil
}

// Functions to satisfy the interfaces are below

// Config returns the plugin configuration
func (plugin RememberPlugin) Config() IgorConfig {
	return plugin.config
}

type rememberConfig struct {
	languages      map[string]config.LanguagePluginDetails
	chosenLanguage string
	Dynamodb       string
	Admins         []string
	Blacklist      []string
}

func parseRememberConfig() (rememberConfig, error) {
	pluginConfig := struct {
		Remember rememberConfig
	}{}

	err := config.ParseConfig(&pluginConfig)
	if err != nil {
		return pluginConfig.Remember, err
	}

	return pluginConfig.Remember, nil
}

type rememberDetails struct {
	Dynamodb string
}

// Languages returns the languages available for the plugin
func (config rememberConfig) Languages() map[string]config.LanguagePluginDetails {
	return config.languages
}

// ChosenLanguage returns the language active for this plugin
func (config rememberConfig) ChosenLanguage() string {
	return config.chosenLanguage
}

// Description returns a global description of the plugin
func (plugin RememberPlugin) Description(language string) string {
	return getDescriptionText(plugin, language)
}

// Name returns the name of the plugin
func (plugin RememberPlugin) Name() string {
	return plugin.name
}

// Message returns a formatted version of the original message
func (plugin RememberPlugin) Message() string {
	return strings.ToLower(plugin.request.Text)
}
