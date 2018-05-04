#!/usr/bin/env bash

TAG="[get-ldflags.sh]"

USAGE="Usage: WORKSPACE={RELATIVE_WORKSPACE_PATH} CONFIG_PACKAGE_PATH=${CONFIG_PKG_NAME} $0"


if [ "$WORKSPACE" == "" ] || [ "$CONFIG_PACKAGE_PATH" == "" ]; then
    echo ${USAGE} >&2
    exit 1
fi

WORKSPACE=$(echo ${WORKSPACE} | tr -d " ")
CONFIG_PACKAGE_PATH=$(echo ${CONFIG_PACKAGE_PATH} | tr -d " ")

if ! [ -x "$(command -v bazel)" ]; then
  echo 'Error: buildozer is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v buildozer)" ]; then
  echo 'Error: buildozer is not installed.' >&2
  exit 1
fi

LDFLAG_VERSION=$(cat ./VERSION | tr -d " ")
LDFLAG_BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAG_GIT_COMMIT=$(git rev-parse --short HEAD | tr -d " ")
LDFLAG_GIT_BRANCH=$(git symbolic-ref -q --short HEAD | tr -d " ")
LDFLAG_GIT_STATE="dirty"

DIRTY_FILES=$(git status --porcelain | grep -v package-lock.json | wc -l | tr -d " ")
if [ "$DIRTY_FILES" == "0" ]; then
  LDFLAG_GIT_STATE="clean"
fi

RULES=$(bazel query "kind(go_library, //${CONFIG_PACKAGE_PATH}:go_default_library)")
buildozer -quiet "remove x_defs" $RULES
DICT="{\"Version\":\ \"${LDFLAG_VERSION}\""
DICT="${DICT},\ \"BuildDate\":\ \"${LDFLAG_BUILD_DATE}\""
DICT="${DICT},\ \"GitCommit\":\ \"${LDFLAG_GIT_COMMIT}\""
DICT="${DICT},\ \"GitBranch\":\ \"${LDFLAG_GIT_BRANCH}\""
DICT="${DICT},\ \"GitState\":\ \"${LDFLAG_GIT_STATE}\""
DICT="${DICT}}"
buildozer -quiet "set x_defs ${DICT}" $RULES 1> /dev/null

