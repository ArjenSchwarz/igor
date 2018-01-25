#!/usr/bin/env bash
set -e
ARN_ROLE=$1
if [ -z "$ARN_ROLE" ]
then
    echo "You need to provide the ARN for the Role you wish to use. Try running createiamrole.sh if you don't have one."
    exit
fi
if [ -z "$3" ]
then
    REGION="us-east-1"
else
    REGION=$3
fi
if [ -z "$2" ]
then
    NAME="igor"
else
    NAME=$2
fi
APINAME="${NAME}Api"

# Create the Lambda function
echo "Creating the Lambda function"
aws lambda create-function --function-name "${NAME}" --runtime go1.x --role ${ARN_ROLE} --handler main --zip-file fileb://igor.zip --region ${REGION}
LAMBDAARN=$(aws lambda list-functions --query "Functions[?FunctionName==\`${NAME}\`].FunctionArn" --output text --region ${REGION})

# Create the API Gateway
echo "Creating the API Gateway"
aws apigateway create-rest-api --name "${APINAME}" --description "Api for ${NAME}" --region ${REGION}
APIID=$(aws apigateway get-rest-apis --query "items[?name==\`${APINAME}\`].id" --output text --region ${REGION})
PARENTRESOURCEID=$(aws apigateway get-resources --rest-api-id ${APIID} --query 'items[?path==`/`].id' --output text --region ${REGION})

# Add the resource
aws apigateway create-resource --rest-api-id ${APIID} --parent-id ${PARENTRESOURCEID} --path-part igor --region ${REGION}
RESOURCEID=$(aws apigateway get-resources --rest-api-id ${APIID} --query 'items[?path==`/igor`].id' --output text --region ${REGION})

# Add the POST method
aws apigateway put-method --rest-api-id ${APIID} --resource-id ${RESOURCEID} --http-method POST --authorization-type NONE --region ${REGION}

# Method request config
aws apigateway put-integration --rest-api-id ${APIID} \
--resource-id ${RESOURCEID} \
--http-method POST \
--type AWS \
--integration-http-method POST \
--uri arn:aws:apigateway:${REGION}:lambda:path/2015-03-31/functions/${LAMBDAARN}/invocations \
--request-templates '{"application/x-www-form-urlencoded":"{\"body\": $input.json(\"$\")}"}' \
--region ${REGION}

# Method response config
aws apigateway put-method-response \
--rest-api-id ${APIID} \
--resource-id ${RESOURCEID} \
--http-method POST \
--status-code 200 \
--response-models "{}" \
--region ${REGION}

aws apigateway put-integration-response \
--rest-api-id ${APIID} \
--resource-id ${RESOURCEID} \
--http-method POST \
--status-code 200 \
--selection-pattern ".*" \
--region ${REGION}

# Deploy Gateway
aws apigateway create-deployment \
--rest-api-id ${APIID} \
--stage-name prod \
--region ${REGION}

# Create permissions
APIARN=$(echo ${LAMBDAARN} | sed -e 's/lambda/execute-api/' -e "s/function:${NAME}/${APIID}/")
aws lambda add-permission \
--function-name ${NAME} \
--statement-id apigateway-igor-test-2 \
--action lambda:InvokeFunction \
--principal apigateway.amazonaws.com \
--source-arn "${APIARN}/*/POST/igor" \
--region ${REGION}

aws lambda add-permission \
--function-name ${NAME} \
--statement-id apigateway-igor-prod-2 \
--action lambda:InvokeFunction \
--principal apigateway.amazonaws.com \
--source-arn "${APIARN}/prod/POST/igor" \
--region ${REGION}

echo "The url you have to use in your Slack settings is:
https://${APIID}.execute-api.${REGION}.amazonaws.com/prod/igor"
