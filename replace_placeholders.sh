#!/bin/bash

TASKDEF_JSON="taskdef.json"

sed -i "s/\${ACCOUNT_ID}/$ACCOUNT_ID/g" $TASKDEF_JSON
sed -i "s/\${AWS_REGION}/$AWS_REGION/g" $TASKDEF_JSON
