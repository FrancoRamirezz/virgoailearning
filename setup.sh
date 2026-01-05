#!/bin/bash

echo "Setting up the Backend project..."

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Check if dependencies were downloaded successfully
if [ $? -eq 0 ]; then
    echo "Dependencies downloaded successfully!"
else
    echo "Failed to download dependencies"
    exit 1
fi

# Build the project
echo "Building the project..."
go build -o backend .

if [ $? -eq 0 ]; then
    echo "Project built successfully!"
    echo "You can now run: ./backend"
else
    echo "Build failed"
    exit 1
fi
