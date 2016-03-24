package plugins_test

import (
	// "github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
	"os"
	"testing"
)

func TestHelp(t *testing.T) {
	plugin := plugins.Help()
	if plugin.Name() == "" {
		t.Error("No name is set for the plugin")
	}
	if plugin.Description() == "" {
		t.Error("No description is set for the plugin")
	}
}

func TestDescribe(t *testing.T) {
	plugin := plugins.Help()
	descriptions := plugin.Describe()
	if len(descriptions) != 4 {
		t.Error("Expected 4 descriptions")
	}
	expectedCommands := []string{"help", "introduce yourself", "tell me about yourself"}
	for _, command := range expectedCommands {
		if _, ok := descriptions[command]; !ok {
			t.Error("Expected the '" + command + "' command")
		}
	}
}

func TestWork(t *testing.T) {
	// Make sure it doesn't try to read the config file
	err := os.Setenv("IGOR_CONFIG", "{\"token\": \"testtoken\"}")
	if err != nil {
		t.Error("Problem setting environment variable")
	}
	request := slack.Request{}
	plugin := plugins.Help()
	// No result test
	request.Text = "fail"
	_, err = plugin.Work(request)
	if err == nil {
		t.Error("Expected failure")
	}
	// Help call, lowercase
	request.Text = "help"
	response, err := plugin.Work(request)
	if err != nil {
		t.Error("Unexpected error for help", err.Error())
	}
	if response.Text == "" {
		t.Error("Empty response")
	}
	if len(response.Attachments) == 0 {
		t.Error("Expected attachments")
	}
	if response.IsPublic() {
		t.Error("Help should not give a public response")
	}
	lowercaseText := response.Text
	// Help call, mixed case
	request.Text = "Help"
	response, err = plugin.Work(request)
	if err != nil {
		t.Error("Unexpected error for help")
	}
	if response.Text != lowercaseText {
		t.Error("Mixed case help should get same result as lowercase")
	}

	// Introduce yourself call
	request.Text = "introduce yourself"
	response, err = plugin.Work(request)
	if err != nil {
		t.Error("Unexpected error for help")
	}
	if response.Text == "" {
		t.Error("Empty response")
	}
	if len(response.Attachments) != 1 {
		t.Error("Expected an attachment")
	}
	if !response.IsPublic() {
		t.Error("Introduce yourself should not give a public response")
	}
}
