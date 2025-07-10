package handlers

import (
	"net/http"

	"kenyan-real-estate-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// KenyanFeaturesHandler handles Kenyan-specific features and data
type KenyanFeaturesHandler struct{}

// NewKenyanFeaturesHandler creates a new Kenyan features handler
func NewKenyanFeaturesHandler() *KenyanFeaturesHandler {
	return &KenyanFeaturesHandler{}
}

// GetAmenities returns all available property amenities
// @Summary Get property amenities
// @Description Get all available property amenities, optionally filtered by category
// @Tags Kenyan Features
// @Accept json
// @Produce json
// @Param category query string false "Filter by category"
// @Success 200 {object} object{amenities=object} "List of amenities"
// @Router /amenities [get]
func (h *KenyanFeaturesHandler) GetAmenities(c *gin.Context) {
	category := c.Query("category")
	
	if category != "" {
		amenities := utils.GetAmenitiesByCategory(category)
		c.JSON(http.StatusOK, gin.H{
			"category":  category,
			"amenities": amenities,
		})
		return
	}

	allAmenities := utils.GetAllAmenities()
	c.JSON(http.StatusOK, gin.H{
		"amenities": allAmenities,
	})
}

// GetPropertyTypes returns all property types with descriptions
// @Summary Get property types
// @Description Get all available property types with their descriptions
// @Tags Kenyan Features
// @Accept json
// @Produce json
// @Success 200 {object} object{property_types=object} "List of property types"
// @Router /property-types [get]
func (h *KenyanFeaturesHandler) GetPropertyTypes(c *gin.Context) {
	propertyTypes := make(map[string]interface{})
	
	for propertyType := range utils.KenyanPropertyTypes {
		propertyTypes[propertyType] = map[string]string{
			"name":        propertyType,
			"description": utils.GetPropertyTypeDescription(propertyType),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"property_types": propertyTypes,
	})
}

// GetUtilities returns default utilities information
// @Summary Get utilities information
// @Description Get default utilities information for properties
// @Tags Kenyan Features
// @Accept json
// @Produce json
// @Success 200 {object} object{utilities=object} "List of utilities"
// @Router /utilities [get]
func (h *KenyanFeaturesHandler) GetUtilities(c *gin.Context) {
	utilities := utils.GetDefaultUtilities()
	c.JSON(http.StatusOK, gin.H{
		"utilities": utilities,
	})
}

// GetRentalTerms returns common rental terms in Kenya
// @Summary Get rental terms
// @Description Get common rental terms used in Kenya
// @Tags Kenyan Features
// @Accept json
// @Produce json
// @Success 200 {object} object{rental_terms=object} "List of rental terms"
// @Router /rental-terms [get]
func (h *KenyanFeaturesHandler) GetRentalTerms(c *gin.Context) {
	terms := utils.GetRentalTerms()
	c.JSON(http.StatusOK, gin.H{
		"rental_terms": terms,
	})
}

// GetPopularAreas returns popular residential areas by county
// @Summary Get popular areas
// @Description Get popular residential areas, optionally filtered by county
// @Tags Kenyan Features
// @Accept json
// @Produce json
// @Param county query string false "Filter by county name"
// @Success 200 {object} object{areas=object} "List of popular areas"
// @Failure 404 {object} object{error=string} "County not found"
// @Router /popular-areas [get]
func (h *KenyanFeaturesHandler) GetPopularAreas(c *gin.Context) {
	county := c.Query("county")
	areas := utils.GetPopularKenyanAreas()
	
	if county != "" {
		if countyAreas, exists := areas[county]; exists {
			c.JSON(http.StatusOK, gin.H{
				"county": county,
				"areas":  countyAreas,
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": "County not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"areas": areas,
	})
}

// ValidatePhoneNumber validates a Kenyan phone number
// @Summary Validate phone number
// @Description Validate a Kenyan phone number format
// @Tags Kenyan Features
// @Accept json
// @Produce json
// @Param phone_data body object{phone_number=string} true "Phone number to validate"
// @Success 200 {object} object{phone_number=string,is_valid=bool} "Validation result"
// @Failure 400 {object} object{error=string,details=string} "Invalid request data"
// @Router /validate-phone [post]
func (h *KenyanFeaturesHandler) ValidatePhoneNumber(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	isValid := utils.ValidateKenyanPhoneNumber(req.PhoneNumber)
	c.JSON(http.StatusOK, gin.H{
		"phone_number": req.PhoneNumber,
		"is_valid":     isValid,
	})
}

// FormatCurrency formats amount in Kenyan Shillings
// @Summary Format currency
// @Description Format an amount in Kenyan Shillings (KES)
// @Tags Kenyan Features
// @Accept json
// @Produce json
// @Param amount_data body object{amount=number} true "Amount to format"
// @Success 200 {object} object{amount=number,formatted=string} "Formatted currency"
// @Failure 400 {object} object{error=string,details=string} "Invalid request data"
// @Router /format-currency [post]
func (h *KenyanFeaturesHandler) FormatCurrency(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	formatted := utils.FormatKenyanCurrency(req.Amount)
	c.JSON(http.StatusOK, gin.H{
		"amount":    req.Amount,
		"formatted": formatted,
	})
}

