# Installing Igor

Both Slack and Lambda require details from each other. To provide the easiest step by step instructions, this therefore means you'll be switching between the two.

# Configure Slack

Follow these steps to configure the slash command in Slack:

1. Navigate to https://your-team-domain.slack.com/services/new
2. Search for and select "Slash Commands".
3. Enter a name for your command (I recommend igor) and click "Add Slash Command Integration".
4. Copy the token string from the integration settings and use it in the next section.
5. Leave the page open for now

# Download and configure

1. Download the [latest igor zip file](https://github.com/ArjenSchwarz/igor/releases/download/latest/igor.zip).
2. Unzip this file.
3. Edit the *config.yml* file in there by replacing the placeholder token string with yours.
4. Make any other configuration changes you wish to make.
5. Zip it again, making sure all 3 files are in the new zip file.

# Set up AWS

A bash script is provided to automate this part. This is a lot easier, but the manual steps are described as well in case you are unable to use this.

## Automated

The automated installation requires that you have the [AWS CLI tools](https://aws.amazon.com/cli/) installed, and assumes that the zipfile of the source is in the same directory and named `igor.zip`.

1. Download the [latest installation scripts](https://github.com/ArjenSchwarz/igor/releases/download/latest/installation.zip).
2. Unzip.
3. Either run the `createiamrole.sh` script and copy the resulting ARN or copy the ARN from a role you wish to use.
4. Run the `setupaws.sh` script, providing the role ARN and optionally a name and region. It will default to use "igor" and "us-east-1". The URL you need for the next step will be shown at the end of the script.

Example:
```bash
./setupaws.sh "arn:aws:iam::ACCOUNT_ID:role/IgorRole" "igor" "us-east-1"
```

## Manually

### Lambda

1. In your AWS Console, create a new Lambda function. When asked for a blueprint, skip this (the option is at the bottom).
2. Provide the name and description to the project, and ensure that the Runtime is set to Go1.x and the handler is set to `main`. 
3. For the source you can upload the zipfile you just created.
4. Set the role. If you don't have one yet, you can select to create a Basic Execution Role and use that. Currently no other permissions are required.
5. All other settings you can leave at their default values, and you can continue to create the function.

### API endpoint

At this point you will be brought to an overview of the function. Here you will need to select the API endpoints tab and then follow the remaining steps.

1. Click on add a new API endpoint.
2. Choose the API Gateway option and a resource name that you're happy with. Also, make sure to set the Method to **POST**, and the Security to **OPEN**.

You are then returned to the function's API endpoint overview.

3. Make a note of the API endpoint as we'll need it later.
4. Click on the *prod* link to configure the remaining details.
5. Click on the link for Resources and then select the POST under your endpoint's name (for example /igor).
6. Select the Integration Request, and add a new Mapping Template. The Content-Type for this should be: `application/x-www-form-urlencoded`.
7. After you've created this template, change it from *Input Passthrough* to *Mapping template* and use `{ "body": $input.json("$") }` as the mapping template.

You've now made changes to the API, so you will have to deploy it again. There is a button for that. When deploying, make sure to deploy to the *prod* environment.

# Finish Slack

On the page you left open, fill in the URL for the API endpoint and configure everything else as you see fit before you save the integration.

# Try it out

In Slack (presuming you chose the trigger igor) you can now run **/igor help** to get an overview of what it can do.
