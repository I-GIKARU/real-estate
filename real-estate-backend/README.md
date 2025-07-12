# Real Estate Backend API

A comprehensive Go backend for a real estate rental platform featuring JWT authentication, property management, user roles, and password reset functionality.

## Features

- **User Authentication**: JWT-based authentication with role-based access control (admins, agents, and tenants)
- **Property Management**: Full CRUD operations for properties with search and filtering
- **Password Reset**: Secure password reset functionality with email verification
- **User Management**: 
  - Admin approval system for agents
  - Email verification for new users
  - Role-based permissions
  - User profile management
- **Email Services**: Automated email notifications for verification and password reset
- **Web Forms**: HTML forms for password reset (embedded in backend)
- **Image Management**: Property image upload and management
- **Location Services**: County and sub-county data

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL
- **Authentication**: JWT tokens
- **Payment**: M-Pesa Daraja API
- **Dependencies**:
  - `github.com/gin-gonic/gin` - HTTP web framework
  - `github.com/lib/pq` - PostgreSQL driver
  - `github.com/golang-jwt/jwt/v5` - JWT implementation
  - `github.com/google/uuid` - UUID generation
  - `golang.org/x/crypto/bcrypt` - Password hashing
  - `github.com/joho/godotenv` - Environment variable loading

## Project Structure

```
kenyan-real-estate-backend/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── database/
│   │   └── connection.go          # Database connection utilities
│   ├── handlers/
│   │   ├── user.go                # User authentication handlers
│   │   ├── property.go            # Property management handlers
│   │   ├── location.go            # Location data handlers
│   │   ├── payment.go             # Payment processing handlers
│   │   └── kenyan_features.go     # Kenyan-specific features
│   ├── middleware/
│   │   └── auth.go                # Authentication middleware
│   ├── models/
│   │   ├── user.go                # User model and repository
│   │   ├── property.go            # Property model and repository
│   │   ├── location.go            # County/SubCounty models
│   │   └── rental.go              # Rental and payment models
│   └── services/
│       └── mpesa.go               # M-Pesa integration service
├── pkg/
│   ├── auth/
│   │   └── jwt.go                 # JWT utilities
│   └── utils/
│       └── kenyan_features.go     # Kenyan-specific utilities
├── migrations/
│   ├── 001_initial_schema.sql     # Database schema
│   └── 002_kenyan_counties_data.sql # Kenyan location data
├── docs/
│   └── database_schema.md         # Database documentation
├── .env.example                   # Environment variables template
├── go.mod                         # Go module definition
└── README.md                      # This file
```

## Installation and Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- M-Pesa Developer Account (for payment integration)

### 1. Clone the Repository

```bash
git clone <repository-url>
cd kenyan-real-estate-backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

Create a PostgreSQL database and run the migrations:

```bash
# Create database
createdb kenyan_real_estate

# Run migrations
psql -d kenyan_real_estate -f migrations/001_initial_schema.sql
psql -d kenyan_real_estate -f migrations/002_kenyan_counties_data.sql
```

### 4. Environment Configuration

Copy the example environment file and configure your settings:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
APP_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kenyan_real_estate
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24

# M-Pesa Configuration (Get from Safaricom Developer Portal)
MPESA_CONSUMER_KEY=your_consumer_key
MPESA_CONSUMER_SECRET=your_consumer_secret
MPESA_ENVIRONMENT=sandbox
MPESA_PASS_KEY=your_pass_key
MPESA_SHORT_CODE=your_short_code

# Upload Configuration
MAX_FILE_SIZE=10485760
UPLOAD_DIR=./uploads
```

### 5. Run the Application

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### User Roles

- **landlord**: Can create, update, and manage properties
- **tenant**: Can view properties and make rental payments

## API Endpoints

### Authentication Endpoints

#### Register User
```http
POST /api/v1/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "0712345678",
  "user_type": "tenant"
}
```

#### Login
```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Refresh Token
```http
POST /api/v1/refresh-token
Authorization: Bearer <token>
```

### Property Endpoints

#### Get Public Properties (with filtering)
```http
GET /api/v1/properties?county_id=1&property_type=apartment&min_rent=10000&max_rent=50000
```

#### Get Single Property
```http
GET /api/v1/properties/{id}
```

#### Create Property (Landlord only)
```http
POST /api/v1/properties
Authorization: Bearer <landlord-token>
Content-Type: application/json

{
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
  "location_details": "Near Yaya Centre",
  "amenities": ["24/7 Security", "Swimming Pool", "Gym"],
  "utilities_included": ["water", "security"],
  "parking_spaces": 1,
  "is_furnished": true,
  "availability_date": "2024-01-01T00:00:00Z"
}
```

### Location Endpoints

#### Get All Counties
```http
GET /api/v1/counties
```

#### Get Sub-counties by County
```http
GET /api/v1/counties/{county_id}/sub-counties
```

### Kenyan Features Endpoints

#### Get Property Amenities
```http
GET /api/v1/amenities
GET /api/v1/amenities?category=security
```

#### Get Property Types
```http
GET /api/v1/property-types
```

#### Get Popular Areas
```http
GET /api/v1/popular-areas
GET /api/v1/popular-areas?county=Nairobi
```

#### Validate Phone Number
```http
POST /api/v1/validate-phone
Content-Type: application/json

{
  "phone_number": "0712345678"
}
```

### Payment Endpoints

#### Initiate Rent Payment (Tenant only)
```http
POST /api/v1/payments/initiate
Authorization: Bearer <tenant-token>
Content-Type: application/json

{
  "lease_id": "uuid-here",
  "amount": 45000,
  "phone_number": "254712345678",
  "payment_type": "rent"
}
```

#### Query Payment Status
```http
GET /api/v1/payments/status/{checkout_request_id}
Authorization: Bearer <token>
```

#### Get Payments by Lease
```http
GET /api/v1/payments/lease/{lease_id}
Authorization: Bearer <token>
```

## M-Pesa Integration

The application integrates with Safaricom's M-Pesa Daraja API for payment processing:

### Features:
- **STK Push**: Initiate payments from customer's phone
- **Payment Validation**: Validate Kenyan phone numbers
- **Transaction Tracking**: Track payment status and history
- **Callback Handling**: Process M-Pesa payment confirmations

### Setup:
1. Register at [Safaricom Developer Portal](https://developer.safaricom.co.ke/)
2. Create an app and get Consumer Key and Secret
3. Configure your shortcode and passkey
4. Set up callback URLs for payment confirmations

## Database Schema

The application uses PostgreSQL with the following main tables:

- **users**: User accounts (landlords and tenants)
- **counties**: Kenyan counties (47 counties)
- **sub_counties**: Sub-counties within each county
- **properties**: Property listings with full details
- **property_images**: Property photos and media
- **leases**: Rental agreements between landlords and tenants
- **payments**: Payment records and M-Pesa transactions

## Security Features

- **Password Hashing**: bcrypt for secure password storage
- **JWT Authentication**: Stateless authentication with role-based access
- **Input Validation**: Comprehensive request validation
- **CORS Support**: Configurable cross-origin resource sharing
- **Environment Variables**: Sensitive data stored in environment variables

## Kenyan Market Features

### Counties and Locations
- Complete list of all 47 Kenyan counties
- Sub-counties for detailed location targeting
- Popular residential areas by county

### Property Types
- Bedsitter, Studio, Apartment, Maisonette
- Bungalow, Villa, Townhouse, Penthouse
- Commercial properties

### Local Amenities
- Security features (24/7 security, CCTV, electric fence)
- Utilities (borehole water, solar power, backup generator)
- Facilities (swimming pool, gym, children playground)

### Payment Integration
- M-Pesa STK Push for rent payments
- Kenyan phone number validation
- KES currency formatting

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o kenyan-real-estate-backend cmd/server/main.go
```

### Docker Support
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Support

For support and questions:
- Email: support@example.com
- Documentation: [API Docs](./docs/)
- Issues: [GitHub Issues](https://github.com/your-repo/issues)

