#!/bin/bash

# Enable API and start the chain
echo "Starting Academic Token blockchain with API enabled..."

# Navigate to blockchain directory
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken

# Start chain with API enabled
ignite chain serve \
  --reset-once \
  --api.enable \
  --api.swagger \
  --api.address tcp://0.0.0.0:1318 \
  --api.enabled-unsafe-cors
