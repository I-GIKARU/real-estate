package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"real-estate-backend/internal/config"
	"real-estate-backend/internal/models"
	"real-estate-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PropertyHandler handles property-related HTTP requests
type PropertyHandler struct {
	propertyRepo      *models.PropertyRepository
	propertyImageRepo *models.PropertyImageRepository
	cloudinaryService *services.CloudinaryService
	uploadConfig      *config.UploadConfig
}

// NewPropertyHandler creates a new property handler
func NewPropertyHandler(propertyRepo *models.PropertyRepository, propertyImageRepo *models.PropertyImageRepository, cloudinaryService *services.CloudinaryService, uploadConfig *config.UploadConfig) *PropertyHandler {
	return &PropertyHandler{
		propertyRepo:      propertyRepo,
		propertyImageRepo: propertyImageRepo,
		cloudinaryService: cloudinaryService,
		uploadConfig:      uploadConfig,
	}
}

// CreateProperty handles property creation (agent only)
// @Summary Create a new property
// @Description Create a new property listing (agent only)
// @Tags Properties
// @Accept json
// @Produce json
// @Security Bearer
// @Param property body models.CreatePropertyRequest true "Property data"
// @Success 201 {object} object{message=string,property=models.Property} "Property created successfully"
// @Failure 400 {object} object{error=string,details=string} "Invalid request data"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /properties [post]
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	agentID, ok := userID.(uuid.UUID)
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
		AgentID:           agentID,
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
// @Summary Get a property by ID
// @Description Get detailed information about a specific property including images
// @Tags Properties
// @Accept json
// @Produce json
// @Param id path string true "Property ID" Format(uuid)
// @Success 200 {object} object{property=models.Property} "Property details"
// @Failure 400 {object} object{error=string} "Invalid property ID"
// @Failure 404 {object} object{error=string} "Property not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /properties/{id} [get]
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
// @Summary Get public property listings
// @Description Get a list of available properties with optional filtering and pagination
// @Tags Properties
// @Accept json
// @Produce json
// @Param county_id query int false "Filter by county ID"
// @Param sub_county_id query int false "Filter by sub-county ID"
// @Param property_type query string false "Filter by property type" Enums(apartment,house,bedsitter,studio,maisonette,bungalow,villa,commercial)
// @Param min_rent query number false "Minimum rent amount"
// @Param max_rent query number false "Maximum rent amount"
// @Param min_bedrooms query int false "Minimum number of bedrooms"
// @Param max_bedrooms query int false "Maximum number of bedrooms"
// @Param min_bathrooms query int false "Minimum number of bathrooms"
// @Param is_furnished query boolean false "Filter by furnished status"
// @Param has_parking query boolean false "Filter by parking availability"
// @Param limit query int false "Number of results per page" default(20)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} object{properties=[]models.Property,total=int,limit=int,offset=int} "List of properties"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /properties [get]
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

// GetMyProperties handles getting properties for the authenticated agent
// @Summary Get my properties
// @Description Get all properties managed by the authenticated agent
// @Tags Properties
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query int false "Number of results per page" default(20)
// @Param offset query int false "Number of results to skip" default(0)
// @Success 200 {object} object{properties=[]models.Property} "List of landlord's properties"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /my-properties [get]
func (h *PropertyHandler) GetMyProperties(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	agentID, ok := userID.(uuid.UUID)
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

	properties, err := h.propertyRepo.GetByAgentID(agentID, limit, offset)
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
// @Summary Update property
// @Description Update property information (landlord only)
// @Tags Properties
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Property ID" Format(uuid)
// @Param property body models.UpdatePropertyRequest true "Property update data"
// @Success 200 {object} object{message=string,property=models.Property} "Property updated successfully"
// @Failure 400 {object} object{error=string,details=string} "Invalid request data"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "You can only update your own properties"
// @Failure 404 {object} object{error=string} "Property not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /properties/{id} [put]
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	agentID, ok := userID.(uuid.UUID)
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
	if property.AgentID != agentID {
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
// @Summary Delete property
// @Description Delete a property (landlord only)
// @Tags Properties
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Property ID" Format(uuid)
// @Success 200 {object} object{message=string} "Property deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid property ID"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "You can only delete your own properties"
// @Failure 404 {object} object{error=string} "Property not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /properties/{id} [delete]
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	agentID, ok := userID.(uuid.UUID)
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
	if property.AgentID != agentID {
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
// @Summary Add property image
// @Description Add an image to a property (landlord only)
// @Tags Properties
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param id path string true "Property ID" Format(uuid)
// @Param image formData file true "Image file"
// @Param caption formData string false "Image caption"
// @Param is_primary formData boolean false "Set as primary image"
// @Param display_order formData int false "Display order"
// @Success 201 {object} object{message=string,image=models.PropertyImage} "Image added successfully"
// @Failure 400 {object} object{error=string} "Invalid request data or file"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "You can only add images to your own properties"
// @Failure 404 {object} object{error=string} "Property not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /properties/{id}/images [post]
func (h *PropertyHandler) AddPropertyImage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	agentID, ok := userID.(uuid.UUID)
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

	if property.AgentID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only add images to your own properties",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image file is required",
		})
		return
	}

	defer file.Close()

	if err := h.cloudinaryService.ValidateImageFile(header, h.uploadConfig.AllowedTypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var req models.CreatePropertyImageRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid form data",
			"details": err.Error(),
		})
		return
	}

	uploadResponse, err := h.cloudinaryService.UploadImage(c.Request.Context(), file, header, propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	image := &models.PropertyImage{
		PropertyID:   propertyID,
		ImageURL:     uploadResponse.URL,
		SecureURL:    uploadResponse.SecureURL,
		PublicID:     uploadResponse.PublicID,
		Width:        &uploadResponse.Width,
		Height:       &uploadResponse.Height,
		Format:       &uploadResponse.Format,
		Bytes:        &uploadResponse.Bytes,
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

// DeletePropertyImage handles deleting a property image
// @Summary Delete property image
// @Description Delete an image from a property (landlord only)
// @Tags Properties
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Property ID" Format(uuid)
// @Param image_id path string true "Image ID" Format(uuid)
// @Success 200 {object} object{message=string} "Image deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid property or image ID"
// @Failure 401 {object} object{error=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "You can only delete images from your own properties"
// @Failure 404 {object} object{error=string} "Property or image not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /properties/{id}/images/{image_id} [delete]
func (h *PropertyHandler) DeletePropertyImage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	agentID, ok := userID.(uuid.UUID)
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

	imageIDStr := c.Param("image_id")
	imageID, err := uuid.Parse(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid image ID",
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

	if property.AgentID != agentID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete images from your own properties",
		})
		return
	}

	// Get the image to get the public_id for Cloudinary cleanup
	image, err := h.propertyImageRepo.GetByID(imageID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Image not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get image",
		})
		return
	}

	// Verify the image belongs to the property
	if image.PropertyID != propertyID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image does not belong to this property",
		})
		return
	}

	// Delete from Cloudinary
	if err := h.cloudinaryService.DeleteImage(c.Request.Context(), image.PublicID); err != nil {
		// Log the error but don't fail the request - the image might already be deleted
		// or the public_id might be invalid
	}

	// Delete from database
	if err := h.propertyImageRepo.Delete(imageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete image",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}

