package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PropertyType represents different types of properties in Kenya
type PropertyType string

const (
	PropertyTypeApartment  PropertyType = "apartment"
	PropertyTypeHouse      PropertyType = "house"
	PropertyTypeBedsitter  PropertyType = "bedsitter"
	PropertyTypeStudio     PropertyType = "studio"
	PropertyTypeMaisonette PropertyType = "maisonette"
	PropertyTypeBungalow   PropertyType = "bungalow"
	PropertyTypeVilla      PropertyType = "villa"
	PropertyTypeCommercial PropertyType = "commercial"
)

// Amenities represents property amenities as JSON
type Amenities map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (a Amenities) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan implements the sql.Scanner interface for database retrieval
func (a *Amenities) Scan(value interface{}) error {
	if value == nil {
		*a = make(Amenities)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	default:
		return fmt.Errorf("cannot scan %T into Amenities", value)
	}
}

// UtilitiesIncluded represents utilities included in rent
type UtilitiesIncluded map[string]interface{}

// Value implements the driver.Valuer interface for database storage
func (u UtilitiesIncluded) Value() (driver.Value, error) {
	return json.Marshal(u)
}

// Scan implements the sql.Scanner interface for database retrieval
func (u *UtilitiesIncluded) Scan(value interface{}) error {
	if value == nil {
		*u = make(UtilitiesIncluded)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, u)
	case string:
		return json.Unmarshal([]byte(v), u)
	default:
		return fmt.Errorf("cannot scan %T into UtilitiesIncluded", value)
	}
}

// Property represents a rental property
type Property struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	LandlordID        uuid.UUID         `json:"landlord_id" db:"landlord_id"`
	Title             string            `json:"title" db:"title"`
	Description       *string           `json:"description,omitempty" db:"description"`
	PropertyType      PropertyType      `json:"property_type" db:"property_type"`
	Bedrooms          int               `json:"bedrooms" db:"bedrooms"`
	Bathrooms         int               `json:"bathrooms" db:"bathrooms"`
	SquareMeters      *float64          `json:"square_meters,omitempty" db:"square_meters"`
	RentAmount        float64           `json:"rent_amount" db:"rent_amount"`
	DepositAmount     *float64          `json:"deposit_amount,omitempty" db:"deposit_amount"`
	CountyID          int               `json:"county_id" db:"county_id"`
	SubCountyID       *int              `json:"sub_county_id,omitempty" db:"sub_county_id"`
	LocationDetails   *string           `json:"location_details,omitempty" db:"location_details"`
	Latitude          *float64          `json:"latitude,omitempty" db:"latitude"`
	Longitude         *float64          `json:"longitude,omitempty" db:"longitude"`
	Amenities         Amenities         `json:"amenities" db:"amenities"`
	UtilitiesIncluded UtilitiesIncluded `json:"utilities_included" db:"utilities_included"`
	ParkingSpaces     int               `json:"parking_spaces" db:"parking_spaces"`
	IsFurnished       bool              `json:"is_furnished" db:"is_furnished"`
	IsAvailable       bool              `json:"is_available" db:"is_available"`
	AvailabilityDate  *time.Time        `json:"availability_date,omitempty" db:"availability_date"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
	
	// Joined fields
	County    *County    `json:"county,omitempty"`
	SubCounty *SubCounty `json:"sub_county,omitempty"`
	Landlord  *User      `json:"landlord,omitempty"`
	Images    []*PropertyImage `json:"images,omitempty"`
}

// CreatePropertyRequest represents the request to create a new property
type CreatePropertyRequest struct {
	Title             string            `json:"title" binding:"required"`
	Description       *string           `json:"description,omitempty"`
	PropertyType      PropertyType      `json:"property_type" binding:"required"`
	Bedrooms          int               `json:"bedrooms" binding:"min=0"`
	Bathrooms         int               `json:"bathrooms" binding:"min=0"`
	SquareMeters      *float64          `json:"square_meters,omitempty" binding:"omitempty,min=0"`
	RentAmount        float64           `json:"rent_amount" binding:"required,min=0"`
	DepositAmount     *float64          `json:"deposit_amount,omitempty" binding:"omitempty,min=0"`
	CountyID          int               `json:"county_id" binding:"required"`
	SubCountyID       *int              `json:"sub_county_id,omitempty"`
	LocationDetails   *string           `json:"location_details,omitempty"`
	Latitude          *float64          `json:"latitude,omitempty"`
	Longitude         *float64          `json:"longitude,omitempty"`
	Amenities         Amenities         `json:"amenities"`
	UtilitiesIncluded UtilitiesIncluded `json:"utilities_included"`
	ParkingSpaces     int               `json:"parking_spaces" binding:"min=0"`
	IsFurnished       bool              `json:"is_furnished"`
	AvailabilityDate  *time.Time        `json:"availability_date,omitempty"`
}

// UpdatePropertyRequest represents the request to update a property
type UpdatePropertyRequest struct {
	Title             *string           `json:"title,omitempty"`
	Description       *string           `json:"description,omitempty"`
	Bedrooms          *int              `json:"bedrooms,omitempty" binding:"omitempty,min=0"`
	Bathrooms         *int              `json:"bathrooms,omitempty" binding:"omitempty,min=0"`
	SquareMeters      *float64          `json:"square_meters,omitempty" binding:"omitempty,min=0"`
	RentAmount        *float64          `json:"rent_amount,omitempty" binding:"omitempty,min=0"`
	DepositAmount     *float64          `json:"deposit_amount,omitempty" binding:"omitempty,min=0"`
	LocationDetails   *string           `json:"location_details,omitempty"`
	Latitude          *float64          `json:"latitude,omitempty"`
	Longitude         *float64          `json:"longitude,omitempty"`
	Amenities         *Amenities        `json:"amenities,omitempty"`
	UtilitiesIncluded *UtilitiesIncluded `json:"utilities_included,omitempty"`
	ParkingSpaces     *int              `json:"parking_spaces,omitempty" binding:"omitempty,min=0"`
	IsFurnished       *bool             `json:"is_furnished,omitempty"`
	IsAvailable       *bool             `json:"is_available,omitempty"`
	AvailabilityDate  *time.Time        `json:"availability_date,omitempty"`
}

// PropertySearchFilters represents search filters for properties
type PropertySearchFilters struct {
	CountyID         *int          `json:"county_id,omitempty"`
	SubCountyID      *int          `json:"sub_county_id,omitempty"`
	PropertyType     *PropertyType `json:"property_type,omitempty"`
	MinRent          *float64      `json:"min_rent,omitempty"`
	MaxRent          *float64      `json:"max_rent,omitempty"`
	MinBedrooms      *int          `json:"min_bedrooms,omitempty"`
	MaxBedrooms      *int          `json:"max_bedrooms,omitempty"`
	MinBathrooms     *int          `json:"min_bathrooms,omitempty"`
	IsFurnished      *bool         `json:"is_furnished,omitempty"`
	HasParkingSpaces *bool         `json:"has_parking_spaces,omitempty"`
	IsAvailable      *bool         `json:"is_available,omitempty"`
	Limit            int           `json:"limit,omitempty"`
	Offset           int           `json:"offset,omitempty"`
}

// PropertyRepository handles database operations for properties
type PropertyRepository struct {
	db *sql.DB
}

// NewPropertyRepository creates a new property repository
func NewPropertyRepository(db *sql.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

// Create creates a new property
func (r *PropertyRepository) Create(property *Property) error {
	query := `
		INSERT INTO properties (
			id, landlord_id, title, description, property_type, bedrooms, bathrooms,
			square_meters, rent_amount, deposit_amount, county_id, sub_county_id,
			location_details, latitude, longitude, amenities, utilities_included,
			parking_spaces, is_furnished, is_available, availability_date, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
		) RETURNING id, created_at, updated_at`

	property.ID = uuid.New()
	property.IsAvailable = true
	now := time.Now()
	property.CreatedAt = now
	property.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		property.ID,
		property.LandlordID,
		property.Title,
		property.Description,
		property.PropertyType,
		property.Bedrooms,
		property.Bathrooms,
		property.SquareMeters,
		property.RentAmount,
		property.DepositAmount,
		property.CountyID,
		property.SubCountyID,
		property.LocationDetails,
		property.Latitude,
		property.Longitude,
		property.Amenities,
		property.UtilitiesIncluded,
		property.ParkingSpaces,
		property.IsFurnished,
		property.IsAvailable,
		property.AvailabilityDate,
		property.CreatedAt,
		property.UpdatedAt,
	).Scan(&property.ID, &property.CreatedAt, &property.UpdatedAt)

	return err
}

// GetByID retrieves a property by ID
func (r *PropertyRepository) GetByID(id uuid.UUID) (*Property, error) {
	property := &Property{}
	query := `
		SELECT p.id, p.landlord_id, p.title, p.description, p.property_type, p.bedrooms,
			   p.bathrooms, p.square_meters, p.rent_amount, p.deposit_amount, p.county_id,
			   p.sub_county_id, p.location_details, p.latitude, p.longitude, p.amenities,
			   p.utilities_included, p.parking_spaces, p.is_furnished, p.is_available,
			   p.availability_date, p.created_at, p.updated_at
		FROM properties p
		WHERE p.id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&property.ID,
		&property.LandlordID,
		&property.Title,
		&property.Description,
		&property.PropertyType,
		&property.Bedrooms,
		&property.Bathrooms,
		&property.SquareMeters,
		&property.RentAmount,
		&property.DepositAmount,
		&property.CountyID,
		&property.SubCountyID,
		&property.LocationDetails,
		&property.Latitude,
		&property.Longitude,
		&property.Amenities,
		&property.UtilitiesIncluded,
		&property.ParkingSpaces,
		&property.IsFurnished,
		&property.IsAvailable,
		&property.AvailabilityDate,
		&property.CreatedAt,
		&property.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return property, nil
}

// GetByLandlordID retrieves properties by landlord ID
func (r *PropertyRepository) GetByLandlordID(landlordID uuid.UUID, limit, offset int) ([]*Property, error) {
	query := `
		SELECT p.id, p.landlord_id, p.title, p.description, p.property_type, p.bedrooms,
			   p.bathrooms, p.square_meters, p.rent_amount, p.deposit_amount, p.county_id,
			   p.sub_county_id, p.location_details, p.latitude, p.longitude, p.amenities,
			   p.utilities_included, p.parking_spaces, p.is_furnished, p.is_available,
			   p.availability_date, p.created_at, p.updated_at
		FROM properties p
		WHERE p.landlord_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, landlordID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []*Property
	for rows.Next() {
		property := &Property{}
		err := rows.Scan(
			&property.ID,
			&property.LandlordID,
			&property.Title,
			&property.Description,
			&property.PropertyType,
			&property.Bedrooms,
			&property.Bathrooms,
			&property.SquareMeters,
			&property.RentAmount,
			&property.DepositAmount,
			&property.CountyID,
			&property.SubCountyID,
			&property.LocationDetails,
			&property.Latitude,
			&property.Longitude,
			&property.Amenities,
			&property.UtilitiesIncluded,
			&property.ParkingSpaces,
			&property.IsFurnished,
			&property.IsAvailable,
			&property.AvailabilityDate,
			&property.CreatedAt,
			&property.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}

	return properties, nil
}

// Update updates a property
func (r *PropertyRepository) Update(property *Property) error {
	query := `
		UPDATE properties
		SET title = $2, description = $3, bedrooms = $4, bathrooms = $5, square_meters = $6,
			rent_amount = $7, deposit_amount = $8, location_details = $9, latitude = $10,
			longitude = $11, amenities = $12, utilities_included = $13, parking_spaces = $14,
			is_furnished = $15, is_available = $16, availability_date = $17, updated_at = $18
		WHERE id = $1`

	property.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		property.ID,
		property.Title,
		property.Description,
		property.Bedrooms,
		property.Bathrooms,
		property.SquareMeters,
		property.RentAmount,
		property.DepositAmount,
		property.LocationDetails,
		property.Latitude,
		property.Longitude,
		property.Amenities,
		property.UtilitiesIncluded,
		property.ParkingSpaces,
		property.IsFurnished,
		property.IsAvailable,
		property.AvailabilityDate,
		property.UpdatedAt,
	)

	return err
}

// Delete deletes a property
func (r *PropertyRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM properties WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Search searches for properties based on filters
func (r *PropertyRepository) Search(filters *PropertySearchFilters) ([]*Property, error) {
	baseQuery := `
		SELECT p.id, p.landlord_id, p.title, p.description, p.property_type, p.bedrooms,
			   p.bathrooms, p.square_meters, p.rent_amount, p.deposit_amount, p.county_id,
			   p.sub_county_id, p.location_details, p.latitude, p.longitude, p.amenities,
			   p.utilities_included, p.parking_spaces, p.is_furnished, p.is_available,
			   p.availability_date, p.created_at, p.updated_at
		FROM properties p
		WHERE 1=1`

	var args []interface{}
	argIndex := 1

	// Build dynamic WHERE clause based on filters
	if filters.CountyID != nil {
		baseQuery += fmt.Sprintf(" AND p.county_id = $%d", argIndex)
		args = append(args, *filters.CountyID)
		argIndex++
	}

	if filters.SubCountyID != nil {
		baseQuery += fmt.Sprintf(" AND p.sub_county_id = $%d", argIndex)
		args = append(args, *filters.SubCountyID)
		argIndex++
	}

	if filters.PropertyType != nil {
		baseQuery += fmt.Sprintf(" AND p.property_type = $%d", argIndex)
		args = append(args, *filters.PropertyType)
		argIndex++
	}

	if filters.MinRent != nil {
		baseQuery += fmt.Sprintf(" AND p.rent_amount >= $%d", argIndex)
		args = append(args, *filters.MinRent)
		argIndex++
	}

	if filters.MaxRent != nil {
		baseQuery += fmt.Sprintf(" AND p.rent_amount <= $%d", argIndex)
		args = append(args, *filters.MaxRent)
		argIndex++
	}

	if filters.MinBedrooms != nil {
		baseQuery += fmt.Sprintf(" AND p.bedrooms >= $%d", argIndex)
		args = append(args, *filters.MinBedrooms)
		argIndex++
	}

	if filters.MaxBedrooms != nil {
		baseQuery += fmt.Sprintf(" AND p.bedrooms <= $%d", argIndex)
		args = append(args, *filters.MaxBedrooms)
		argIndex++
	}

	if filters.MinBathrooms != nil {
		baseQuery += fmt.Sprintf(" AND p.bathrooms >= $%d", argIndex)
		args = append(args, *filters.MinBathrooms)
		argIndex++
	}

	if filters.IsFurnished != nil {
		baseQuery += fmt.Sprintf(" AND p.is_furnished = $%d", argIndex)
		args = append(args, *filters.IsFurnished)
		argIndex++
	}

	if filters.HasParkingSpaces != nil && *filters.HasParkingSpaces {
		baseQuery += " AND p.parking_spaces > 0"
	}

	if filters.IsAvailable != nil {
		baseQuery += fmt.Sprintf(" AND p.is_available = $%d", argIndex)
		args = append(args, *filters.IsAvailable)
		argIndex++
	}

	baseQuery += " ORDER BY p.created_at DESC"

	// Add pagination
	limit := 20
	if filters.Limit > 0 {
		limit = filters.Limit
	}
	baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, limit)
	argIndex++

	if filters.Offset > 0 {
		baseQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []*Property
	for rows.Next() {
		property := &Property{}
		err := rows.Scan(
			&property.ID,
			&property.LandlordID,
			&property.Title,
			&property.Description,
			&property.PropertyType,
			&property.Bedrooms,
			&property.Bathrooms,
			&property.SquareMeters,
			&property.RentAmount,
			&property.DepositAmount,
			&property.CountyID,
			&property.SubCountyID,
			&property.LocationDetails,
			&property.Latitude,
			&property.Longitude,
			&property.Amenities,
			&property.UtilitiesIncluded,
			&property.ParkingSpaces,
			&property.IsFurnished,
			&property.IsAvailable,
			&property.AvailabilityDate,
			&property.CreatedAt,
			&property.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}

	return properties, nil
}

