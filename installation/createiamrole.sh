#!/usr/bin/env bash
set -e
ROLE_ARN=$(aws iam list-roles --query 'Roles[?RoleName==`IgorRole`].Arn' --output text)
if ! [ -z "$ROLE_ARN" ]
then
    echo "An existing role called IgorRole was found. The ARN for this is:
${ROLE_ARN}"
else
    aws iam create-role --role-name "IgorRole" --assume-role-policy-document file://iamtrustdocument.json
    aws iam put-role-policy --role-name "IgorRole" --policy-name "IgorRolePolicy" --policy-document file://basiciamrole.json
    ROLE_ARN=$(aws iam get-role --role-name "IgorRole" --query Role.Arn --output text)
    echo "Created a new role called IgorRole. The ARN for this is:
${ROLE_ARN}"
fi
