#!/usr/bin/env bash
CONFIG_TOKEN=$(grep "^token:" config.yml | sed -e 's/token://' -e 's/ //g' -e 's/"//g')
TEXT=$(echo $1 | sed -e 's/ /%20/')
CHANNEL=BBBB2222
if [[ -n $2 ]]; then
  CHANNEL=${2}
fi
KMS=$(grep "^kms:" config.yml | sed -e 's/kms://' -e 's/ //g' -e 's/"//g')
if [[ ${KMS} == "true" ]]; then
  echo "KMS doesn't work while testing, please disable it"
  exit
fi
IGOR_VAR="'{\"body\":\"token=${CONFIG_TOKEN}&team_id=CCCC3333&team_domain=testdomain&channel_id=${CHANNEL}&channel_name=igor-testing&user_id=AAAA1111&user_name=testuser&command=%2Figor&text=${TEXT}&response_url=http://slackhook\" }'"
/usr/bin/env bash -c "./igor ${IGOR_VAR}"
