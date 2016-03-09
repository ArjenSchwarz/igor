package plugins

import (
	"github.com/ArjenSchwarz/igor/slack"
)

type IgorPlugin interface {
	Work(slack.SlackRequest) (slack.SlackResponse, error)
	Describe() map[string]string
	Name() string
	Description() string
}

//TODO Ensure plugins can be disabled
func GetPlugins() map[string]IgorPlugin {
	plugins := make(map[string]IgorPlugin)
	plugins["help"] = Help()
	plugins["weather"] = Weather()
	plugins["tumblr"] = RandomTumblr()
	return plugins
}
