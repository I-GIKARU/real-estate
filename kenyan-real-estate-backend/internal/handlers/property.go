package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"kenyan-real-estate-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PropertyHandler handles property-related HTTP requests
type PropertyHandler struct {
	propertyRepo      *models.PropertyRepository
	propertyImageRepo *models.PropertyImageRepository
}

// NewPropertyHandler creates a new property handler
func NewPropertyHandler(propertyRepo *models.PropertyRepository, propertyImageRepo *models.PropertyImageRepository) *PropertyHandler {
	return &PropertyHandler{
		propertyRepo:      propertyRepo,
		propertyImageRepo: propertyImageRepo,
	}
}

// CreateProperty handles property creation (landlord only)
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	landlordID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req models.CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Create property
	property := &models.Property{
		LandlordID:        landlordID,
		Title:             req.Title,
		Description:       req.Description,
		PropertyType:      req.PropertyType,
		Bedrooms:          req.Bedrooms,
		Bathrooms:         req.Bathrooms,
		SquareMeters:      req.SquareMeters,
		RentAmount:        req.RentAmount,
		DepositAmount:     req.DepositAmount,
		CountyID:          req.CountyID,
		SubCountyID:       req.SubCountyID,
		LocationDetails:   req.LocationDetails,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		Amenities:         req.Amenities,
		UtilitiesIncluded: req.UtilitiesIncluded,
		ParkingSpaces:     req.ParkingSpaces,
		IsFurnished:       req.IsFurnished,
		AvailabilityDate:  req.AvailabilityDate,
	}

	if err := h.propertyRepo.Create(property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create property",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Property created successfully",
		"property": property,
	})
}

// GetProperty handles getting a single property by ID
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	propertyIDStr := c.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid property ID",
		})
		return
	}

	property, err := h.propertyRepo.GetByID(propertyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Property not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get property",
		})
		return
	}

	// Get property images
	images, err := h.propertyImageRepo.GetByPropertyID(propertyID)
	if err != nil {
		// Log error but don't fail the request
		images = []*models.PropertyImage{}
	}
	property.Images = images

	c.JSON(http.StatusOK, gin.H{
		"property": property,
	})
}

// GetPublicProperties handles getting public property listings with search and filtering
func (h *PropertyHandler) GetPublicProperties(c *gin.Context) {
	// Parse query parameters for filtering
	filters := &models.PropertySearchFilters{}

	if countyIDStr := c.Query("county_id"); countyIDStr != "" {
		if countyID, err := strconv.Atoi(countyIDStr); err == nil {
			filters.CountyID = &countyID
		}
	}

	if subCountyIDStr := c.Query("sub_county_id"); subCountyIDStr != "" {
		if subCountyID, err := strconv.Atoi(subCountyIDStr); err == nil {
			filters.SubCountyID = &subCountyID
		}
	}

	if propertyTypeStr := c.Query("property_type"); propertyTypeStr != "" {
		propertyType := models.PropertyType(propertyTypeStr)
		filters.PropertyType = &propertyType
	}

	if minRentStr := c.Query("min_rent"); minRentStr != "" {
		if minRent, err := strconv.ParseFloat(minRentStr, 64); err == nil {
			filters.MinRent = &minRent
		}
	}

	if maxRentStr := c.Query("max_rent"); maxRentStr != "" {
		if maxRent, err := strconv.ParseFloat(maxRentStr, 64); err == nil {
			filters.MaxRent = &maxRent
		}
	}

	if minBedroomsStr := c.Query("min_bedrooms"); minBedroomsStr != "" {
		if minBedrooms, err := strconv.Atoi(minBedroomsStr); err == nil {
			filters.MinBedrooms = &minBedrooms
		}
	}

	if maxBedroomsStr := c.Query("max_bedrooms"); maxBedroomsStr != "" {
		if maxBedrooms, err := strconv.Atoi(maxBedroomsStr); err == nil {
			filters.MaxBedrooms = &maxBedrooms
		}
	}

	if minBathroomsStr := c.Query("min_bathrooms"); minBathroomsStr != "" {
		if minBathrooms, err := strconv.Atoi(minBathroomsStr); err == nil {
			filters.MinBathrooms = &minBathrooms
		}
	}

	if isFurnishedStr := c.Query("is_furnished"); isFurnishedStr != "" {
		if isFurnished, err := strconv.ParseBool(isFurnishedStr); err == nil {
			filters.IsFurnished = &isFurnished
		}
	}

	if hasParkingStr := c.Query("has_parking"); hasParkingStr != "" {
		if hasParking, err := strconv.ParseBool(hasParkingStr); err == nil {
			filters.HasParkingSpaces = &hasParking
		}
	}

	// Pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}
	if filters.Limit == 0 {
		filters.Limit = 20 // Default limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Only show available properties for public listings
	isAvailable := true
	filters.IsAvailable = &isAvailable

	properties, err := h.propertyRepo.Search(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search properties",
		})
		return
	}

	// Get images for each property
	for _, property := range properties {
		images, err := h.propertyImageRepo.GetByPropertyID(property.ID)
		if err == nil {
			property.Images = images
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"properties": properties,
		"filters":    filters,
	})
}

// GetMyProperties handles getting properties for the authenticated landlord
func (h *PropertyHandler) GetMyProperties(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	landlordID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Pagination
	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	properties, err := h.propertyRepo.GetByLandlordID(landlordID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get properties",
		})
		return
	}

	// Get images for each property
	for _, property := range properties {
		images, err := h.propertyImageRepo.GetByPropertyID(property.ID)
		if err == nil {
			property.Images = images
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"properties": properties,
	})
}

// UpdateProperty handles property updates (landlord only)
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	landlordID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	propertyIDStr := c.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid property ID",
		})
		return
	}

	// Get existing property
	property, err := h.propertyRepo.GetByID(propertyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Property not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get property",
		})
		return
	}

	// Check if user owns the property
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only update your own properties",
		})
		return
	}

	var req models.UpdatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Update fields if provided
	if req.Title != nil {
		property.Title = *req.Title
	}
	if req.Description != nil {
		property.Description = req.Description
	}
	if req.Bedrooms != nil {
		property.Bedrooms = *req.Bedrooms
	}
	if req.Bathrooms != nil {
		property.Bathrooms = *req.Bathrooms
	}
	if req.SquareMeters != nil {
		property.SquareMeters = req.SquareMeters
	}
	if req.RentAmount != nil {
		property.RentAmount = *req.RentAmount
	}
	if req.DepositAmount != nil {
		property.DepositAmount = req.DepositAmount
	}
	if req.LocationDetails != nil {
		property.LocationDetails = req.LocationDetails
	}
	if req.Latitude != nil {
		property.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		property.Longitude = req.Longitude
	}
	if req.Amenities != nil {
		property.Amenities = *req.Amenities
	}
	if req.UtilitiesIncluded != nil {
		property.UtilitiesIncluded = *req.UtilitiesIncluded
	}
	if req.ParkingSpaces != nil {
		property.ParkingSpaces = *req.ParkingSpaces
	}
	if req.IsFurnished != nil {
		property.IsFurnished = *req.IsFurnished
	}
	if req.IsAvailable != nil {
		property.IsAvailable = *req.IsAvailable
	}
	if req.AvailabilityDate != nil {
		property.AvailabilityDate = req.AvailabilityDate
	}

	if err := h.propertyRepo.Update(property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update property",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Property updated successfully",
		"property": property,
	})
}

// DeleteProperty handles property deletion (landlord only)
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	landlordID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	propertyIDStr := c.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid property ID",
		})
		return
	}

	// Get existing property to check ownership
	property, err := h.propertyRepo.GetByID(propertyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Property not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get property",
		})
		return
	}

	// Check if user owns the property
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own properties",
		})
		return
	}

	if err := h.propertyRepo.Delete(propertyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete property",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Property deleted successfully",
	})
}

// AddPropertyImage handles adding images to a property
func (h *PropertyHandler) AddPropertyImage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	landlordID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	propertyIDStr := c.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid property ID",
		})
		return
	}

	// Check if user owns the property
	property, err := h.propertyRepo.GetByID(propertyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Property not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get property",
		})
		return
	}

	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only add images to your own properties",
		})
		return
	}

	var req models.CreatePropertyImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	image := &models.PropertyImage{
		PropertyID:   propertyID,
		ImageURL:     req.ImageURL,
		Caption:      req.Caption,
		IsPrimary:    req.IsPrimary,
		DisplayOrder: req.DisplayOrder,
	}

	if err := h.propertyImageRepo.Create(image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add image",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Image added successfully",
		"image":   image,
	})
}

