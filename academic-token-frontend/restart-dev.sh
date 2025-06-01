#!/bin/bash

echo "🔄 Restarting Next.js to clear cache..."

# Kill any existing Next.js processes
pkill -f "next dev" || true
pkill -f "next-server" || true

# Clear Next.js cache
rm -rf .next

echo "✅ Cache cleared. Starting fresh development server..."
npm run dev
