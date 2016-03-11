package plugins_test

import (
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
	"strings"
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
	if len(descriptions) != 2 {
		t.Error("Expected 2 descriptions")
	}
	expectedCommands := []string{"help", "introduce yourself"}
	for _, command := range expectedCommands {
		if _, ok := descriptions[command]; !ok {
			t.Error("Expected the '" + command + "' command")
		}
	}
}

func TestWork(t *testing.T) {
	request := slack.Request{}
	plugin := plugins.Help()
	// No result test
	request.Text = "fail"
	_, err := plugin.Work(request)
	if err == nil {
		t.Error("Expected failure")
	}
	// Help call, lowercase
	request.Text = "help"
	response, err := plugin.Work(request)
	if err != nil {
		t.Error("Unexpected error for help")
	}
	if response.Text == "" {
		t.Error("Empty response")
	}
	if len(response.Attachments) != 1 {
		t.Error("Expected an attachment")
	}
	if !strings.Contains(response.Attachments[0].Text, "help") {
		t.Error("Expected help command in description")
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
