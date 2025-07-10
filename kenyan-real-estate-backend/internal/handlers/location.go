package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"kenyan-real-estate-backend/internal/models"

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
func (h *LocationHandler) GetSubCounties(c *gin.Context) {
	countyIDStr := c.Param("county_id")
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

