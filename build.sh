#!/bin/bash

# Clean and build the project
echo "Cleaning previous build..."
make clean

echo "Building the project..."
make install

echo "Build complete!"
