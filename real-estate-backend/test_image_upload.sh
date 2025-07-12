#!/bin/bash

# Test script for image upload endpoint
# Replace these variables with your actual values

PROPERTY_ID="f476792c-3664-412e-89b5-c36b7954fa6b"
BASE_URL="https://real-estate-backend-840370620772.us-central1.run.app/api/v1"
TOKEN="your-auth-token-here"

# Create a test image file
echo "Creating test image file..."
convert -size 300x200 xc:lightblue test_image.jpg 2>/dev/null || {
    echo "ImageMagick not available, creating a simple text file as test..."
    echo "This is a test image file" > test_image.txt
    cp test_image.txt test_image.jpg
}

# Test the upload endpoint
echo "Testing image upload endpoint..."
curl -X POST \
  "${BASE_URL}/properties/${PROPERTY_ID}/images" \
  -H "Authorization: Bearer ${TOKEN}" \
  -F "image=@test_image.jpg" \
  -F "caption=Test image upload" \
  -F "is_primary=false" \
  -F "display_order=1" \
  -v

# Clean up
rm -f test_image.jpg test_image.txt

echo "Test completed."
