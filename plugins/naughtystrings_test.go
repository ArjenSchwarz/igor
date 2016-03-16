package plugins_test

import (
	"bytes"
	"encoding/json"
	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/plugins"
	"github.com/ArjenSchwarz/igor/slack"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestNaughtyStrings(t *testing.T) {
	config.SetRawConfig([]byte("token: testtoken"))
	var list []string

	filename, _ := filepath.Abs("../devtools/blns.json")
	c, _ := ioutil.ReadFile(filename)
	dec := json.NewDecoder(bytes.NewReader(c))
	dec.Decode(&list)
	log.Println(len(list))
	for _, plugin := range plugins.GetPlugins(config.GeneralConfig()) {
		for _, string := range list {
			request := slack.Request{Text: string}
			_, err := plugin.Work(request)
			if err != nil && err.Error() != "No match" {
				t.Error("Failed naughty string: " + string + " - " + plugin.Name() + " > " + err.Error())
			}
		}
	}
}
