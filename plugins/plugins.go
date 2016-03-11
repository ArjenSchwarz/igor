package plugins

import (
	"github.com/ArjenSchwarz/igor/config"
	"github.com/ArjenSchwarz/igor/slack"
)

// IgorPlugin is the interface that needs to be followed by all plugins
type IgorPlugin interface {
	Work(slack.Request) (slack.Response, error)
	Describe() map[string]string
	Name() string
	Description() string
}

// GetPlugins retrieves all the plugins that are activated. It checks the
// config for a whitelist and blacklist as well.
func GetPlugins(config config.Config) map[string]IgorPlugin {
	plugins := make(map[string]IgorPlugin)
	plugins["help"] = Help()
	plugins["weather"] = Weather()
	plugins["tumblr"] = RandomTumblr()
	plugins["status"] = Status()

	// Whitelist plugins
	if config.Whitelist != nil {
		whitelist := make(map[string]IgorPlugin)
		whitelist["help"] = Help() //Help is always required
		for _, allowedPlugin := range config.Whitelist {
			whitelist[allowedPlugin] = plugins[allowedPlugin]
		}
		plugins = whitelist
	}

	// Blacklist plugins
	if config.Blacklist != nil {
		for _, pluginname := range config.Blacklist {
			if pluginname != "help" { // Help is always required
				delete(plugins, pluginname)
			}
		}
	}
	return plugins
}
