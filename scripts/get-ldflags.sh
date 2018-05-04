#!/usr/bin/env bash

export LDFLAG_GIT_COMMIT=$(git rev-parse --short HEAD)
export LDFLAG_GIT_STATE

echo $LDFLAG_GIT_COMMIT
