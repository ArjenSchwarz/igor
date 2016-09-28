[![wercker status](https://app.wercker.com/status/eea144a7251e1b84d514904e19eff205/s/master "wercker status")](https://app.wercker.com/project/byKey/eea144a7251e1b84d514904e19eff205)

# Igor

Igor is a Slack Slash command that acts like a bot. The code is written in Go and is designed to be run on AWS Lambda through a NodeJS wrapper or as a Docker container. All commands are handled through plugins, making it extendable.

The name is based on Sir Terry Pratchett's wonderful use of Dr. Frankenstein's servant. No disrespect intended.

Igor is currently early in development, and can't do much yet, but it is usable.

# Available Plugins

* [Help](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Help), explains Igor
* [Weather](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Weather), get the current weather and forecasts
* [(Random) Tumblr](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-(Random)-Tumblr) image, get a random image from a Tumblr blog
* [Status](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Status), get the current status of webservices like GitHub and Bitbucket
* XKCD, get the latest XKCD comic

# Language support

Igor is built to understand multiple languages. The language files are stored in the language directory, and are yaml files. If you wish to add a language create a file to put in there following the structure of the existing files. If you don't wish to provide a translation for a specific plugin you can leave it out as it will gracefully fall back to the default language. The default language is defined in the configuration as `defaultlanguage: yourlanguage` and defaults to `english`.

If you set the default language to a language that doesn't have all plugins implemented, it will be possible to make Igor unable to comply. This will make Igor sad and it might even crash. So this is not recommended.

# TODO

Many things, have a look at the [Roadmap](https://github.com/ArjenSchwarz/igor/wiki/Roadmap) for the current ideas.

# Installation

Please have a look at the [Wiki](https://github.com/ArjenSchwarz/igor/wiki) for a full list of installation and running options.

# Contribute

If you wish to contribute in any way (reporting bugs, requesting features, writing code), feel free to do so either by opening Issues or Pull Request. For Pull Requests, just follow the standard pattern.

1. Fork the repository
2. Make your changes
3. Make a pull request that explains what it does

To make plugin development easier, there is a snippet for Sublime Text included in the devtools directory. If you copy this to your User package you can easily create the skeleton for a plugin with it.

You can also test your commands locally using `bin/testcommand.sh`. This script will read your config.yml file and based on that it will generate a correctly formatted json string and provide that to the binary.

For example:

```bash
$ bin/testcommand.sh "introduce yourself"
{"text":"I am Igor, reprethenting We-R-Igors.","response_type":"in_channel","attachments":[{"title":"A Spare Hand When Needed","text":"We come from Überwald, but are alwayth where we are needed motht.\nRun */igor help* to see which Igors are currently available.","mrkdwn_in":["text"]}]}
```
