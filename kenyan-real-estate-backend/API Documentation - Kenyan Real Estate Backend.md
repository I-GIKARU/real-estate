# API Documentation - Kenyan Real Estate Backend

## Overview

This document provides detailed information about the REST API endpoints for the Kenyan Real Estate Backend platform.

**Base URL**: `http://localhost:8080/api/v1`

**Authentication**: JWT Bearer Token (where required)

**Content-Type**: `application/json`

## Authentication

### Register User

Creates a new user account.

**Endpoint**: `POST /register`

**Request Body**:
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "0712345678",
  "user_type": "tenant"
}
```

**Response** (201 Created):
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid-here",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "254712345678",
    "user_type": "tenant",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "token": "jwt-token-here"
}
```

### Login

Authenticates a user and returns a JWT token.

**Endpoint**: `POST /login`

**Request Body**:
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Response** (200 OK):
```json
{
  "message": "Login successful",
  "user": {
    "id": "uuid-here",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "user_type": "tenant"
  },
  "token": "jwt-token-here"
}
```

### Refresh Token

Refreshes an existing JWT token.

**Endpoint**: `POST /refresh-token`

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "token": "new-jwt-token-here"
}
```

## User Profile

### Get Profile

Gets the authenticated user's profile.

**Endpoint**: `GET /profile`

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "user": {
    "id": "uuid-here",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "254712345678",
    "user_type": "tenant",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Update Profile

Updates the authenticated user's profile.

**Endpoint**: `PUT /profile`

**Headers**: `Authorization: Bearer <token>`

**Request Body**:
```json
{
  "first_name": "John",
  "last_name": "Smith",
  "phone_number": "0722345678"
}
```

## Properties

### Get Public Properties

Retrieves public property listings with optional filtering.

**Endpoint**: `GET /properties`

**Query Parameters**:
- `county_id` (integer): Filter by county ID
- `sub_county_id` (integer): Filter by sub-county ID
- `property_type` (string): Filter by property type
- `min_rent` (number): Minimum rent amount
- `max_rent` (number): Maximum rent amount
- `min_bedrooms` (integer): Minimum number of bedrooms
- `max_bedrooms` (integer): Maximum number of bedrooms
- `min_bathrooms` (integer): Minimum number of bathrooms
- `is_furnished` (boolean): Filter by furnished status
- `has_parking` (boolean): Filter by parking availability
- `limit` (integer): Number of results per page (default: 20)
- `offset` (integer): Number of results to skip (default: 0)

**Example**: `GET /properties?county_id=1&property_type=apartment&min_rent=20000&max_rent=80000&limit=10`

**Response** (200 OK):
```json
{
  "properties": [
    {
      "id": "uuid-here",
      "title": "Modern 2BR Apartment in Kilimani",
      "description": "Beautiful apartment with modern amenities",
      "property_type": "apartment",
      "bedrooms": 2,
      "bathrooms": 2,
      "square_meters": 120.5,
      "rent_amount": 45000,
      "deposit_amount": 90000,
      "county_id": 1,
      "sub_county_id": 1,
      "location_details": "Near Yaya Centre, Kilimani",
      "latitude": -1.2921,
      "longitude": 36.8219,
      "amenities": ["24/7 Security", "Swimming Pool", "Gym"],
      "utilities_included": ["water", "security"],
      "parking_spaces": 1,
      "is_furnished": true,
      "is_available": true,
      "availability_date": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "county": {
        "id": 1,
        "name": "Nairobi",
        "code": "001"
      },
      "sub_county": {
        "id": 1,
        "name": "Westlands",
        "county_id": 1
      },
      "images": [
        {
          "id": "uuid-here",
          "image_url": "https://example.com/image1.jpg",
          "caption": "Living room",
          "is_primary": true,
          "display_order": 1
        }
      ]
    }
  ],
  "filters": {
    "county_id": 1,
    "property_type": "apartment",
    "min_rent": 20000,
    "max_rent": 80000,
    "limit": 10,
    "offset": 0
  }
}
```

### Get Single Property

Retrieves details of a specific property.

**Endpoint**: `GET /properties/{id}`

**Response** (200 OK):
```json
{
  "property": {
    "id": "uuid-here",
    "title": "Modern 2BR Apartment in Kilimani",
    "description": "Beautiful apartment with modern amenities...",
    "property_type": "apartment",
    "bedrooms": 2,
    "bathrooms": 2,
    "square_meters": 120.5,
    "rent_amount": 45000,
    "deposit_amount": 90000,
    "county_id": 1,
    "sub_county_id": 1,
    "location_details": "Near Yaya Centre, Kilimani",
    "latitude": -1.2921,
    "longitude": 36.8219,
    "amenities": ["24/7 Security", "Swimming Pool", "Gym"],
    "utilities_included": ["water", "security"],
    "parking_spaces": 1,
    "is_furnished": true,
    "is_available": true,
    "availability_date": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "county": {
      "id": 1,
      "name": "Nairobi",
      "code": "001"
    },
    "sub_county": {
      "id": 1,
      "name": "Westlands",
      "county_id": 1
    },
    "landlord": {
      "id": "uuid-here",
      "first_name": "Jane",
      "last_name": "Smith",
      "phone_number": "254722345678"
    },
    "images": [
      {
        "id": "uuid-here",
        "image_url": "https://example.com/image1.jpg",
        "caption": "Living room",
        "is_primary": true,
        "display_order": 1
      }
    ]
  }
}
```

### Create Property (Landlord Only)

Creates a new property listing.

**Endpoint**: `POST /properties`

**Headers**: `Authorization: Bearer <landlord-token>`

**Request Body**:
```json
{
  "title": "Modern 2BR Apartment in Kilimani",
  "description": "Beautiful apartment with modern amenities, close to shopping centers and public transport.",
  "property_type": "apartment",
  "bedrooms": 2,
  "bathrooms": 2,
  "square_meters": 120.5,
  "rent_amount": 45000,
  "deposit_amount": 90000,
  "county_id": 1,
  "sub_county_id": 1,
  "location_details": "Near Yaya Centre, Kilimani",
  "latitude": -1.2921,
  "longitude": 36.8219,
  "amenities": ["24/7 Security", "Swimming Pool", "Gym", "Parking"],
  "utilities_included": ["water", "security", "garbage"],
  "parking_spaces": 1,
  "is_furnished": true,
  "availability_date": "2024-02-01T00:00:00Z"
}
```

**Response** (201 Created):
```json
{
  "message": "Property created successfully",
  "property": {
    "id": "uuid-here",
    "title": "Modern 2BR Apartment in Kilimani",
    "landlord_id": "uuid-here",
    "is_available": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### Update Property (Landlord Only)

Updates an existing property.

**Endpoint**: `PUT /properties/{id}`

**Headers**: `Authorization: Bearer <landlord-token>`

**Request Body** (partial update):
```json
{
  "rent_amount": 50000,
  "is_available": false,
  "amenities": ["24/7 Security", "Swimming Pool", "Gym", "Parking", "Generator"]
}
```

### Delete Property (Landlord Only)

Deletes a property listing.

**Endpoint**: `DELETE /properties/{id}`

**Headers**: `Authorization: Bearer <landlord-token>`

**Response** (200 OK):
```json
{
  "message": "Property deleted successfully"
}
```

### Get My Properties (Landlord Only)

Gets properties owned by the authenticated landlord.

**Endpoint**: `GET /my-properties`

**Headers**: `Authorization: Bearer <landlord-token>`

**Query Parameters**:
- `limit` (integer): Number of results per page (default: 20)
- `offset` (integer): Number of results to skip (default: 0)

### Add Property Image (Landlord Only)

Adds an image to a property.

**Endpoint**: `POST /properties/{id}/images`

**Headers**: `Authorization: Bearer <landlord-token>`

**Request Body**:
```json
{
  "image_url": "https://example.com/property-image.jpg",
  "caption": "Master bedroom",
  "is_primary": false,
  "display_order": 2
}
```

## Location Services

### Get Counties

Retrieves all Kenyan counties.

**Endpoint**: `GET /counties`

**Response** (200 OK):
```json
{
  "counties": [
    {
      "id": 1,
      "name": "Nairobi",
      "code": "001"
    },
    {
      "id": 2,
      "name": "Mombasa",
      "code": "002"
    }
  ]
}
```

### Get County

Retrieves a specific county.

**Endpoint**: `GET /counties/{id}`

### Get Sub-Counties

Retrieves sub-counties for a specific county.

**Endpoint**: `GET /counties/{county_id}/sub-counties`

**Response** (200 OK):
```json
{
  "sub_counties": [
    {
      "id": 1,
      "name": "Westlands",
      "county_id": 1
    },
    {
      "id": 2,
      "name": "Dagoretti North",
      "county_id": 1
    }
  ],
  "county_id": 1
}
```

## Kenyan Features

### Get Amenities

Retrieves property amenities by category or all amenities.

**Endpoint**: `GET /amenities`

**Query Parameters**:
- `category` (string): Filter by category (security, utilities, kitchen, etc.)

**Response** (200 OK):
```json
{
  "amenities": {
    "security": [
      "24/7 Security",
      "CCTV Surveillance",
      "Electric Fence",
      "Security Guards"
    ],
    "utilities": [
      "Borehole Water",
      "Mains Water",
      "Backup Generator",
      "Solar Water Heating"
    ]
  }
}
```

### Get Property Types

Retrieves all property types with descriptions.

**Endpoint**: `GET /property-types`

**Response** (200 OK):
```json
{
  "property_types": {
    "bedsitter": {
      "name": "bedsitter",
      "description": "A single room with a small kitchen area and private bathroom"
    },
    "apartment": {
      "name": "apartment",
      "description": "Multi-room unit in a building with shared facilities"
    }
  }
}
```

### Get Popular Areas

Retrieves popular residential areas by county.

**Endpoint**: `GET /popular-areas`

**Query Parameters**:
- `county` (string): Filter by county name

**Response** (200 OK):
```json
{
  "areas": {
    "Nairobi": [
      "Westlands",
      "Karen",
      "Kilimani",
      "Lavington"
    ],
    "Mombasa": [
      "Nyali",
      "Bamburi",
      "Shanzu"
    ]
  }
}
```

### Validate Phone Number

Validates a Kenyan phone number format.

**Endpoint**: `POST /validate-phone`

**Request Body**:
```json
{
  "phone_number": "0712345678"
}
```

**Response** (200 OK):
```json
{
  "phone_number": "0712345678",
  "is_valid": true
}
```

## Payment System

### Initiate Rent Payment (Tenant Only)

Initiates a rent payment via M-Pesa STK Push.

**Endpoint**: `POST /payments/initiate`

**Headers**: `Authorization: Bearer <tenant-token>`

**Request Body**:
```json
{
  "lease_id": "uuid-here",
  "amount": 45000,
  "phone_number": "254712345678",
  "payment_type": "rent"
}
```

**Response** (200 OK):
```json
{
  "message": "Payment initiated successfully",
  "payment_id": "uuid-here",
  "checkout_request_id": "ws_CO_123456789",
  "customer_message": "Please enter your M-Pesa PIN to complete the payment"
}
```

### Query Payment Status

Queries the status of an M-Pesa payment.

**Endpoint**: `GET /payments/status/{checkout_request_id}`

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "status": {
    "ResponseCode": "0",
    "ResponseDescription": "Success",
    "MerchantRequestID": "123456",
    "CheckoutRequestID": "ws_CO_123456789",
    "ResultCode": "0",
    "ResultDesc": "The service request is processed successfully."
  }
}
```

### Get Payments by Lease

Retrieves payment history for a specific lease.

**Endpoint**: `GET /payments/lease/{lease_id}`

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "payments": [
    {
      "id": "uuid-here",
      "lease_id": "uuid-here",
      "amount": 45000,
      "payment_type": "rent",
      "payment_method": "mpesa",
      "status": "completed",
      "transaction_id": "MPesa123456",
      "created_at": "2024-01-01T00:00:00Z",
      "completed_at": "2024-01-01T00:05:00Z"
    }
  ],
  "lease_id": "uuid-here"
}
```

### M-Pesa Callback (Internal)

Handles M-Pesa payment callbacks from Safaricom.

**Endpoint**: `POST /payments/mpesa/callback`

**Note**: This endpoint is called by Safaricom's servers and should not be called directly.

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request data",
  "details": "Validation error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing authentication token"
}
```

### 403 Forbidden
```json
{
  "error": "Forbidden",
  "message": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "message": "An unexpected error occurred"
}
```

## Rate Limiting

The API implements rate limiting to prevent abuse:
- **Authentication endpoints**: 5 requests per minute per IP
- **General endpoints**: 100 requests per minute per authenticated user
- **Payment endpoints**: 10 requests per minute per user

## Data Types

### Property Types
- `bedsitter`: Single room with kitchenette
- `studio`: Open plan living space
- `apartment`: Multi-room unit
- `maisonette`: Two-story unit
- `bungalow`: Single-story house
- `villa`: Luxury house
- `townhouse`: Multi-story shared-wall house
- `penthouse`: Top-floor luxury unit
- `duplex`: Two-unit building
- `commercial`: Business property

### Payment Types
- `rent`: Monthly rent payment
- `deposit`: Security deposit
- `utility`: Utility bills
- `maintenance`: Maintenance fees

### Payment Methods
- `mpesa`: M-Pesa mobile money
- `bank_transfer`: Bank transfer
- `cash`: Cash payment
- `cheque`: Bank cheque

### Payment Status
- `pending`: Payment initiated but not completed
- `completed`: Payment successfully processed
- `failed`: Payment failed
- `cancelled`: Payment cancelled by user
- `refunded`: Payment refunded

## Webhooks

### M-Pesa Payment Confirmation

When a payment is completed, the system processes the M-Pesa callback and updates the payment status. You can implement webhook listeners to get notified of payment status changes.

## SDK and Libraries

### Go Client Example
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type Client struct {
    BaseURL string
    Token   string
}

func (c *Client) GetProperties() (*PropertiesResponse, error) {
    req, _ := http.NewRequest("GET", c.BaseURL+"/properties", nil)
    req.Header.Set("Authorization", "Bearer "+c.Token)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result PropertiesResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}
```

## Testing

### Health Check

**Endpoint**: `GET /health`

**Response** (200 OK):
```json
{
  "status": "healthy",
  "service": "kenyan-real-estate-backend"
}
```

Use this endpoint to verify that the API is running and accessible.

