// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"net/http"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LandmarkHandler struct {
	landmarkService services.LandmarkService
}

func NewLandmarkHandler(landmarkService services.LandmarkService) *LandmarkHandler {
	return &LandmarkHandler{landmarkService: landmarkService}
}

func (h *LandmarkHandler) CheckIn(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	landmarkIDStr := c.Param("landmark_id")
	landmarkID, err := uuid.Parse(landmarkIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid landmark_id format"})
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
		IsMocked  bool    `json:"is_mocked"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.landmarkService.CheckIn(c.Request.Context(), userID, landmarkID, req.Latitude, req.Longitude, req.IsMocked)
	if err != nil {
		if err == services.ErrFakeGPS {
			c.JSON(http.StatusForbidden, gin.H{
				"code":      403,
				"message":   "ERR_FAKE_GPS_DETECTED",
			})
			return
		}
		if err == services.ErrOutOfRange {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "ERR_OUT_OF_RANGE",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Check-in successful. +10 points",
	})
}
