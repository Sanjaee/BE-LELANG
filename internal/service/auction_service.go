package service

import (
	"errors"
	"fmt"
	"time"

	"yourapp/internal/model"
	"yourapp/internal/repository"

	"github.com/shopspring/decimal"
)

type AuctionService interface {
	// Seller
	CreateSeller(req CreateSellerRequest) (*model.Seller, error)
	GetSeller(id string) (*model.Seller, error)
	GetAllSellers() ([]model.Seller, error)

	// Organizer
	CreateOrganizer(req CreateOrganizerRequest) (*model.Organizer, error)
	GetOrganizer(id uint) (*model.Organizer, error)
	GetAllOrganizers() ([]model.Organizer, error)

	// Category
	CreateCategory(req CreateCategoryRequest) (*model.ItemCategory, error)
	GetCategory(id uint) (*model.ItemCategory, error)
	GetAllCategories() ([]model.ItemCategory, error)

	// Auction Item
	CreateAuctionItem(req CreateAuctionItemRequest) (*model.AuctionItem, error)
	GetAuctionItem(id uint) (*model.AuctionItem, error)
	GetPublishedAuctions(filters repository.AuctionItemFilters) ([]model.AuctionItem, int64, error)
	UpdateAuctionItem(id uint, req UpdateAuctionItemRequest) (*model.AuctionItem, error)
	PublishAuctionItem(id uint) error
	DeleteAuctionItem(id uint) error

	// Bidding
	PlaceBid(req PlaceBidRequest) (*model.Bid, error)
	GetItemBids(itemID uint) ([]model.Bid, error)
	GetUserBids(userID string) ([]model.Bid, error)
}

// ========== REQUEST/RESPONSE STRUCTS ==========

type CreateSellerRequest struct {
	SellerName    string           `json:"seller_name" binding:"required"`
	SellerType    model.SellerType `json:"seller_type"`
	Address       string           `json:"address"`
	Phone         string           `json:"phone"`
	Email         string           `json:"email"`
	ContactPerson string           `json:"contact_person"`
}

type CreateOrganizerRequest struct {
	OrganizerName string              `json:"organizer_name" binding:"required"`
	OrganizerCode string              `json:"organizer_code"`
	OrganizerType model.OrganizerType `json:"organizer_type"`
	Address       string              `json:"address"`
	City          string              `json:"city"`
	Province      string              `json:"province"`
	Phone         string              `json:"phone"`
	Email         string              `json:"email"`
}

type CreateCategoryRequest struct {
	CategoryName     string `json:"category_name" binding:"required"`
	ParentCategoryID *uint  `json:"parent_category_id"`
	Description      string `json:"description"`
}

type CreateAuctionItemRequest struct {
	LotCode             string              `json:"lot_code" binding:"required"`
	ItemName            string              `json:"item_name" binding:"required"`
	CategoryID          uint                `json:"category_id" binding:"required"`
	SellerID            string              `json:"seller_id" binding:"required"`
	OrganizerID         uint                `json:"organizer_id" binding:"required"`
	ItemType            model.ItemType      `json:"item_type" binding:"required"`
	SubType             string              `json:"sub_type"`
	Description         string              `json:"description"`
	DetailedDescription string              `json:"detailed_description"`
	LimitPrice          float64             `json:"limit_price" binding:"required"`
	DepositAmount       float64             `json:"deposit_amount" binding:"required"`
	StartingPrice       float64             `json:"starting_price"`
	IncrementAmount     float64             `json:"increment_amount"`
	AuctionMethod       model.AuctionMethod `json:"auction_method"`
	Images              []ImageRequest      `json:"images"`
	Schedule            *ScheduleRequest    `json:"schedule"`
}

type UpdateAuctionItemRequest struct {
	ItemName            string              `json:"item_name"`
	CategoryID          uint                `json:"category_id"`
	ItemType            model.ItemType      `json:"item_type"`
	SubType             string              `json:"sub_type"`
	Description         string              `json:"description"`
	DetailedDescription string              `json:"detailed_description"`
	LimitPrice          float64             `json:"limit_price"`
	DepositAmount       float64             `json:"deposit_amount"`
	StartingPrice       float64             `json:"starting_price"`
	IncrementAmount     float64             `json:"increment_amount"`
	AuctionMethod       model.AuctionMethod `json:"auction_method"`
	Status              model.AuctionStatus `json:"status"`
	Images              []ImageRequest      `json:"images"`
	Schedule            *ScheduleRequest    `json:"schedule"`
}

type ImageRequest struct {
	ImageURL     string          `json:"image_url" binding:"required"`
	ImageType    model.ImageType `json:"image_type"`
	DisplayOrder int             `json:"display_order"`
	Caption      string          `json:"caption"`
}

type ScheduleRequest struct {
	RegistrationStart time.Time `json:"registration_start"`
	RegistrationEnd   time.Time `json:"registration_end"`
	DepositDeadline   time.Time `json:"deposit_deadline" binding:"required"`
	AuctionStart      time.Time `json:"auction_start" binding:"required"`
	AuctionEnd        time.Time `json:"auction_end" binding:"required"`
	AnnouncementDate  time.Time `json:"announcement_date"`
}

type PlaceBidRequest struct {
	ItemID    uint    `json:"item_id" binding:"required"`
	UserID    string  `json:"user_id" binding:"required"`
	BidAmount float64 `json:"bid_amount" binding:"required"`
	IPAddress string  `json:"ip_address"`
	UserAgent string  `json:"user_agent"`
}

// ========== SERVICE IMPLEMENTATION ==========

type auctionService struct {
	sellerRepo    repository.SellerRepository
	organizerRepo repository.OrganizerRepository
	categoryRepo  repository.CategoryRepository
	itemRepo      repository.AuctionItemRepository
	imageRepo     repository.ItemImageRepository
	scheduleRepo  repository.AuctionScheduleRepository
	bidRepo       repository.BidRepository
	userRepo      repository.UserRepository
}

func NewAuctionService(
	sellerRepo repository.SellerRepository,
	organizerRepo repository.OrganizerRepository,
	categoryRepo repository.CategoryRepository,
	itemRepo repository.AuctionItemRepository,
	imageRepo repository.ItemImageRepository,
	scheduleRepo repository.AuctionScheduleRepository,
	bidRepo repository.BidRepository,
	userRepo repository.UserRepository,
) AuctionService {
	return &auctionService{
		sellerRepo:    sellerRepo,
		organizerRepo: organizerRepo,
		categoryRepo:  categoryRepo,
		itemRepo:      itemRepo,
		imageRepo:     imageRepo,
		scheduleRepo:  scheduleRepo,
		bidRepo:       bidRepo,
		userRepo:      userRepo,
	}
}

// ========== SELLER ==========

func (s *auctionService) CreateSeller(req CreateSellerRequest) (*model.Seller, error) {
	seller := &model.Seller{
		SellerName:    req.SellerName,
		SellerType:    req.SellerType,
		Address:       stringPtr(req.Address),
		Phone:         stringPtr(req.Phone),
		Email:         stringPtr(req.Email),
		ContactPerson: stringPtr(req.ContactPerson),
	}

	if err := s.sellerRepo.Create(seller); err != nil {
		return nil, err
	}

	return seller, nil
}

func (s *auctionService) GetSeller(id string) (*model.Seller, error) {
	return s.sellerRepo.FindByID(id)
}

func (s *auctionService) GetAllSellers() ([]model.Seller, error) {
	return s.sellerRepo.FindAll()
}

// ========== ORGANIZER ==========

func (s *auctionService) CreateOrganizer(req CreateOrganizerRequest) (*model.Organizer, error) {
	organizer := &model.Organizer{
		OrganizerName: req.OrganizerName,
		OrganizerCode: stringPtr(req.OrganizerCode),
		OrganizerType: req.OrganizerType,
		Address:       stringPtr(req.Address),
		City:          stringPtr(req.City),
		Province:      stringPtr(req.Province),
		Phone:         stringPtr(req.Phone),
		Email:         stringPtr(req.Email),
	}

	if err := s.organizerRepo.Create(organizer); err != nil {
		return nil, err
	}

	return organizer, nil
}

func (s *auctionService) GetOrganizer(id uint) (*model.Organizer, error) {
	return s.organizerRepo.FindByID(id)
}

func (s *auctionService) GetAllOrganizers() ([]model.Organizer, error) {
	return s.organizerRepo.FindAll()
}

// ========== CATEGORY ==========

func (s *auctionService) CreateCategory(req CreateCategoryRequest) (*model.ItemCategory, error) {
	category := &model.ItemCategory{
		CategoryName:     req.CategoryName,
		ParentCategoryID: req.ParentCategoryID,
		Description:      stringPtr(req.Description),
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *auctionService) GetCategory(id uint) (*model.ItemCategory, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *auctionService) GetAllCategories() ([]model.ItemCategory, error) {
	return s.categoryRepo.FindAll()
}

// ========== AUCTION ITEM ==========

func (s *auctionService) CreateAuctionItem(req CreateAuctionItemRequest) (*model.AuctionItem, error) {
	// Verify category exists
	if _, err := s.categoryRepo.FindByID(req.CategoryID); err != nil {
		return nil, errors.New("category not found")
	}

	// Verify seller exists
	if _, err := s.sellerRepo.FindByID(req.SellerID); err != nil {
		return nil, errors.New("seller not found")
	}

	// Verify organizer exists
	if _, err := s.organizerRepo.FindByID(req.OrganizerID); err != nil {
		return nil, errors.New("organizer not found")
	}

	item := &model.AuctionItem{
		LotCode:             req.LotCode,
		ItemName:            req.ItemName,
		CategoryID:          req.CategoryID,
		SellerID:            req.SellerID,
		OrganizerID:         req.OrganizerID,
		ItemType:            req.ItemType,
		SubType:             stringPtr(req.SubType),
		Description:         stringPtr(req.Description),
		DetailedDescription: stringPtr(req.DetailedDescription),
		LimitPrice:          decimal.NewFromFloat(req.LimitPrice),
		DepositAmount:       decimal.NewFromFloat(req.DepositAmount),
		StartingPrice:       decimal.NewFromFloat(req.StartingPrice),
		CurrentHighestBid:   decimal.NewFromFloat(req.StartingPrice),
		IncrementAmount:     decimal.NewFromFloat(req.IncrementAmount),
		AuctionMethod:       req.AuctionMethod,
		Status:              model.AuctionStatusDraft,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}

	// Create images if provided
	if len(req.Images) > 0 {
		for _, img := range req.Images {
			image := &model.ItemImage{
				ItemID:       item.ID,
				ImageURL:     img.ImageURL,
				ImageType:    img.ImageType,
				DisplayOrder: img.DisplayOrder,
				Caption:      stringPtr(img.Caption),
			}
			if err := s.imageRepo.Create(image); err != nil {
				return nil, err
			}
		}
	}

	// Create schedule if provided
	if req.Schedule != nil {
		schedule := &model.AuctionSchedule{
			ItemID:            item.ID,
			RegistrationStart: &req.Schedule.RegistrationStart,
			RegistrationEnd:   &req.Schedule.RegistrationEnd,
			DepositDeadline:   req.Schedule.DepositDeadline,
			AuctionStart:      req.Schedule.AuctionStart,
			AuctionEnd:        req.Schedule.AuctionEnd,
			AnnouncementDate:  &req.Schedule.AnnouncementDate,
		}
		if err := s.scheduleRepo.Create(schedule); err != nil {
			return nil, err
		}
	}

	// Fetch complete item with relations
	return s.itemRepo.FindByID(item.ID)
}

func (s *auctionService) GetAuctionItem(id uint) (*model.AuctionItem, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Increment view count
	_ = s.itemRepo.IncrementViewCount(id)

	return item, nil
}

func (s *auctionService) GetPublishedAuctions(filters repository.AuctionItemFilters) ([]model.AuctionItem, int64, error) {
	return s.itemRepo.FindPublished(filters)
}

func (s *auctionService) UpdateAuctionItem(id uint, req UpdateAuctionItemRequest) (*model.AuctionItem, error) {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("auction item not found")
	}

	// Only allow updates on draft items
	if item.Status != model.AuctionStatusDraft {
		return nil, errors.New("cannot update item that is not in draft status")
	}

	// Update fields
	if req.ItemName != "" {
		item.ItemName = req.ItemName
	}
	if req.CategoryID != 0 {
		item.CategoryID = req.CategoryID
	}
	if req.ItemType != "" {
		item.ItemType = req.ItemType
	}
	if req.SubType != "" {
		item.SubType = stringPtr(req.SubType)
	}
	if req.Description != "" {
		item.Description = stringPtr(req.Description)
	}
	if req.DetailedDescription != "" {
		item.DetailedDescription = stringPtr(req.DetailedDescription)
	}
	if req.LimitPrice > 0 {
		item.LimitPrice = decimal.NewFromFloat(req.LimitPrice)
	}
	if req.DepositAmount > 0 {
		item.DepositAmount = decimal.NewFromFloat(req.DepositAmount)
	}
	if req.StartingPrice > 0 {
		item.StartingPrice = decimal.NewFromFloat(req.StartingPrice)
	}
	if req.IncrementAmount > 0 {
		item.IncrementAmount = decimal.NewFromFloat(req.IncrementAmount)
	}
	if req.AuctionMethod != "" {
		item.AuctionMethod = req.AuctionMethod
	}

	if err := s.itemRepo.Update(item); err != nil {
		return nil, err
	}

	// Update images if provided
	if len(req.Images) > 0 {
		// Delete old images
		_ = s.imageRepo.DeleteByItemID(id)

		// Create new images
		for _, img := range req.Images {
			image := &model.ItemImage{
				ItemID:       id,
				ImageURL:     img.ImageURL,
				ImageType:    img.ImageType,
				DisplayOrder: img.DisplayOrder,
				Caption:      stringPtr(img.Caption),
			}
			if err := s.imageRepo.Create(image); err != nil {
				return nil, err
			}
		}
	}

	// Update schedule if provided
	if req.Schedule != nil {
		schedule, err := s.scheduleRepo.FindByItemID(id)
		if err != nil {
			// Create new schedule
			schedule = &model.AuctionSchedule{
				ItemID:            id,
				RegistrationStart: &req.Schedule.RegistrationStart,
				RegistrationEnd:   &req.Schedule.RegistrationEnd,
				DepositDeadline:   req.Schedule.DepositDeadline,
				AuctionStart:      req.Schedule.AuctionStart,
				AuctionEnd:        req.Schedule.AuctionEnd,
				AnnouncementDate:  &req.Schedule.AnnouncementDate,
			}
			if err := s.scheduleRepo.Create(schedule); err != nil {
				return nil, err
			}
		} else {
			// Update existing schedule
			schedule.RegistrationStart = &req.Schedule.RegistrationStart
			schedule.RegistrationEnd = &req.Schedule.RegistrationEnd
			schedule.DepositDeadline = req.Schedule.DepositDeadline
			schedule.AuctionStart = req.Schedule.AuctionStart
			schedule.AuctionEnd = req.Schedule.AuctionEnd
			schedule.AnnouncementDate = &req.Schedule.AnnouncementDate
			if err := s.scheduleRepo.Update(schedule); err != nil {
				return nil, err
			}
		}
	}

	return s.itemRepo.FindByID(id)
}

func (s *auctionService) PublishAuctionItem(id uint) error {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return errors.New("auction item not found")
	}

	if item.Status != model.AuctionStatusDraft {
		return errors.New("can only publish draft items")
	}

	// Check if schedule exists
	if item.Schedule == nil {
		return errors.New("auction schedule is required before publishing")
	}

	return s.itemRepo.UpdateStatus(id, model.AuctionStatusPublished)
}

func (s *auctionService) DeleteAuctionItem(id uint) error {
	item, err := s.itemRepo.FindByID(id)
	if err != nil {
		return errors.New("auction item not found")
	}

	// Only allow deletion of draft items
	if item.Status != model.AuctionStatusDraft {
		return errors.New("can only delete draft items")
	}

	// Delete related data
	_ = s.imageRepo.DeleteByItemID(id)

	return s.itemRepo.Delete(id)
}

// ========== BIDDING ==========

func (s *auctionService) PlaceBid(req PlaceBidRequest) (*model.Bid, error) {
	// Get item
	item, err := s.itemRepo.FindByID(req.ItemID)
	if err != nil {
		return nil, errors.New("auction item not found")
	}

	// Check if auction is ongoing
	if item.Status != model.AuctionStatusOngoing && item.Status != model.AuctionStatusPublished {
		return nil, errors.New("auction is not active")
	}

	// Check auction schedule
	if item.Schedule != nil {
		now := time.Now()
		if now.Before(item.Schedule.AuctionStart) {
			return nil, errors.New("auction has not started yet")
		}
		if now.After(item.Schedule.AuctionEnd) {
			return nil, errors.New("auction has ended")
		}
	}

	// Check minimum bid
	bidAmount := decimal.NewFromFloat(req.BidAmount)
	minBid := item.CurrentHighestBid.Add(item.IncrementAmount)

	if item.BidCount == 0 {
		// First bid must be at least starting price
		if bidAmount.LessThan(item.StartingPrice) {
			return nil, fmt.Errorf("bid must be at least the starting price: %s", item.StartingPrice.String())
		}
	} else {
		if bidAmount.LessThan(minBid) {
			return nil, fmt.Errorf("bid must be at least %s", minBid.String())
		}
	}

	// Get user
	user, err := s.userRepo.FindByID(req.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check user balance
	if user.Balance.LessThan(bidAmount) {
		return nil, errors.New("insufficient balance")
	}

	// Create bid
	bid := &model.Bid{
		ItemID:    req.ItemID,
		UserID:    req.UserID,
		BidAmount: bidAmount,
		BidType:   model.BidTypeManual,
		BidStatus: model.BidStatusWinning,
		IsHighest: true,
		IPAddress: stringPtr(req.IPAddress),
		UserAgent: stringPtr(req.UserAgent),
	}

	if err := s.bidRepo.Create(bid); err != nil {
		return nil, err
	}

	// Mark previous bids as outbid
	_ = s.bidRepo.MarkAllAsOutbid(req.ItemID, bid.ID)

	// Update item with new highest bid
	bidAmountFloat, _ := bidAmount.Float64()
	_ = s.itemRepo.UpdateBidInfo(req.ItemID, bidAmountFloat, item.BidCount+1)

	// Update item status to ongoing if it was published
	if item.Status == model.AuctionStatusPublished {
		_ = s.itemRepo.UpdateStatus(req.ItemID, model.AuctionStatusOngoing)
	}

	return bid, nil
}

func (s *auctionService) GetItemBids(itemID uint) ([]model.Bid, error) {
	return s.bidRepo.FindByItemID(itemID)
}

func (s *auctionService) GetUserBids(userID string) ([]model.Bid, error) {
	return s.bidRepo.FindByUserID(userID)
}

// ========== HELPER FUNCTIONS ==========

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
