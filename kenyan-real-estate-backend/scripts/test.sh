#!/bin/bash

# Test script for Kenyan Real Estate Backend
# This script tests basic functionality of the API

echo "🏠 Kenyan Real Estate Backend - Test Script"
echo "==========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
API_BASE_URL="http://localhost:8080/api/v1"
SERVER_PID=""

# Function to print colored output
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Function to start the server
start_server() {
    print_info "Starting the server..."
    
    # Check if server is already running
    if curl -s "$API_BASE_URL/../health" > /dev/null 2>&1; then
        print_status 0 "Server is already running"
        return 0
    fi
    
    # Start the server in background
    go run cmd/server/main.go &
    SERVER_PID=$!
    
    # Wait for server to start
    for i in {1..30}; do
        if curl -s "$API_BASE_URL/../health" > /dev/null 2>&1; then
            print_status 0 "Server started successfully (PID: $SERVER_PID)"
            return 0
        fi
        sleep 1
    done
    
    print_status 1 "Failed to start server"
    return 1
}

# Function to stop the server
stop_server() {
    if [ ! -z "$SERVER_PID" ]; then
        print_info "Stopping server (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
        print_status 0 "Server stopped"
    fi
}

# Function to test health endpoint
test_health() {
    print_info "Testing health endpoint..."
    
    response=$(curl -s -w "%{http_code}" "$API_BASE_URL/../health")
    http_code="${response: -3}"
    body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        print_status 0 "Health endpoint returned 200 OK"
        echo "Response: $body"
        return 0
    else
        print_status 1 "Health endpoint failed (HTTP $http_code)"
        return 1
    fi
}

# Function to test public endpoints
test_public_endpoints() {
    print_info "Testing public endpoints..."
    
    # Test counties endpoint
    response=$(curl -s -w "%{http_code}" "$API_BASE_URL/counties")
    http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        print_status 0 "Counties endpoint works"
    else
        print_status 1 "Counties endpoint failed (HTTP $http_code)"
    fi
    
    # Test amenities endpoint
    response=$(curl -s -w "%{http_code}" "$API_BASE_URL/amenities")
    http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        print_status 0 "Amenities endpoint works"
    else
        print_status 1 "Amenities endpoint failed (HTTP $http_code)"
    fi
    
    # Test property types endpoint
    response=$(curl -s -w "%{http_code}" "$API_BASE_URL/property-types")
    http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        print_status 0 "Property types endpoint works"
    else
        print_status 1 "Property types endpoint failed (HTTP $http_code)"
    fi
}

# Function to test phone validation
test_phone_validation() {
    print_info "Testing phone number validation..."
    
    # Test valid Kenyan phone number
    response=$(curl -s -w "%{http_code}" -X POST "$API_BASE_URL/validate-phone" \
        -H "Content-Type: application/json" \
        -d '{"phone_number": "0712345678"}')
    http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        print_status 0 "Phone validation endpoint works"
    else
        print_status 1 "Phone validation endpoint failed (HTTP $http_code)"
    fi
}

# Function to test user registration
test_user_registration() {
    print_info "Testing user registration..."
    
    # Generate random email to avoid conflicts
    random_email="test$(date +%s)@example.com"
    
    response=$(curl -s -w "%{http_code}" -X POST "$API_BASE_URL/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$random_email\",
            \"password\": \"testpassword123\",
            \"first_name\": \"Test\",
            \"last_name\": \"User\",
            \"phone_number\": \"0712345678\",
            \"user_type\": \"tenant\"
        }")
    http_code="${response: -3}"
    body="${response%???}"
    
    if [ "$http_code" = "201" ]; then
        print_status 0 "User registration works"
        # Extract token for further tests
        USER_TOKEN=$(echo "$body" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        if [ ! -z "$USER_TOKEN" ]; then
            print_status 0 "JWT token received"
        fi
    else
        print_status 1 "User registration failed (HTTP $http_code)"
        echo "Response: $body"
    fi
}

# Function to test protected endpoints
test_protected_endpoints() {
    if [ -z "$USER_TOKEN" ]; then
        print_info "Skipping protected endpoint tests (no token available)"
        return
    fi
    
    print_info "Testing protected endpoints..."
    
    # Test profile endpoint
    response=$(curl -s -w "%{http_code}" "$API_BASE_URL/profile" \
        -H "Authorization: Bearer $USER_TOKEN")
    http_code="${response: -3}"
    
    if [ "$http_code" = "200" ]; then
        print_status 0 "Profile endpoint works"
    else
        print_status 1 "Profile endpoint failed (HTTP $http_code)"
    fi
}

# Function to run compilation test
test_compilation() {
    print_info "Testing Go compilation..."
    
    if go build -o /tmp/kenyan-real-estate-test ./cmd/server; then
        print_status 0 "Code compiles successfully"
        rm -f /tmp/kenyan-real-estate-test
    else
        print_status 1 "Compilation failed"
        return 1
    fi
}

# Function to check dependencies
check_dependencies() {
    print_info "Checking dependencies..."
    
    # Check if Go is installed
    if command -v go &> /dev/null; then
        print_status 0 "Go is installed ($(go version))"
    else
        print_status 1 "Go is not installed"
        return 1
    fi
    
    # Check if curl is available
    if command -v curl &> /dev/null; then
        print_status 0 "curl is available"
    else
        print_status 1 "curl is not available"
        return 1
    fi
    
    # Check if go.mod exists
    if [ -f "go.mod" ]; then
        print_status 0 "go.mod file exists"
    else
        print_status 1 "go.mod file not found"
        return 1
    fi
}

# Function to run all tests
run_tests() {
    echo ""
    echo "🧪 Running Tests..."
    echo "=================="
    
    # Check dependencies first
    check_dependencies || return 1
    
    # Test compilation
    test_compilation || return 1
    
    # Start server
    start_server || return 1
    
    # Wait a moment for server to fully initialize
    sleep 2
    
    # Run API tests
    test_health
    test_public_endpoints
    test_phone_validation
    test_user_registration
    test_protected_endpoints
    
    # Stop server
    stop_server
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTION]"
    echo ""
    echo "Options:"
    echo "  test      Run all tests"
    echo "  start     Start the server"
    echo "  stop      Stop the server"
    echo "  health    Check health endpoint"
    echo "  help      Show this help message"
    echo ""
}

# Trap to ensure server is stopped on script exit
trap stop_server EXIT

# Main script logic
case "${1:-test}" in
    "test")
        run_tests
        ;;
    "start")
        start_server
        echo "Server is running. Press Ctrl+C to stop."
        wait
        ;;
    "stop")
        stop_server
        ;;
    "health")
        test_health
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    *)
        echo "Unknown option: $1"
        show_usage
        exit 1
        ;;
esac

echo ""
echo "🏁 Test script completed!"

