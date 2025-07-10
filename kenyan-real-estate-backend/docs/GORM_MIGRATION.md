# GORM Migration Guide

## Overview

This document describes the migration from raw SQL to GORM ORM in the Kenyan Real Estate Backend application.

## What Changed

### 1. Database Connection
- **Before**: Used `database/sql` with `lib/pq` PostgreSQL driver
- **After**: Uses GORM with `gorm.io/driver/postgres`

### 2. Models
All models have been updated with GORM struct tags:

#### User Model
```go
type User struct {
    ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Email           string     `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash    string     `json:"-" gorm:"not null"`
    // ... other fields with GORM tags
    DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relationships
    Properties []Property `json:"properties,omitempty" gorm:"foreignKey:LandlordID"`
}
```

#### Property Model
```go
type Property struct {
    ID                uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    LandlordID        uuid.UUID         `json:"landlord_id" gorm:"type:uuid;not null"`
    // ... other fields
    Amenities         Amenities         `json:"amenities" gorm:"type:jsonb"`
    UtilitiesIncluded UtilitiesIncluded `json:"utilities_included" gorm:"type:jsonb"`
    DeletedAt         gorm.DeletedAt    `json:"-" gorm:"index"`
    
    // Relationships
    County    *County         `json:"county,omitempty" gorm:"foreignKey:CountyID"`
    SubCounty *SubCounty      `json:"sub_county,omitempty" gorm:"foreignKey:SubCountyID"`
    Landlord  *User           `json:"landlord,omitempty" gorm:"foreignKey:LandlordID"`
    Images    []*PropertyImage `json:"images,omitempty" gorm:"foreignKey:PropertyID"`
}
```

### 3. Repositories
Repositories now use GORM methods instead of raw SQL:

#### Before (Raw SQL)
```go
func (r *UserRepository) Create(user *User) error {
    query := `INSERT INTO users (...) VALUES (...)`
    // Complex SQL execution
}
```

#### After (GORM)
```go
func (r *UserRepository) Create(user *User) error {
    return r.db.Create(user).Error
}
```

### 4. Auto-Migration
GORM now handles database schema creation/updates automatically:

```go
// Run auto-migration
if err := database.AutoMigrate(
    &models.User{},
    &models.County{},
    &models.SubCounty{},
    &models.Property{},
    &models.PropertyImage{},
); err != nil {
    log.Fatal("Failed to run database migrations:", err)
}
```

## Key Benefits

### 1. **Simplified Database Operations**
- No more complex SQL queries
- Built-in validation and type safety
- Automatic relationship loading with `Preload()`

### 2. **Soft Deletes**
- Built-in soft delete support with `gorm.DeletedAt`
- Automatically excludes deleted records from queries

### 3. **Automatic Migrations**
- Schema changes are handled automatically
- No manual SQL migration scripts needed

### 4. **Better Relationships**
- Easy relationship definitions with struct tags
- Automatic foreign key constraints
- Eager/lazy loading support

### 5. **Hooks and Callbacks**
- `BeforeCreate`, `AfterCreate`, etc.
- Automatic ID generation and timestamps

## Database Schema Features

### 1. **UUID Primary Keys**
```go
ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
```

### 2. **JSON Fields (PostgreSQL)**
```go
Amenities Amenities `gorm:"type:jsonb"`
```

### 3. **Indexes and Constraints**
```go
Email string `gorm:"uniqueIndex;not null"`
```

### 4. **Automatic Timestamps**
```go
CreatedAt time.Time `gorm:"autoCreateTime"`
UpdatedAt time.Time `gorm:"autoUpdateTime"`
```

## Setup Instructions

### 1. **Environment Variables**
Ensure your `.env` file includes:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kenyan_real_estate
DB_SSL_MODE=disable
APP_ENV=development
```

### 2. **Database Setup**
1. Create PostgreSQL database
2. Run the application - GORM will auto-migrate tables
3. Optionally run the seed script:
   ```bash
   psql -h localhost -U postgres -d kenyan_real_estate -f scripts/seed_data.sql
   ```

### 3. **Running the Application**
```bash
go mod tidy
go build -o bin/server cmd/server/main.go
./bin/server
```

## GORM Query Examples

### 1. **Basic Queries**
```go
// Find by ID
var user User
db.First(&user, "id = ?", userID)

// Find with conditions
var properties []Property
db.Where("rent_amount >= ? AND bedrooms >= ?", minRent, minBedrooms).Find(&properties)
```

### 2. **Relationships**
```go
// Preload relationships
var property Property
db.Preload("County").Preload("SubCounty").Preload("Landlord").Preload("Images").First(&property, id)

// Create with associations
property := Property{
    Title: "Beautiful Apartment",
    LandlordID: landlordID,
    CountyID: countyID,
}
db.Create(&property)
```

### 3. **Advanced Queries**
```go
// Search with multiple conditions
query := db.Model(&Property{})
if filters.MinRent != nil {
    query = query.Where("rent_amount >= ?", *filters.MinRent)
}
if filters.PropertyType != nil {
    query = query.Where("property_type = ?", *filters.PropertyType)
}
query.Find(&properties)
```

### 4. **Transactions**
```go
db.Transaction(func(tx *gorm.DB) error {
    // Multiple operations in transaction
    if err := tx.Model(&PropertyImage{}).Where("property_id = ?", propertyID).Update("is_primary", false).Error; err != nil {
        return err
    }
    return tx.Model(&PropertyImage{}).Where("id = ?", imageID).Update("is_primary", true).Error
})
```

## Migration Notes

### 1. **Backward Compatibility**
- Database schema remains the same
- Existing data is preserved
- Table names unchanged (`users`, `properties`, etc.)

### 2. **Custom Types**
The `Amenities` and `UtilitiesIncluded` types still implement `driver.Valuer` and `sql.Scanner` for JSON handling.

### 3. **Foreign Key Relationships**
GORM automatically creates foreign key constraints based on struct tags.

## Performance Considerations

### 1. **Connection Pooling**
```go
sqlDB, err := DB.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(25)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### 2. **Selective Loading**
Use `Select()` to load only needed fields:
```go
db.Select("id", "title", "rent_amount").Find(&properties)
```

### 3. **Preloading**
Only preload relationships when needed to avoid N+1 queries.

## Troubleshooting

### 1. **Migration Errors**
If auto-migration fails, check:
- Database permissions
- PostgreSQL version compatibility
- Environment variables

### 2. **Query Issues**
Enable GORM logging for debugging:
```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

### 3. **Relationship Problems**
Ensure foreign key fields match:
- `LandlordID` in Property → `ID` in User
- `CountyID` in Property → `ID` in County

## Next Steps

1. **Add Remaining Models**: Lease, Payment, RentalApplication models
2. **Implement Caching**: Redis for frequently accessed data
3. **Add Validations**: Custom GORM validators
4. **Optimize Queries**: Add database indexes for common queries
5. **Add Tests**: Unit tests for repository methods

## Resources

- [GORM Documentation](https://gorm.io/docs/)
- [PostgreSQL GORM Driver](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
- [GORM Associations](https://gorm.io/docs/associations.html)
