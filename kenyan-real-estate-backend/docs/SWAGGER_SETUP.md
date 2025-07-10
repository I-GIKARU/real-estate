# Swagger API Documentation

## Overview

This project now includes comprehensive Swagger/OpenAPI documentation for all API endpoints. This makes it easy to explore, test, and integrate with the Kenyan Real Estate API.

## Accessing Swagger UI

Once your server is running, you can access the interactive Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## Features

### 🎯 **Interactive API Explorer**
- Test all endpoints directly from the browser
- View request/response schemas
- See example payloads
- Authentication testing with JWT tokens

### 📝 **Comprehensive Documentation**
- All endpoints documented with descriptions
- Request/response examples
- Parameter validation rules
- Error response documentation

### 🔐 **Authentication Support**
- JWT Bearer token authentication
- Easy token testing in the UI
- Protected endpoint indicators

## API Overview

### Authentication Endpoints
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login
- `POST /api/v1/refresh-token` - Token refresh

### User Management
- `GET /api/v1/profile` - Get user profile (requires auth)
- `PUT /api/v1/profile` - Update user profile (requires auth)

### Property Management
- `GET /api/v1/properties` - Get public property listings (with filtering)
- `GET /api/v1/properties/{id}` - Get specific property details
- `POST /api/v1/properties` - Create property (landlord only)
- `PUT /api/v1/properties/{id}` - Update property (landlord only)
- `DELETE /api/v1/properties/{id}` - Delete property (landlord only)

### Location Data
- `GET /api/v1/counties` - Get all Kenya counties
- `GET /api/v1/counties/{id}` - Get specific county
- `GET /api/v1/counties/{id}/sub-counties` - Get sub-counties
- `GET /api/v1/sub-counties/{id}` - Get specific sub-county

### Kenyan Features
- `GET /api/v1/amenities` - Get available amenities
- `GET /api/v1/property-types` - Get property types
- `GET /api/v1/utilities` - Get utility options
- `POST /api/v1/validate-phone` - Validate Kenyan phone numbers
- `POST /api/v1/format-currency` - Format KES currency

## How to Use Swagger UI

### 1. **Exploring Endpoints**
1. Navigate to the Swagger UI
2. Browse endpoints by category (tags)
3. Click on any endpoint to expand details
4. View required parameters and response schemas

### 2. **Testing Endpoints**

#### **Public Endpoints (No Auth Required)**
1. Click "Try it out" on any public endpoint
2. Fill in required parameters
3. Click "Execute"
4. View the response

#### **Protected Endpoints (Auth Required)**
1. First, register or login to get a JWT token:
   ```bash
   # Example registration
   curl -X POST "http://localhost:8080/api/v1/register" \
   -H "Content-Type: application/json" \
   -d '{
     "email": "user@example.com",
     "password": "password123",
     "first_name": "John",
     "last_name": "Doe",
     "phone_number": "+254700000000",
     "user_type": "tenant"
   }'
   ```

2. Copy the token from the response
3. Click the "Authorize" button at the top of Swagger UI
4. Enter: `Bearer YOUR_JWT_TOKEN`
5. Click "Authorize"
6. Now you can test protected endpoints

### 3. **Property Search Examples**

#### **Basic Property Search**
```
GET /api/v1/properties
```

#### **Filtered Search**
```
GET /api/v1/properties?county_id=47&min_rent=20000&max_rent=100000&bedrooms=2&property_type=apartment
```

#### **Search in Nairobi (County ID: 47)**
```
GET /api/v1/properties?county_id=47&sub_county_id=1&is_furnished=true
```

## Model Schemas

### User Models
- **CreateUserRequest**: Registration data
- **LoginRequest**: Login credentials  
- **UserResponse**: User profile data (no sensitive info)

### Property Models
- **CreatePropertyRequest**: Property creation data
- **Property**: Complete property information
- **PropertySearchFilters**: Search and filter options

### Location Models
- **County**: Kenya county information
- **SubCounty**: Sub-county details
- **PropertyImage**: Property image data

## Property Types (Kenyan Context)

The API supports these property types specific to the Kenyan market:
- `apartment` - Multi-unit residential building
- `house` - Standalone residential house
- `bedsitter` - Single room with kitchenette
- `studio` - Open plan living space
- `maisonette` - Multi-level apartment/townhouse
- `bungalow` - Single-story house
- `villa` - Luxury residential property
- `commercial` - Commercial properties

## Response Examples

### Success Response (Property Listing)
```json
{
  "properties": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Beautiful 2BR Apartment in Westlands",
      "property_type": "apartment",
      "bedrooms": 2,
      "bathrooms": 2,
      "rent_amount": 75000,
      "county": {
        "id": 47,
        "name": "Nairobi",
        "code": "047"
      },
      "amenities": {
        "wifi": true,
        "parking": true,
        "gym": false,
        "swimming_pool": true
      }
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

### Error Response
```json
{
  "error": "Invalid request data",
  "details": "rent_amount must be greater than 0"
}
```

## Development Workflow

### Adding New Endpoints

1. **Add Swagger Annotations** to your handler:
```go
// CreateProperty handles property creation
// @Summary Create a new property
// @Description Create a new property listing (landlord only)
// @Tags Properties
// @Accept json
// @Produce json
// @Security Bearer
// @Param property body models.CreatePropertyRequest true "Property data"
// @Success 201 {object} object{message=string,property=models.Property}
// @Failure 400 {object} object{error=string}
// @Router /properties [post]
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
    // Handler implementation
}
```

2. **Regenerate Documentation**:
```bash
~/go/bin/swag init -g cmd/server/main.go -o ./docs
```

3. **Rebuild and Test**:
```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

### Swagger Annotation Tags

- `@Summary` - Brief endpoint description
- `@Description` - Detailed endpoint description  
- `@Tags` - Group endpoints by category
- `@Accept` - Request content type
- `@Produce` - Response content type
- `@Security` - Authentication requirement
- `@Param` - Request parameters
- `@Success` - Success response format
- `@Failure` - Error response format
- `@Router` - Endpoint path and method

## Best Practices

### 1. **Consistent Error Responses**
Always use the same error format:
```go
c.JSON(http.StatusBadRequest, gin.H{
    "error": "Short error message",
    "details": "Detailed error information"
})
```

### 2. **Descriptive Documentation**
- Use clear, descriptive summaries
- Include examples in descriptions
- Document all possible error cases

### 3. **Group Related Endpoints**
Use tags to organize endpoints:
- `Authentication` - Login, register, tokens
- `Users` - User management
- `Properties` - Property CRUD operations
- `Locations` - County/sub-county data

### 4. **Security Documentation**
- Mark protected endpoints with `@Security Bearer`
- Document authentication requirements
- Show example authorization headers

## Troubleshooting

### Swagger UI Not Loading
1. Check server is running on correct port
2. Verify `/swagger/index.html` endpoint
3. Check browser console for errors

### Documentation Not Updating
1. Regenerate docs: `~/go/bin/swag init -g cmd/server/main.go -o ./docs`
2. Rebuild application
3. Restart server

### Authentication Issues
1. Ensure JWT token is valid
2. Check token format: `Bearer YOUR_TOKEN`
3. Verify token hasn't expired

## Additional Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Guide](https://swagger.io/tools/swagger-ui/)

## Screenshots

When you access `http://localhost:8080/swagger/index.html`, you'll see:

1. **API Overview** - Title, description, and version
2. **Endpoint Categories** - Organized by tags
3. **Interactive Testing** - Try endpoints directly
4. **Authentication** - JWT token support
5. **Model Schemas** - Request/response structures
6. **Response Examples** - Sample API responses

This makes your Kenyan Real Estate API much easier to work with for frontend developers, mobile app developers, and API consumers!
