package plugins_test

import (
	"github.com/arjenschwarz/igor/plugins"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	plugin := plugins.Help()
	if plugin.Name() == "" {
		t.Error("No name is set for the plugin")
	}
	if plugin.Version() == "" {
		t.Error("No version is set for the plugin")
	}
	if plugin.Description() == "" {
		t.Error("No description is set for the plugin")
	}
	if plugin.Author() == "" {
		t.Error("No author is set for the plugin")
	}
}

func TestDescriptions(t *testing.T) {
	plugin := plugins.Help()
	descriptions := plugin.Descriptions()
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

func TestResponse(t *testing.T) {
	plugin := plugins.Help()
	// No result test
	_, err := plugin.Response("fail")
	if err == nil {
		t.Error("Expected failure")
	}
	// Help call, lowercase
	response, err := plugin.Response("help")
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
	response, err = plugin.Response("Help")
	if err != nil {
		t.Error("Unexpected error for help")
	}
	if response.Text != lowercaseText {
		t.Error("Mixed case help should get same result as lowercase")
	}

	// Introduce yourself call
	response, err = plugin.Response("introduce yourself")
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
