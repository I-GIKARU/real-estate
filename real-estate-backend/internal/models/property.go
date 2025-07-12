package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	ID                uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AgentID           uuid.UUID         `json:"agent_id" gorm:"type:uuid;not null"`
	Title             string            `json:"title" gorm:"not null"`
	Description       *string           `json:"description,omitempty"`
	PropertyType      PropertyType      `json:"property_type" gorm:"not null;type:varchar(20)"`
	Bedrooms          int               `json:"bedrooms" gorm:"not null"`
	Bathrooms         int               `json:"bathrooms" gorm:"not null"`
	SquareMeters      *float64          `json:"square_meters,omitempty"`
	RentAmount        float64           `json:"rent_amount" gorm:"not null"`
	DepositAmount     *float64          `json:"deposit_amount,omitempty"`
	CountyID          int               `json:"county_id" gorm:"not null"`
	SubCountyID       *int              `json:"sub_county_id,omitempty"`
	LocationDetails   *string           `json:"location_details,omitempty"`
	Latitude          *float64          `json:"latitude,omitempty"`
	Longitude         *float64          `json:"longitude,omitempty"`
	Amenities         Amenities         `json:"amenities" gorm:"type:jsonb"`
	UtilitiesIncluded UtilitiesIncluded `json:"utilities_included" gorm:"type:jsonb"`
	ParkingSpaces     int               `json:"parking_spaces"`
	IsFurnished       bool              `json:"is_furnished" gorm:"default:false"`
	IsAvailable       bool              `json:"is_available" gorm:"default:true"`
	AvailabilityDate  *time.Time        `json:"availability_date,omitempty"`
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt    `json:"-" gorm:"index"`

	// Relationships
	County    *County         `json:"county,omitempty" gorm:"foreignKey:CountyID"`
	SubCounty *SubCounty      `json:"sub_county,omitempty" gorm:"foreignKey:SubCountyID"`
	Agent     *User           `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	Images    []*PropertyImage `json:"images,omitempty" gorm:"foreignKey:PropertyID"`
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

// BeforeCreate GORM hook to set ID
func (p *Property) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for Property model
func (Property) TableName() string {
	return "properties"
}

// PropertyRepository handles database operations for properties
type PropertyRepository struct {
	db *gorm.DB
}

// NewPropertyRepository creates a new property repository
func NewPropertyRepository(db *gorm.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

// Create creates a new property
func (r *PropertyRepository) Create(property *Property) error {
	return r.db.Create(property).Error
}

// GetByID retrieves a property by ID
func (r *PropertyRepository) GetByID(id uuid.UUID) (*Property, error) {
	var property Property
	err := r.db.Preload("County").Preload("SubCounty").Preload("Agent").Preload("Images").First(&property, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}

// GetByAgentID retrieves properties by agent ID
func (r *PropertyRepository) GetByAgentID(agentID uuid.UUID, limit, offset int) ([]*Property, error) {
	var properties []*Property
	query := r.db.Preload("County").Preload("SubCounty").Preload("Agent").Preload("Images").Where("agent_id = ?", agentID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&properties)
	return properties, result.Error
}

// Update updates a property
func (r *PropertyRepository) Update(property *Property) error {
	return r.db.Save(property).Error
}

// Delete deletes a property
func (r *PropertyRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&Property{}, id).Error
}

// Search searches for properties based on filters
func (r *PropertyRepository) Search(filters *PropertySearchFilters) ([]*Property, error) {
	var properties []*Property
	query := r.db.Model(&Property{}).Preload("County").Preload("SubCounty").Preload("Agent").Preload("Images")

	if filters.CountyID != nil {
		query = query.Where("county_id = ?", *filters.CountyID)
	}

	if filters.SubCountyID != nil {
		query = query.Where("sub_county_id = ?", *filters.SubCountyID)
	}

	if filters.PropertyType != nil {
		query = query.Where("property_type = ?", *filters.PropertyType)
	}

	if filters.MinRent != nil {
		query = query.Where("rent_amount >= ?", *filters.MinRent)
	}

	if filters.MaxRent != nil {
		query = query.Where("rent_amount <= ?", *filters.MaxRent)
	}

	if filters.MinBedrooms != nil {
		query = query.Where("bedrooms >= ?", *filters.MinBedrooms)
	}

	if filters.MaxBedrooms != nil {
		query = query.Where("bedrooms <= ?", *filters.MaxBedrooms)
	}

	if filters.MinBathrooms != nil {
		query = query.Where("bathrooms >= ?", *filters.MinBathrooms)
	}

	if filters.IsFurnished != nil {
		query = query.Where("is_furnished = ?", *filters.IsFurnished)
	}

	if filters.HasParkingSpaces != nil && *filters.HasParkingSpaces {
		query = query.Where("parking_spaces > 0")
	}

	if filters.IsAvailable != nil {
		query = query.Where("is_available = ?", *filters.IsAvailable)
	}

	query = query.Order("created_at DESC")

	// Set default limit if not provided
	limit := 20
	if filters.Limit > 0 {
		limit = filters.Limit
	}
	query = query.Limit(limit)

	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&properties).Error
	return properties, err
}
