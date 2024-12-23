#!/bin/bash

if ! [ -x "$(command -v templ)" ]; then
    echo "Templ is not installed. Exiting."
    pause 10
    exit 1
fi

echo "Generating templated files..."

if ! templ generate; then
    echo "Templ failed. Exiting."
    pause 10
    exit 1
fi

# Check if Go is installed
if ! [ -x "$(command -v go)" ]; then
    echo "Go is not installed. Exiting."
    pause 10
    exit 1
fi

# Remove the old binary
echo "Removing the old binary..."
rm -f ./bin/Sudoku.exe

# Build the Go application
echo "Building the application..."

if ! go build -o ./bin/Sudoku.exe; then
    echo "Build failed. Exiting."
    pause 10
    exit 1
fi

# Run the application
echo "Running the application..."
./bin/Sudoku.exe -debug

pause 10