#!/bin/bash

if ! [ -x "$(command -v templ)" ]; then
    echo "Templ is not installed. Exiting."
    pause 10
    exit 1
fi

if ! [ -x "$(command -v go)" ]; then
    echo "Go is not installed. Exiting."
    pause 10
    exit 1
fi

if ! [ -x "$(command -v air)" ]; then
    echo "Air is not installed. Exiting."
    pause 10
    exit 1
fi

if ! [ -x "$(command -v npm)" ]; then
    echo "NPM is not installed. Exiting."
    pause 10
    exit 1
fi

FILE="./.air.toml"

if [ -f "$FILE" ]; then
    echo "File '$FILE' exists."
else
    echo "File '$FILE' does not exist."
    exit 1
fi

air