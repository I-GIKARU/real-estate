package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// County represents a Kenyan county
type County struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex"`
	Code      string    `json:"code" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// Relationships
	SubCounties []SubCounty `json:"sub_counties,omitempty" gorm:"foreignKey:CountyID"`
	Properties  []Property  `json:"properties,omitempty" gorm:"foreignKey:CountyID"`
}

// SubCounty represents a sub-county within a county
type SubCounty struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CountyID  int       `json:"county_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// Relationships
	County     *County    `json:"county,omitempty" gorm:"foreignKey:CountyID"`
	Properties []Property `json:"properties,omitempty" gorm:"foreignKey:SubCountyID"`
}

// PropertyImage represents an image associated with a property
type PropertyImage struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PropertyID     uuid.UUID `json:"property_id" gorm:"type:uuid;not null"`
	ImageURL       string    `json:"image_url" gorm:"not null"`
	SecureURL      string    `json:"secure_url" gorm:"not null"`
	PublicID       string    `json:"public_id" gorm:"not null"`
	Caption        *string   `json:"caption,omitempty"`
	IsPrimary      bool      `json:"is_primary" gorm:"default:false"`
	DisplayOrder   int       `json:"display_order" gorm:"default:0"`
	Width          *int      `json:"width,omitempty"`
	Height         *int      `json:"height,omitempty"`
	Format         *string   `json:"format,omitempty"`
	Bytes          *int      `json:"bytes,omitempty"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// Relationships
	Property *Property `json:"property,omitempty" gorm:"foreignKey:PropertyID"`
}

// CreatePropertyImageRequest represents the request to add an image to a property
type CreatePropertyImageRequest struct {
	Caption      *string `form:"caption,omitempty"`
	IsPrimary    bool    `form:"is_primary"`
	DisplayOrder int     `form:"display_order"`
}

// BeforeCreate GORM hook for PropertyImage
func (pi *PropertyImage) BeforeCreate(tx *gorm.DB) error {
	if pi.ID == uuid.Nil {
		pi.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for County model
func (County) TableName() string {
	return "counties"
}

// TableName returns the table name for SubCounty model
func (SubCounty) TableName() string {
	return "sub_counties"
}

// TableName returns the table name for PropertyImage model
func (PropertyImage) TableName() string {
	return "property_images"
}

// CountyRepository handles database operations for counties
type CountyRepository struct {
	db *gorm.DB
}

// NewCountyRepository creates a new county repository
func NewCountyRepository(db *gorm.DB) *CountyRepository {
	return &CountyRepository{db: db}
}

// GetAll retrieves all counties
func (r *CountyRepository) GetAll() ([]*County, error) {
	var counties []*County
	err := r.db.Order("name").Find(&counties).Error
	return counties, err
}

// GetByID retrieves a county by ID
func (r *CountyRepository) GetByID(id int) (*County, error) {
	var county County
	err := r.db.First(&county, id).Error
	if err != nil {
		return nil, err
	}
	return &county, nil
}

// SubCountyRepository handles database operations for sub-counties
type SubCountyRepository struct {
	db *gorm.DB
}

// NewSubCountyRepository creates a new sub-county repository
func NewSubCountyRepository(db *gorm.DB) *SubCountyRepository {
	return &SubCountyRepository{db: db}
}

// GetByCountyID retrieves sub-counties by county ID
func (r *SubCountyRepository) GetByCountyID(countyID int) ([]*SubCounty, error) {
	var subCounties []*SubCounty
	err := r.db.Where("county_id = ?", countyID).Order("name").Find(&subCounties).Error
	return subCounties, err
}

// GetByID retrieves a sub-county by ID
func (r *SubCountyRepository) GetByID(id int) (*SubCounty, error) {
	var subCounty SubCounty
	err := r.db.Preload("County").First(&subCounty, id).Error
	if err != nil {
		return nil, err
	}
	return &subCounty, nil
}

// PropertyImageRepository handles database operations for property images
type PropertyImageRepository struct {
	db *gorm.DB
}

// NewPropertyImageRepository creates a new property image repository
func NewPropertyImageRepository(db *gorm.DB) *PropertyImageRepository {
	return &PropertyImageRepository{db: db}
}

// Create creates a new property image
func (r *PropertyImageRepository) Create(image *PropertyImage) error {
	return r.db.Create(image).Error
}

// GetByPropertyID retrieves images for a property
func (r *PropertyImageRepository) GetByPropertyID(propertyID uuid.UUID) ([]*PropertyImage, error) {
	var images []*PropertyImage
	err := r.db.Where("property_id = ?", propertyID).
		Order("is_primary DESC, display_order ASC, created_at ASC").
		Find(&images).Error
	return images, err
}

// GetByID retrieves a property image by ID
func (r *PropertyImageRepository) GetByID(id uuid.UUID) (*PropertyImage, error) {
	var image PropertyImage
	err := r.db.First(&image, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// Delete deletes a property image
func (r *PropertyImageRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&PropertyImage{}, id).Error
}

// SetPrimary sets an image as primary and unsets others for the same property
func (r *PropertyImageRepository) SetPrimary(imageID uuid.UUID, propertyID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Unset all primary images for the property
		if err := tx.Model(&PropertyImage{}).Where("property_id = ?", propertyID).Update("is_primary", false).Error; err != nil {
			return err
		}
		
		// Set the specified image as primary
		return tx.Model(&PropertyImage{}).Where("id = ?", imageID).Update("is_primary", true).Error
	})
}

