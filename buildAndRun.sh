#!/bin/bash

# Check if Go is installed
if ! [ -x "$(command -v go)" ]; then
    echo "Go is not installed. Exiting."
    exit 1
fi

# Remove the old binary
echo "Removing the old binary..."
rm -f ./bin/Sudoku.exe

# Build the Go application
echo "Building the application..."

if ! go build -o ./bin/Sudoku.exe; then
    echo "Build failed. Exiting."
    exit 1
fi

# Run the application
echo "Running the application..."
./bin/Sudoku.exe -debug
