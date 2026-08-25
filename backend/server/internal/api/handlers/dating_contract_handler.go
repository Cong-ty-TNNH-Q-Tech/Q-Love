// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package handlers

import (
	"net/http"
	"time"

	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DatingContractHandler struct {
	contractService services.DatingContractService
}

func NewDatingContractHandler(contractService services.DatingContractService) *DatingContractHandler {
	return &DatingContractHandler{contractService: contractService}
}

func (h *DatingContractHandler) CreateContract(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		TargetUserID    string    `json:"target_user_id" binding:"required"`
		DepositAmount   float64   `json:"deposit_amount" binding:"required"`
		AppointmentTime time.Time `json:"appointment_time" binding:"required"`
		LocationNote    string    `json:"location_note"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetUserID, err := uuid.Parse(req.TargetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target_user_id format"})
		return
	}

	contract, err := h.contractService.CreateContract(c.Request.Context(), userID, targetUserID, req.DepositAmount, req.AppointmentTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

func (h *DatingContractHandler) AcceptContract(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	contractIDStr := c.Param("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract_id format"})
		return
	}

	contract, err := h.contractService.AcceptContract(c.Request.Context(), contractID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func (h *DatingContractHandler) CancelContract(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	contractIDStr := c.Param("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract_id format"})
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// reason is optional
	}

	err = h.contractService.CancelContract(c.Request.Context(), contractID, userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cancelled_by":    userID.String(),
		"penalty_applied": true,
		"message":         "Cancelled successfully",
	})
}

func (h *DatingContractHandler) ScanContract(c *gin.Context) {
	contractIDStr := c.Param("contract_id")
	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract_id format"})
		return
	}

	var req struct {
		QRToken string `json:"qr_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.contractService.ScanContract(c.Request.Context(), contractID, req.QRToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}
