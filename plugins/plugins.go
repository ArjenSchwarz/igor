package plugins

import (
	"github.com/ArjenSchwarz/igor/slack"
)

type IgorPlugin interface {
	Response(string) (slack.SlackResponse, error)
	Descriptions() map[string]string
	Name() string
	Version() string
	Description() string
	Author() string
}

//TODO Ensure plugins can be disabled
func GetPlugins() map[string]IgorPlugin {
	plugins := make(map[string]IgorPlugin)
	plugins["help"] = Help()
	return plugins
}
