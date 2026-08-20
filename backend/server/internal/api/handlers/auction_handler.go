// Copyright (c) 2026 Q-Tech. All rights reserved.
// Licensed under the GNU AGPLv3 License.

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/repository"
	"github.com/Cong-ty-TNNH-Q-Tech/Q-Love/backend/server/internal/services"
)

type AuctionHandler struct {
	auctionService services.AuctionService
	auctionRepo    repository.AuctionRepository
}

func NewAuctionHandler(auctionService services.AuctionService, auctionRepo repository.AuctionRepository) *AuctionHandler {
	return &AuctionHandler{
		auctionService: auctionService,
		auctionRepo:    auctionRepo,
	}
}

type BidRequest struct {
	Amount float64 `json:"amount"`
}

func (h *AuctionHandler) GetActiveAuctions(c *fiber.Ctx) error {
	ctx := c.Context()
	offset := c.QueryInt("offset", 0)
	limit := c.QueryInt("limit", 100)
	
	auctions, err := h.auctionRepo.GetActiveAuctions(ctx, offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"auctions": auctions})
}

func (h *AuctionHandler) PlaceBid(c *fiber.Ctx) error {
	ctx := c.Context()
	auctionIDStr := c.Params("id")
	auctionID, err := uuid.Parse(auctionIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid auction id"})
	}

	var req BidRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Assuming we extract UserID from JWT middleware
	userIDVal := c.Locals("userID")
	if userIDVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := userIDVal.(uuid.UUID)

	if err := h.auctionService.PlaceBid(ctx, auctionID, userID, req.Amount); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "bid placed successfully"})
}
