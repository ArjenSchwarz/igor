#!/usr/bin/env bash
CONFIG_TOKEN=$(grep "^token:" config.yml | sed -e 's/token://' -e 's/ //g' -e 's/"//g')
TEXT=$(echo $1 | sed 's/ /%20/')
IGOR_VAR="'{\"body\":\"token=${CONFIG_TOKEN}&team_id=CCCC3333&team_domain=testdomain&channel_id=BBBB2222&channel_name=igor-testing&user_id=AAAA1111&user_name=testuser&command=%2Figor&text=${TEXT}&response_url=http://slackhook\" }'"
/usr/bin/env bash -c "./igor ${IGOR_VAR}"