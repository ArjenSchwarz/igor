[![wercker status](https://app.wercker.com/status/eea144a7251e1b84d514904e19eff205/s/master "wercker status")](https://app.wercker.com/project/byKey/eea144a7251e1b84d514904e19eff205)

# Igor

Igor is a Slack Slash command that acts like a bot. The code is written in Go and is designed to run on Lambda or as a Docker container. All commands are handled through plugins, making it extendable.

The name is based on Sir Terry Pratchett's wonderful use of Dr. Frankenstein's servant. No disrespect intended.

Igor is currently early in development, and can't do much yet, but it is usable.

# Available Plugins

* [Help](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Help), explains Igor
* [Weather](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Weather), get the current weather and forecasts
* [(Random) Tumblr](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-(Random)-Tumblr) image, get a random image from a Tumblr blog
* [Status](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Status), get the current status of webservices like GitHub and Bitbucket
* [XKCD](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-XKCD), get the latest (or a specific/random) XKCD comic
* [Remember](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Remember), save and display links to photos

# Language support

Igor is built to understand multiple languages. The language files are stored in the language directory, and are yaml files. If you wish to add a language create a file to put in there following the structure of the existing files. If you don't wish to provide a translation for a specific plugin you can leave it out as it will gracefully fall back to the default language. The default language is defined in the configuration as `defaultlanguage: yourlanguage` and defaults to `english`.

If you set the default language to a language that doesn't have all plugins implemented, it will be possible to make Igor unable to comply. This will make Igor sad and it might even crash. So this is not recommended.

# KMS support

There is the option to encrypt the various tokens in you config file using KMS. This means you create a new KMS key (in the region you run Igor) or use the default one provided by AWS. Once you have access to a key you can encrypt your tokens easily using the AWS CLI:

```bash
aws kms encrypt --key-id alias/igorkey --plaintext YOURTOKEN
```

The resulting output then has a CiphertextBlob containing the encrypted value. You can then put this in the place of your plain text token values in the config file. Additionally you will need to mark the config as using KMS by adding the `kms: true` flag.

Take note! You will have to encode all tokens in your configuration once you enable KMS. At the moment that means:

* token (your Slack token)
* weather:apitoken (your open weathermap token)

The last thing you need to do is ensure that your Igor function has usage access to the key, by allowing the role to have that access.

# DynamoDB support

The Remember plugin uses DynamoDB to store its data. You will need to create a table and give your Igor function access to it. See the [plugin's page](https://github.com/ArjenSchwarz/igor/wiki/Plugin:-Remember) for more details.

# TODO

Many things, have a look at the [Roadmap](https://github.com/ArjenSchwarz/igor/wiki/Roadmap) for the current ideas.

# Installation

Please have a look at the [Wiki](https://github.com/ArjenSchwarz/igor/wiki) for a full list of installation and running options.

# Contribute

If you wish to contribute in any way (reporting bugs, requesting features, writing code), feel free to do so either by opening Issues or Pull Request. For Pull Requests, just follow the standard pattern.

1. Fork the repository
2. Make your changes
3. Make a pull request that explains what it does

To make plugin development easier, there is an example plugin available in the devtools directory (example.go.plugin).

You can also test your commands locally using `bin/testcommand.sh`. This script will read your config.yml file and based on that it will generate a correctly formatted json string and provide that to the binary.

For example:

```bash
$ bin/testcommand.sh "introduce yourself"
{"text":"I am Igor, reprethenting We-R-Igors.","response_type":"in_channel","attachments":[{"title":"A Spare Hand When Needed","text":"We come from Überwald, but are alwayth where we are needed motht.\nRun */igor help* to see which Igors are currently available.","mrkdwn_in":["text"]}]}
```
