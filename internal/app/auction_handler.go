package app

import (
	"net/http"
	"strconv"

	"yourapp/internal/model"
	"yourapp/internal/repository"
	"yourapp/internal/service"

	"github.com/gin-gonic/gin"
)

type AuctionHandler struct {
	auctionService service.AuctionService
	jwtSecret      string
}

func NewAuctionHandler(auctionService service.AuctionService, jwtSecret string) *AuctionHandler {
	return &AuctionHandler{
		auctionService: auctionService,
		jwtSecret:      jwtSecret,
	}
}

// ========== SELLER HANDLERS ==========

func (h *AuctionHandler) CreateSeller(c *gin.Context) {
	var req service.CreateSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller, err := h.auctionService.CreateSeller(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": seller})
}

func (h *AuctionHandler) GetSellers(c *gin.Context) {
	sellers, err := h.auctionService.GetAllSellers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sellers})
}

func (h *AuctionHandler) GetSeller(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid seller id"})
		return
	}

	seller, err := h.auctionService.GetSeller(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "seller not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": seller})
}

// ========== ORGANIZER HANDLERS ==========

func (h *AuctionHandler) CreateOrganizer(c *gin.Context) {
	var req service.CreateOrganizerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organizer, err := h.auctionService.CreateOrganizer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": organizer})
}

func (h *AuctionHandler) GetOrganizers(c *gin.Context) {
	organizers, err := h.auctionService.GetAllOrganizers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": organizers})
}

func (h *AuctionHandler) GetOrganizer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organizer id"})
		return
	}

	organizer, err := h.auctionService.GetOrganizer(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organizer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": organizer})
}

// ========== CATEGORY HANDLERS ==========

func (h *AuctionHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.auctionService.CreateCategory(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

func (h *AuctionHandler) GetCategories(c *gin.Context) {
	categories, err := h.auctionService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (h *AuctionHandler) GetCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	category, err := h.auctionService.GetCategory(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// ========== AUCTION ITEM HANDLERS ==========

func (h *AuctionHandler) CreateAuctionItem(c *gin.Context) {
	var req service.CreateAuctionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.auctionService.CreateAuctionItem(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

func (h *AuctionHandler) GetAuctionItems(c *gin.Context) {
	// Parse query parameters
	filters := repository.AuctionItemFilters{
		Page:      1,
		Limit:     20,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Search:    c.Query("search"),
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filters.Page = page
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 {
		filters.Limit = limit
	}
	if categoryID, err := strconv.ParseUint(c.Query("category_id"), 10, 32); err == nil {
		id := uint(categoryID)
		filters.CategoryID = &id
	}

	items, total, err := h.auctionService.GetPublishedAuctions(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"meta": gin.H{
			"total":       total,
			"page":        filters.Page,
			"limit":       filters.Limit,
			"total_pages": (total + int64(filters.Limit) - 1) / int64(filters.Limit),
		},
	})
}

func (h *AuctionHandler) GetAuctionItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	item, err := h.auctionService.GetAuctionItem(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *AuctionHandler) UpdateAuctionItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req service.UpdateAuctionItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.auctionService.UpdateAuctionItem(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *AuctionHandler) PublishAuctionItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	if err := h.auctionService.PublishAuctionItem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item published successfully"})
}

func (h *AuctionHandler) DeleteAuctionItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	if err := h.auctionService.DeleteAuctionItem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}

// ========== BID HANDLERS ==========

func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	var req service.PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	req.UserID = userID.(string)
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.Request.UserAgent()

	bid, err := h.auctionService.PlaceBid(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": bid})
}

func (h *AuctionHandler) GetItemBids(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	bids, err := h.auctionService.GetItemBids(uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bids})
}

func (h *AuctionHandler) GetUserBids(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	bids, err := h.auctionService.GetUserBids(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bids})
}

// ========== RESPONSE TRANSFORMER ==========

// TransformAuctionItemForFrontend transforms an auction item to match frontend expectations
type AuctionItemResponse struct {
	ID            uint                     `json:"id"`
	LotCode       string                   `json:"lot_code"`
	Name          string                   `json:"name"`
	Category      string                   `json:"category"`
	Image         string                   `json:"image"`
	CurrentBid    float64                  `json:"current_bid"`
	PreviousBid   float64                  `json:"previous_bid"`
	StartingPrice float64                  `json:"starting_price"`
	TotalBids     int                      `json:"total_bids"`
	TimeLeft      string                   `json:"time_left"`
	IsHot         bool                     `json:"is_hot"`
	Status        model.AuctionStatus      `json:"status"`
	Description   string                   `json:"description"`
	Images        []string                 `json:"images"`
	Schedule      *AuctionScheduleResponse `json:"schedule,omitempty"`
}

type AuctionScheduleResponse struct {
	AuctionStart string `json:"auction_start"`
	AuctionEnd   string `json:"auction_end"`
}

func (h *AuctionHandler) GetAuctionItemsForFrontend(c *gin.Context) {
	filters := repository.AuctionItemFilters{
		Page:      1,
		Limit:     20,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Search:    c.Query("search"),
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filters.Page = page
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil && limit > 0 {
		filters.Limit = limit
	}
	if categoryID, err := strconv.ParseUint(c.Query("category_id"), 10, 32); err == nil {
		id := uint(categoryID)
		filters.CategoryID = &id
	}

	items, total, err := h.auctionService.GetPublishedAuctions(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform items for frontend
	var response []AuctionItemResponse
	for _, item := range items {
		transformed := transformAuctionItem(item)
		response = append(response, transformed)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
		"meta": gin.H{
			"total":       total,
			"page":        filters.Page,
			"limit":       filters.Limit,
			"total_pages": (total + int64(filters.Limit) - 1) / int64(filters.Limit),
		},
	})
}

func transformAuctionItem(item model.AuctionItem) AuctionItemResponse {
	currentBid, _ := item.CurrentHighestBid.Float64()
	startingPrice, _ := item.StartingPrice.Float64()

	// Calculate previous bid (current - increment)
	previousBid := currentBid
	if item.BidCount > 0 {
		increment, _ := item.IncrementAmount.Float64()
		previousBid = currentBid - increment
	}

	// Get main image
	mainImage := ""
	var allImages []string
	for _, img := range item.Images {
		allImages = append(allImages, img.ImageURL)
		if img.ImageType == model.ImageTypeMain || mainImage == "" {
			mainImage = img.ImageURL
		}
	}

	// Get category name
	categoryName := ""
	if item.Category != nil {
		categoryName = item.Category.CategoryName
	}

	// Get description
	description := ""
	if item.Description != nil {
		description = *item.Description
	}

	// Calculate time left
	timeLeft := "N/A"
	if item.Schedule != nil {
		timeLeft = calculateTimeLeft(item.Schedule.AuctionEnd)
	}

	// Determine if hot (more than 10 bids or ending soon)
	isHot := item.BidCount > 10

	resp := AuctionItemResponse{
		ID:            item.ID,
		LotCode:       item.LotCode,
		Name:          item.ItemName,
		Category:      categoryName,
		Image:         mainImage,
		CurrentBid:    currentBid,
		PreviousBid:   previousBid,
		StartingPrice: startingPrice,
		TotalBids:     item.BidCount,
		TimeLeft:      timeLeft,
		IsHot:         isHot,
		Status:        item.Status,
		Description:   description,
		Images:        allImages,
	}

	if item.Schedule != nil {
		resp.Schedule = &AuctionScheduleResponse{
			AuctionStart: item.Schedule.AuctionStart.Format("2006-01-02T15:04:05Z"),
			AuctionEnd:   item.Schedule.AuctionEnd.Format("2006-01-02T15:04:05Z"),
		}
	}

	return resp
}

func calculateTimeLeft(endTime interface{}) string {
	// Handle different time formats
	var end interface{}
	switch v := endTime.(type) {
	case string:
		end = v
	default:
		end = v
	}

	// For now, return a simple format
	// In production, calculate actual time difference
	_ = end
	return "2h 30m"
}
