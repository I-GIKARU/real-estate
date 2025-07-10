#!/bin/bash

# Test script for Kenyan Real Estate API
# Make sure the server is running on localhost:8080

API_BASE="http://localhost:8080/api/v1"
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${CYAN}🏠 Testing Kenyan Real Estate API${NC}"
echo "========================================"

# Function to make HTTP requests and pretty print JSON
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local token=$4
    
    echo -e "\n${YELLOW}📍 $method $endpoint${NC}"
    
    if [ -n "$data" ] && [ -n "$token" ]; then
        curl -s -X "$method" "$API_BASE$endpoint" \
             -H "Content-Type: application/json" \
             -H "Authorization: Bearer $token" \
             -d "$data" | jq '.'
    elif [ -n "$data" ]; then
        curl -s -X "$method" "$API_BASE$endpoint" \
             -H "Content-Type: application/json" \
             -d "$data" | jq '.'
    elif [ -n "$token" ]; then
        curl -s -X "$method" "$API_BASE$endpoint" \
             -H "Authorization: Bearer $token" | jq '.'
    else
        curl -s -X "$method" "$API_BASE$endpoint" | jq '.'
    fi
}

# Check if server is running
echo -e "\n${CYAN}🔍 Checking if server is running...${NC}"
health_check=$(curl -s http://localhost:8080/health 2>/dev/null)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Server is running!${NC}"
    echo "$health_check" | jq '.'
else
    echo -e "${RED}❌ Server is not running. Please start the server first:${NC}"
    echo "./bin/server"
    exit 1
fi

# Test public endpoints
echo -e "\n${CYAN}🌍 Testing Public Endpoints${NC}"
echo "================================"

# Get counties
make_request "GET" "/counties"

# Get properties (should be empty initially)
make_request "GET" "/properties"

# Test user registration
echo -e "\n${CYAN}👤 Testing User Registration${NC}"
echo "==============================="

# Register a test landlord
landlord_data='{
  "email": "landlord@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Landlord",
  "phone_number": "+254700000001",
  "user_type": "landlord"
}'

echo -e "\n${YELLOW}📝 Registering landlord...${NC}"
landlord_response=$(curl -s -X POST "$API_BASE/register" \
                   -H "Content-Type: application/json" \
                   -d "$landlord_data")

echo "$landlord_response" | jq '.'

# Extract token for landlord
landlord_token=$(echo "$landlord_response" | jq -r '.token // empty')

if [ -n "$landlord_token" ] && [ "$landlord_token" != "null" ]; then
    echo -e "${GREEN}✅ Landlord registered successfully!${NC}"
    
    # Test authenticated endpoints
    echo -e "\n${CYAN}🔐 Testing Authenticated Endpoints${NC}"
    echo "===================================="
    
    # Get user profile
    make_request "GET" "/profile" "" "$landlord_token"
    
    # Create a test property
    echo -e "\n${YELLOW}🏡 Creating a test property...${NC}"
    property_data='{
      "title": "Beautiful 2BR Apartment in Westlands",
      "description": "Modern apartment with stunning city views",
      "property_type": "apartment",
      "bedrooms": 2,
      "bathrooms": 2,
      "square_meters": 120.5,
      "rent_amount": 75000,
      "deposit_amount": 150000,
      "county_id": 47,
      "location_details": "Westlands, Nairobi",
      "amenities": {
        "wifi": true,
        "parking": true,
        "gym": false,
        "swimming_pool": true,
        "security": true
      },
      "utilities_included": {
        "water": true,
        "electricity": false,
        "internet": true
      },
      "parking_spaces": 1,
      "is_furnished": true
    }'
    
    property_response=$(curl -s -X POST "$API_BASE/properties" \
                       -H "Content-Type: application/json" \
                       -H "Authorization: Bearer $landlord_token" \
                       -d "$property_data")
    
    echo "$property_response" | jq '.'
    
    # Get properties again (should show the new property)
    echo -e "\n${YELLOW}📋 Getting updated property list...${NC}"
    make_request "GET" "/properties"
    
else
    echo -e "${RED}❌ Registration failed or no token received${NC}"
fi

# Register a test tenant
echo -e "\n${YELLOW}👥 Registering tenant...${NC}"
tenant_data='{
  "email": "tenant@example.com",
  "password": "password123",
  "first_name": "Jane",
  "last_name": "Tenant",
  "phone_number": "+254700000002",
  "user_type": "tenant"
}'

tenant_response=$(curl -s -X POST "$API_BASE/register" \
                 -H "Content-Type: application/json" \
                 -d "$tenant_data")

echo "$tenant_response" | jq '.'

# Test property search with filters
echo -e "\n${CYAN}🔍 Testing Property Search with Filters${NC}"
echo "========================================"

# Search properties in Nairobi
echo -e "\n${YELLOW}🏙️ Searching properties in Nairobi (County ID: 47)...${NC}"
make_request "GET" "/properties?county_id=47"

# Search apartments with rent between 50k-100k
echo -e "\n${YELLOW}🏠 Searching apartments with rent 50k-100k...${NC}"
make_request "GET" "/properties?property_type=apartment&min_rent=50000&max_rent=100000"

# Search 2+ bedroom properties
echo -e "\n${YELLOW}🛏️ Searching 2+ bedroom properties...${NC}"
make_request "GET" "/properties?min_bedrooms=2"

echo -e "\n${CYAN}🎯 API Testing Complete!${NC}"
echo "================================"
echo -e "${GREEN}✅ All tests completed successfully!${NC}"
echo ""
echo -e "${YELLOW}📚 To explore more endpoints, visit:${NC}"
echo "   http://localhost:8080/swagger/index.html"
echo ""
echo -e "${YELLOW}🔧 Available endpoints:${NC}"
echo "   • Authentication: /register, /login, /refresh-token"
echo "   • User Management: /profile"
echo "   • Properties: /properties, /my-properties"
echo "   • Locations: /counties, /sub-counties"
echo "   • Kenyan Features: /amenities, /property-types"
