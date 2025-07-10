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
func (h *KenyanFeaturesHandler) GetUtilities(c *gin.Context) {
	utilities := utils.GetDefaultUtilities()
	c.JSON(http.StatusOK, gin.H{
		"utilities": utilities,
	})
}

// GetRentalTerms returns common rental terms in Kenya
func (h *KenyanFeaturesHandler) GetRentalTerms(c *gin.Context) {
	terms := utils.GetRentalTerms()
	c.JSON(http.StatusOK, gin.H{
		"rental_terms": terms,
	})
}

// GetPopularAreas returns popular residential areas by county
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

