package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// County represents a Kenyan county
type County struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Code      string    `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SubCounty represents a sub-county within a county
type SubCounty struct {
	ID        int       `json:"id" db:"id"`
	CountyID  int       `json:"county_id" db:"county_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	County    *County   `json:"county,omitempty"`
}

// PropertyImage represents an image associated with a property
type PropertyImage struct {
	ID           uuid.UUID `json:"id" db:"id"`
	PropertyID   uuid.UUID `json:"property_id" db:"property_id"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	Caption      *string   `json:"caption,omitempty" db:"caption"`
	IsPrimary    bool      `json:"is_primary" db:"is_primary"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// CreatePropertyImageRequest represents the request to add an image to a property
type CreatePropertyImageRequest struct {
	PropertyID   uuid.UUID `json:"property_id" binding:"required"`
	ImageURL     string    `json:"image_url" binding:"required"`
	Caption      *string   `json:"caption,omitempty"`
	IsPrimary    bool      `json:"is_primary"`
	DisplayOrder int       `json:"display_order"`
}

// CountyRepository handles database operations for counties
type CountyRepository struct {
	db *sql.DB
}

// NewCountyRepository creates a new county repository
func NewCountyRepository(db *sql.DB) *CountyRepository {
	return &CountyRepository{db: db}
}

// GetAll retrieves all counties
func (r *CountyRepository) GetAll() ([]*County, error) {
	query := `SELECT id, name, code, created_at FROM counties ORDER BY name`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counties []*County
	for rows.Next() {
		county := &County{}
		err := rows.Scan(&county.ID, &county.Name, &county.Code, &county.CreatedAt)
		if err != nil {
			return nil, err
		}
		counties = append(counties, county)
	}

	return counties, nil
}

// GetByID retrieves a county by ID
func (r *CountyRepository) GetByID(id int) (*County, error) {
	county := &County{}
	query := `SELECT id, name, code, created_at FROM counties WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(&county.ID, &county.Name, &county.Code, &county.CreatedAt)
	if err != nil {
		return nil, err
	}

	return county, nil
}

// SubCountyRepository handles database operations for sub-counties
type SubCountyRepository struct {
	db *sql.DB
}

// NewSubCountyRepository creates a new sub-county repository
func NewSubCountyRepository(db *sql.DB) *SubCountyRepository {
	return &SubCountyRepository{db: db}
}

// GetByCountyID retrieves sub-counties by county ID
func (r *SubCountyRepository) GetByCountyID(countyID int) ([]*SubCounty, error) {
	query := `SELECT id, county_id, name, created_at FROM sub_counties WHERE county_id = $1 ORDER BY name`
	
	rows, err := r.db.Query(query, countyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subCounties []*SubCounty
	for rows.Next() {
		subCounty := &SubCounty{}
		err := rows.Scan(&subCounty.ID, &subCounty.CountyID, &subCounty.Name, &subCounty.CreatedAt)
		if err != nil {
			return nil, err
		}
		subCounties = append(subCounties, subCounty)
	}

	return subCounties, nil
}

// GetByID retrieves a sub-county by ID
func (r *SubCountyRepository) GetByID(id int) (*SubCounty, error) {
	subCounty := &SubCounty{}
	query := `
		SELECT sc.id, sc.county_id, sc.name, sc.created_at,
			   c.id, c.name, c.code, c.created_at
		FROM sub_counties sc
		JOIN counties c ON sc.county_id = c.id
		WHERE sc.id = $1`
	
	county := &County{}
	err := r.db.QueryRow(query, id).Scan(
		&subCounty.ID, &subCounty.CountyID, &subCounty.Name, &subCounty.CreatedAt,
		&county.ID, &county.Name, &county.Code, &county.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	subCounty.County = county
	return subCounty, nil
}

// PropertyImageRepository handles database operations for property images
type PropertyImageRepository struct {
	db *sql.DB
}

// NewPropertyImageRepository creates a new property image repository
func NewPropertyImageRepository(db *sql.DB) *PropertyImageRepository {
	return &PropertyImageRepository{db: db}
}

// Create creates a new property image
func (r *PropertyImageRepository) Create(image *PropertyImage) error {
	query := `
		INSERT INTO property_images (id, property_id, image_url, caption, is_primary, display_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	image.ID = uuid.New()
	image.CreatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		image.ID,
		image.PropertyID,
		image.ImageURL,
		image.Caption,
		image.IsPrimary,
		image.DisplayOrder,
		image.CreatedAt,
	).Scan(&image.ID, &image.CreatedAt)

	return err
}

// GetByPropertyID retrieves images for a property
func (r *PropertyImageRepository) GetByPropertyID(propertyID uuid.UUID) ([]*PropertyImage, error) {
	query := `
		SELECT id, property_id, image_url, caption, is_primary, display_order, created_at
		FROM property_images
		WHERE property_id = $1
		ORDER BY is_primary DESC, display_order ASC, created_at ASC`
	
	rows, err := r.db.Query(query, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*PropertyImage
	for rows.Next() {
		image := &PropertyImage{}
		err := rows.Scan(
			&image.ID,
			&image.PropertyID,
			&image.ImageURL,
			&image.Caption,
			&image.IsPrimary,
			&image.DisplayOrder,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

// Delete deletes a property image
func (r *PropertyImageRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM property_images WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// SetPrimary sets an image as primary and unsets others for the same property
func (r *PropertyImageRepository) SetPrimary(imageID uuid.UUID, propertyID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unset all primary images for the property
	_, err = tx.Exec(`UPDATE property_images SET is_primary = false WHERE property_id = $1`, propertyID)
	if err != nil {
		return err
	}

	// Set the specified image as primary
	_, err = tx.Exec(`UPDATE property_images SET is_primary = true WHERE id = $1`, imageID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

