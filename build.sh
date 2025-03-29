#!/usr/bin/env bash

CGO_ENABLED=1 
GOOS=linux 
GOARCH=amd64 

if [[ ! -d ./dist ]]
then
    echo "Making dist folder..."
    mkdir -p ./dist
fi

for app in ./cmd/*/ ; do
    app_name=$(basename $app)
    echo "Building $app_name with additional tag '$ADDITIONAL_TAGS'..."
    go build $ADDITIONAL_TAGS -tags netgo -ldflags '-w -extldflags "-static"' -o ./dist/$app_name ./cmd/$app_name
done