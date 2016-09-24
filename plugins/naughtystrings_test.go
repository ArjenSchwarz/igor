package plugins_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
)

// TestNaughtyStrings calls every plugin with the list of naughtystrings
// (https://github.com/minimaxir/big-list-of-naughty-strings).
func TestNaughtyStrings(t *testing.T) {
	err := os.Setenv("IGOR_CONFIG", "{\"token\": \"testtoken\", \"languagedir\": \"../language\"}")
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
	request := slack.Request{Text: "string"}
	for _, plugin := range plugins.GetPlugins(request, config) {
		t.Run("Plugin="+plugin.Name(), func(t *testing.T) {
			for _, string := range list {
				_, err := plugin.Work()
				if err != nil {
					switch err.(type) {
					case *plugins.NoMatchError:
					default:
						t.Error(fmt.Sprintf("Failed naughty string: %s - %s > %s",
							string,
							plugin.Name(),
							err.Error()))
					}
				}
			}

		})
	}
}
