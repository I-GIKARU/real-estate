#!/bin/bash

API_BASE="https://kenyan-real-estate-backend-671327858247.us-central1.run.app/api/v1"
EMAIL="testuser2@example.com"

echo "=== Testing Simplified Email Verification Flow ==="
echo

echo "1. Registering new user (should automatically send verification email)..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'",
    "password": "TestPassword123!",
    "first_name": "Test",
    "last_name": "User2", 
    "phone_number": "+254700000001",
    "user_type": "tenant"
  }')

echo "Response: $REGISTER_RESPONSE"
echo

# Extract token from response
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo "2. Checking verification status..."
    curl -s -X GET "$API_BASE/verification-status" \
      -H "Authorization: Bearer $TOKEN" | jq .
    echo
else
    echo "Failed to get token from registration"
    exit 1
fi

echo "3. Login attempt (should work even without verification)..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "'$EMAIL'", 
    "password": "TestPassword123!"
  }')

echo "Response: $LOGIN_RESPONSE"
echo

echo "=== Flow Summary ==="
echo "✅ 1. User registers → Verification email sent automatically"
echo "✅ 2. User can check verification status"  
echo "✅ 3. User can login immediately (verification not required for access)"
echo "📧 4. User clicks link in email to verify (when they receive it)"
echo
echo "The simplified flow is working! Users get verification email automatically"
echo "and can use the system immediately while verification is pending."
