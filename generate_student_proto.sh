#!/bin/bash

# Navigate to project root
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken

echo "Generating protobuf files for Student module..."

# Generate protobuf files
make proto-gen

echo "Protobuf generation completed!"

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo "✅ Student module protobuf files generated successfully!"
else
    echo "❌ Error generating protobuf files. Please check the proto files for syntax errors."
    exit 1
fi
