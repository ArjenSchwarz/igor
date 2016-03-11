# Igor

Igor is a Slack Slash command that acts like a bot. The code is written in Go and is designed to be run on AWS Lambda through a NodeJS wrapper. All commands are handled through plugins, making it extendable.

The name is based on Sir Terry Pratchett's wonderful use of Dr. Frankenstein's servant. No disrespect intended.

Igor is currently early in development, and can't do much yet, but it is usable.

# Build status

[![wercker status](https://app.wercker.com/status/eea144a7251e1b84d514904e19eff205/m "wercker status")](https://app.wercker.com/project/bykey/eea144a7251e1b84d514904e19eff205)

# Available Plugins

* Help
* Weather
* (Random) Tumblr image

# TODO

Many things, including but not limited to:

* Better documentation
* More plugins
* More tests

# Installation

The installation is slightly messy, as both Slack and Lambda require details that can only be provided after you're done. To provide the easiest step by step instructions, this therefore means you'll be switching between the two.

## Configure Slack

Follow these steps to configure the slash command in Slack:

1. Navigate to https://<your-team-domain>.slack.com/services/new
2. Search for and select "Slash Commands".
3. Enter a name for your command (I recommend igor) and click "Add Slash Command Integration".
4. Copy the token string from the integration settings and use it in the next section.
5. Leave the page open for now

## Download and configure

1. Download the [latest igor zip file](https://github.com/ArjenSchwarz/igor/releases/download/latest/igor.zip)
2. Unzip this file
3. Edit the *config.yml* file in there by replacing the placeholder token string with yours.
4. Make any other configuration changes you wish to make.
5. Zip it again, making sure all 3 files are in the new zip file.

## Set up Lambda

1. In your AWS Console, create a new Lambda function. When asked for a blueprint, skip this (the option is at the bottom).
2. Provide the name and description to the project, and ensure that the Runtime is set to NodeJS. 
3. For the source you can upload the zipfile you just created. 
4. Set the role. If you don't have one yet, you can select to create a Basic Execution Role and use that. Currently no other permissions are required.
5. All other settings you can leave at their default values, and you can continue to create the function.

## API endpoint

At this point you will be brought to an overview of the function. Here you will need to select the API endpoints tab and then follow the remaining steps.

1. Click on add a new API endpoint. 
2. Choose the API Gateway option and a resource name that you're happy with. Also, make sure to set the Method to **POST**, and the Security to **OPEN**.

You are then returned to the function's API endpoint overview. 

3. Make a note of the API endpoint as we'll need it later.
4. Click on the *prod* link to configure the remaining details.
5. Click on the link for Resources and then select the POST under your endpoint's name (for example /igor). 
6. Select the Integration Request, and add a new Mapping Template. The Content-Type for this should be: `application/x-www-form-urlencoded`.
7. After you've created this, change it from *Input Passthrough* to *Mapping template* and use `{ "body": $input.json("$") }` as the mapping template.

You've now made changes to the API, so you will have to deploy it again. There is a button for that. When deploying, make sure to deploy to the *prod* environment.

## Slack

On the page you left open, fill in the URL for the API endpoint and configure everything else as you see fit before you save the integration.

## Try it out

In Slack (presuming you chose the trigger igor) you can now run **/igor help** to get an overview of what it can do.

# Contribute

If you wish to contribute in any way (reporting bugs, requesting features, writing code), feel free to do so either by opening Issues or Pull Request. For Pull Requests, just follow the standard pattern.

1. Fork the repository
2. Make your changes
3. Make a pull request that explains what it does

To make plugin development easier, there is a snippet for Sublime Text included in the devtools directory. If you copy this to your User package you can easily create the skeleton for a plugin with it.