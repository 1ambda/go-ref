#!/usr/bin/env bash

VERSION=$(cat VERSION)
SWAGGER_JSON="api/swagger.yml"

sed -i.bak s/[\s]*version:.*/version:\ ${VERSION}/g ${SWAGGER_JSON}
rm ${SWAGGER_JSON}.bak || true
