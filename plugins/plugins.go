package plugins

import (
	"github.com/ArjenSchwarz/igor/slack"
)

type IgorPlugin interface {
	Response(string) (*slack.SlackResponse, error)
	HelpMessages() []string
}

func GetPlugins() map[string]IgorPlugin {
	plugins := make(map[string]IgorPlugin)
	// var plugins []IgorPlugin
	// plugins = append(plugins, new(HelpPlugin))
	// plugins
	plugins["help"] = Help()
	// manager := Manager{
	// 	Plugins: plugins,
	// 	name:    "something",
	// }
	return plugins
}

// type Manager struct {
// 	Plugins map[string]IgorPlugin
// 	name    string
// }

// func (*Manager) getPlugins() map[string]IgorPlugin {
// 	return Manager.plugins
// }
