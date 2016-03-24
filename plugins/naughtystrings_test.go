package plugins_test

import (
	"bytes"
	"encoding/json"
	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestNaughtyStrings calls every plugin with the list of naughtystrings
// (https://github.com/minimaxir/big-list-of-naughty-strings).
func TestNaughtyStrings(t *testing.T) {
	err := os.Setenv("IGOR_CONFIG", "{\"token\": \"testtoken\"}")
	if err != nil {
		t.Error("Problem setting environment variable")
	}
	var list []string

	filename, _ := filepath.Abs("../devtools/blns.json")
	c, _ := ioutil.ReadFile(filename)
	dec := json.NewDecoder(bytes.NewReader(c))
	dec.Decode(&list)
	config, err := config.GeneralConfig()
	if err != nil {
		t.Error("Problem getting config")
	}
	for _, plugin := range plugins.GetPlugins(config) {
		for _, string := range list {
			request := slack.Request{Text: string}
			_, err := plugin.Work(request)
			if err != nil {
				switch err.(type) {
				case *plugins.NoMatchError:
				default:
					t.Error("Failed naughty string: " + string + " - " + plugin.Name() + " > " + err.Error())
				}
			}
		}
	}
}
