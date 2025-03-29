#!/usr/bin/env bash

APP_NAME=$1
APP_PARAMS=$2

if [ -z $APP_NAME ]; then
    echo 'Usage: ./dev.sh <APP_NAME> [<APP_PARAMS>]'
    exit 1
fi

if ! [ -x "$(command -v CompileDaemon)" ]; then
    echo '---------------------------------------------------'
    echo '> CompileDaemon is not installed.                 <'
    echo '> Run the following command to install the binary <'
    echo '> go get github.com/githubnemo/CompileDaemon      <'
    echo '---------------------------------------------------'
    exit 1
fi

function cleanup() {
    rm ./dist/${APP_NAME};
    exit
}
trap cleanup INT SIGHUP SIGINT SIGTERM

# lsof -i :8080 | awk '{print $2}' | grep -v "PID" | xargs kill -9

wire ./cmd/${APP_NAME}/
CompileDaemon -log-prefix=false  \
    -build="go build -x -o ../../dist/${APP_NAME}"   \
    -build-dir="./cmd/${APP_NAME}"  \
    -command="./dist/${APP_NAME} ${APP_PARAMS}"  \
    -pattern="(.+.go|.+.env)$"  \
    -exclude-dir=".git"  \
    -exclude-dir=".idea"  \
    -exclude-dir="vendor"  \
    -exclude-dir="dist"  \
    -exclude-dir="cmd"  \
    -graceful-kill=true  \
    -graceful-timeout=10 \
    -color

