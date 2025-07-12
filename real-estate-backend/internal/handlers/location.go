package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"real-estate-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// LocationHandler handles location-related HTTP requests
type LocationHandler struct {
	countyRepo    *models.CountyRepository
	subCountyRepo *models.SubCountyRepository
}

// NewLocationHandler creates a new location handler
func NewLocationHandler(countyRepo *models.CountyRepository, subCountyRepo *models.SubCountyRepository) *LocationHandler {
	return &LocationHandler{
		countyRepo:    countyRepo,
		subCountyRepo: subCountyRepo,
	}
}

// GetCounties handles getting all counties
// @Summary Get all counties
// @Description Get a list of all counties in Kenya
// @Tags Location
// @Accept json
// @Produce json
// @Success 200 {object} object{counties=[]models.County} "List of counties"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /counties [get]
func (h *LocationHandler) GetCounties(c *gin.Context) {
	counties, err := h.countyRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get counties",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"counties": counties,
	})
}

// GetCounty handles getting a single county by ID
// @Summary Get county by ID
// @Description Get detailed information about a specific county
// @Tags Location
// @Accept json
// @Produce json
// @Param id path int true "County ID"
// @Success 200 {object} object{county=models.County} "County details"
// @Failure 400 {object} object{error=string} "Invalid county ID"
// @Failure 404 {object} object{error=string} "County not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /counties/{id} [get]
func (h *LocationHandler) GetCounty(c *gin.Context) {
	countyIDStr := c.Param("id")
	countyID, err := strconv.Atoi(countyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid county ID",
		})
		return
	}

	county, err := h.countyRepo.GetByID(countyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "County not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get county",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"county": county,
	})
}

// GetSubCounties handles getting sub-counties by county ID
// @Summary Get sub-counties by county ID
// @Description Get all sub-counties within a specific county
// @Tags Location
// @Accept json
// @Produce json
// @Param id path int true "County ID"
// @Success 200 {object} object{sub_counties=[]models.SubCounty,county_id=int} "List of sub-counties"
// @Failure 400 {object} object{error=string} "Invalid county ID"
// @Failure 404 {object} object{error=string} "County not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /counties/{id}/sub-counties [get]
func (h *LocationHandler) GetSubCounties(c *gin.Context) {
	countyIDStr := c.Param("id")
	countyID, err := strconv.Atoi(countyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid county ID",
		})
		return
	}

	// Verify county exists
	_, err = h.countyRepo.GetByID(countyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "County not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to verify county",
		})
		return
	}

	subCounties, err := h.subCountyRepo.GetByCountyID(countyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sub-counties",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sub_counties": subCounties,
		"county_id":    countyID,
	})
}

// GetSubCounty handles getting a single sub-county by ID
// @Summary Get sub-county by ID
// @Description Get detailed information about a specific sub-county
// @Tags Location
// @Accept json
// @Produce json
// @Param id path int true "Sub-county ID"
// @Success 200 {object} object{sub_county=models.SubCounty} "Sub-county details"
// @Failure 400 {object} object{error=string} "Invalid sub-county ID"
// @Failure 404 {object} object{error=string} "Sub-county not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /sub-counties/{id} [get]
func (h *LocationHandler) GetSubCounty(c *gin.Context) {
	subCountyIDStr := c.Param("id")
	subCountyID, err := strconv.Atoi(subCountyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid sub-county ID",
		})
		return
	}

	subCounty, err := h.subCountyRepo.GetByID(subCountyID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Sub-county not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get sub-county",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sub_county": subCounty,
	})
}

