# Database Schema Design - Kenyan Real Estate Platform

## Overview
This document outlines the database schema for a real estate rental platform tailored to the Kenyan market.

## Tables

### 1. Users
Stores user information for both landlords and tenants.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    user_type ENUM('landlord', 'tenant', 'agent') NOT NULL,
    id_number VARCHAR(20) UNIQUE, -- Kenyan ID number
    profile_image_url TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 2. Counties
Kenyan counties for location-based filtering.

```sql
CREATE TABLE counties (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(10) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3. Sub_Counties
Sub-counties within each county.

```sql
CREATE TABLE sub_counties (
    id INT PRIMARY KEY AUTO_INCREMENT,
    county_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (county_id) REFERENCES counties(id) ON DELETE CASCADE
);
```

### 4. Properties
Main properties table with Kenyan-specific features.

```sql
CREATE TABLE properties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    landlord_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    property_type ENUM('apartment', 'house', 'bedsitter', 'studio', 'maisonette', 'bungalow', 'villa', 'commercial') NOT NULL,
    bedrooms INT DEFAULT 0,
    bathrooms INT DEFAULT 0,
    square_meters DECIMAL(10,2),
    rent_amount DECIMAL(12,2) NOT NULL, -- in KES
    deposit_amount DECIMAL(12,2), -- Security deposit
    county_id INT NOT NULL,
    sub_county_id INT,
    location_details TEXT, -- Specific location within sub-county
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    amenities JSON, -- Flexible storage for amenities
    utilities_included JSON, -- Water, electricity, internet, etc.
    parking_spaces INT DEFAULT 0,
    is_furnished BOOLEAN DEFAULT FALSE,
    is_available BOOLEAN DEFAULT TRUE,
    availability_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (landlord_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (county_id) REFERENCES counties(id),
    FOREIGN KEY (sub_county_id) REFERENCES sub_counties(id)
);
```

### 5. Property_Images
Images associated with properties.

```sql
CREATE TABLE property_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL,
    image_url TEXT NOT NULL,
    caption VARCHAR(255),
    is_primary BOOLEAN DEFAULT FALSE,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
);
```

### 6. Rental_Applications
Tenant applications for properties.

```sql
CREATE TABLE rental_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    application_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('pending', 'approved', 'rejected', 'withdrawn') DEFAULT 'pending',
    move_in_date DATE,
    message TEXT,
    monthly_income DECIMAL(12,2),
    employment_status VARCHAR(100),
    references JSON, -- Previous landlord, employer references
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_application (property_id, tenant_id)
);
```

### 7. Leases
Active lease agreements.

```sql
CREATE TABLE leases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    landlord_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    monthly_rent DECIMAL(12,2) NOT NULL,
    deposit_paid DECIMAL(12,2) NOT NULL,
    status ENUM('active', 'expired', 'terminated') DEFAULT 'active',
    lease_terms TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (landlord_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 8. Payments
Payment records for rent and deposits.

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lease_id UUID NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    payment_type ENUM('rent', 'deposit', 'utility', 'maintenance') NOT NULL,
    payment_method ENUM('mpesa', 'bank_transfer', 'cash', 'cheque') NOT NULL,
    mpesa_transaction_id VARCHAR(100), -- For M-Pesa payments
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date DATE,
    status ENUM('pending', 'completed', 'failed', 'refunded') DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (lease_id) REFERENCES leases(id) ON DELETE CASCADE
);
```

### 9. Property_Views
Track property views for analytics.

```sql
CREATE TABLE property_views (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id UUID NOT NULL,
    user_id UUID, -- NULL for anonymous views
    ip_address VARCHAR(45),
    user_agent TEXT,
    viewed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);
```

### 10. Favorites
User favorite properties.

```sql
CREATE TABLE favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    property_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    UNIQUE KEY unique_favorite (user_id, property_id)
);
```

## Indexes

```sql
-- Performance indexes
CREATE INDEX idx_properties_county ON properties(county_id);
CREATE INDEX idx_properties_type ON properties(property_type);
CREATE INDEX idx_properties_rent ON properties(rent_amount);
CREATE INDEX idx_properties_available ON properties(is_available);
CREATE INDEX idx_properties_location ON properties(latitude, longitude);
CREATE INDEX idx_payments_lease ON payments(lease_id);
CREATE INDEX idx_payments_date ON payments(payment_date);
CREATE INDEX idx_applications_property ON rental_applications(property_id);
CREATE INDEX idx_applications_tenant ON rental_applications(tenant_id);
```

## Kenyan-Specific Features

1. **Counties and Sub-Counties**: Based on Kenya's 47 counties
2. **Property Types**: Includes local types like bedsitter, maisonette
3. **M-Pesa Integration**: Support for Kenya's mobile money system
4. **KES Currency**: All amounts in Kenyan Shillings
5. **ID Numbers**: Support for Kenyan national ID numbers
6. **Local Amenities**: Flexible JSON storage for Kenyan-specific amenities

