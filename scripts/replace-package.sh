#!/usr/bin/env bash

TAG="[replace-package.sh]"

echo ""
read -p "${TAG} previous module name: " PREV_MODULE_NAME
read -p "${TAG} new module name: " NEW_MODULE_NAME

if [ "$PREV_MODULE_NAME" == "" ] || [ "$NEW_MODULE_NAME" == "" ]; then
  echo -e "${TAG} ERROR: Got empty module name \n"
  exit 0
fi

echo -e "${TAG} Converting \`${PREV_MODULE_NAME}\` to \`${NEW_MODULE_NAME}\` \n"

find . -type f -name "*.go" | grep -v vendor |      xargs sed -i '' "s/${PREV_MODULE_NAME}/${NEW_MODULE_NAME}/g"
find . -type f -name "Makefile" | grep -v vendor |  xargs sed -i '' "s/${PREV_MODULE_NAME}/${NEW_MODULE_NAME}/g"


