package repository

import (
	"yourapp/internal/model"

	"gorm.io/gorm"
)

// ========== SELLER REPOSITORY ==========

type SellerRepository interface {
	Create(seller *model.Seller) error
	FindByID(id string) (*model.Seller, error)
	FindAll() ([]model.Seller, error)
	Update(seller *model.Seller) error
	Delete(id string) error
}

type sellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) SellerRepository {
	return &sellerRepository{db: db}
}

func (r *sellerRepository) Create(seller *model.Seller) error {
	return r.db.Create(seller).Error
}

func (r *sellerRepository) FindByID(id string) (*model.Seller, error) {
	var seller model.Seller
	err := r.db.Where("id = ?", id).First(&seller).Error
	return &seller, err
}

func (r *sellerRepository) FindAll() ([]model.Seller, error) {
	var sellers []model.Seller
	err := r.db.Find(&sellers).Error
	return sellers, err
}

func (r *sellerRepository) Update(seller *model.Seller) error {
	return r.db.Save(seller).Error
}

func (r *sellerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Seller{}).Error
}

// ========== ORGANIZER REPOSITORY ==========

type OrganizerRepository interface {
	Create(organizer *model.Organizer) error
	FindByID(id uint) (*model.Organizer, error)
	FindAll() ([]model.Organizer, error)
	Update(organizer *model.Organizer) error
	Delete(id uint) error
}

type organizerRepository struct {
	db *gorm.DB
}

func NewOrganizerRepository(db *gorm.DB) OrganizerRepository {
	return &organizerRepository{db: db}
}

func (r *organizerRepository) Create(organizer *model.Organizer) error {
	return r.db.Create(organizer).Error
}

func (r *organizerRepository) FindByID(id uint) (*model.Organizer, error) {
	var organizer model.Organizer
	err := r.db.First(&organizer, id).Error
	return &organizer, err
}

func (r *organizerRepository) FindAll() ([]model.Organizer, error) {
	var organizers []model.Organizer
	err := r.db.Find(&organizers).Error
	return organizers, err
}

func (r *organizerRepository) Update(organizer *model.Organizer) error {
	return r.db.Save(organizer).Error
}

func (r *organizerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Organizer{}, id).Error
}

// ========== CATEGORY REPOSITORY ==========

type CategoryRepository interface {
	Create(category *model.ItemCategory) error
	FindByID(id uint) (*model.ItemCategory, error)
	FindAll() ([]model.ItemCategory, error)
	FindRootCategories() ([]model.ItemCategory, error)
	Update(category *model.ItemCategory) error
	Delete(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.ItemCategory) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uint) (*model.ItemCategory, error) {
	var category model.ItemCategory
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *categoryRepository) FindAll() ([]model.ItemCategory, error) {
	var categories []model.ItemCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindRootCategories() ([]model.ItemCategory, error) {
	var categories []model.ItemCategory
	err := r.db.Where("parent_category_id IS NULL").Preload("SubCategories").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) Update(category *model.ItemCategory) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.ItemCategory{}, id).Error
}

// ========== AUCTION ITEM REPOSITORY ==========

type AuctionItemRepository interface {
	Create(item *model.AuctionItem) error
	FindByID(id uint) (*model.AuctionItem, error)
	FindByLotCode(lotCode string) (*model.AuctionItem, error)
	FindAll(filters AuctionItemFilters) ([]model.AuctionItem, int64, error)
	FindPublished(filters AuctionItemFilters) ([]model.AuctionItem, int64, error)
	Update(item *model.AuctionItem) error
	UpdateStatus(id uint, status model.AuctionStatus) error
	UpdateBidInfo(id uint, highestBid float64, bidCount int) error
	IncrementViewCount(id uint) error
	Delete(id uint) error
}

type AuctionItemFilters struct {
	CategoryID *uint
	SellerID   *string
	Status     *model.AuctionStatus
	Search     string
	Page       int
	Limit      int
	SortBy     string
	SortOrder  string
}

type auctionItemRepository struct {
	db *gorm.DB
}

func NewAuctionItemRepository(db *gorm.DB) AuctionItemRepository {
	return &auctionItemRepository{db: db}
}

func (r *auctionItemRepository) Create(item *model.AuctionItem) error {
	return r.db.Create(item).Error
}

func (r *auctionItemRepository) FindByID(id uint) (*model.AuctionItem, error) {
	var item model.AuctionItem
	err := r.db.
		Preload("Category").
		Preload("Seller").
		Preload("Organizer").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Preload("Schedule").
		First(&item, id).Error
	return &item, err
}

func (r *auctionItemRepository) FindByLotCode(lotCode string) (*model.AuctionItem, error) {
	var item model.AuctionItem
	err := r.db.
		Preload("Category").
		Preload("Seller").
		Preload("Organizer").
		Preload("Images").
		Preload("Schedule").
		Where("lot_code = ?", lotCode).
		First(&item).Error
	return &item, err
}

func (r *auctionItemRepository) FindAll(filters AuctionItemFilters) ([]model.AuctionItem, int64, error) {
	var items []model.AuctionItem
	var total int64

	query := r.db.Model(&model.AuctionItem{})

	// Apply filters
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.SellerID != nil {
		query = query.Where("seller_id = ?", *filters.SellerID)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("item_name ILIKE ? OR lot_code ILIKE ? OR description ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Count total
	query.Count(&total)

	// Apply sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Page > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset)
	}

	// Preload relations
	err := query.
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("image_type = ?", model.ImageTypeMain).Limit(1)
		}).
		Preload("Schedule").
		Find(&items).Error

	return items, total, err
}

func (r *auctionItemRepository) FindPublished(filters AuctionItemFilters) ([]model.AuctionItem, int64, error) {
	// Only show published and ongoing items
	status := model.AuctionStatusPublished
	filters.Status = &status

	// Also include ongoing items
	var items []model.AuctionItem
	var total int64

	query := r.db.Model(&model.AuctionItem{}).
		Where("status IN ?", []model.AuctionStatus{model.AuctionStatusPublished, model.AuctionStatusOngoing})

	// Apply other filters
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("item_name ILIKE ? OR lot_code ILIKE ? OR description ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	// Sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}
	sortOrder := "DESC"
	if filters.SortOrder != "" {
		sortOrder = filters.SortOrder
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Page > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset)
	}

	err := query.
		Preload("Category").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Preload("Schedule").
		Find(&items).Error

	return items, total, err
}

func (r *auctionItemRepository) Update(item *model.AuctionItem) error {
	return r.db.Save(item).Error
}

func (r *auctionItemRepository) UpdateStatus(id uint, status model.AuctionStatus) error {
	return r.db.Model(&model.AuctionItem{}).Where("item_id = ?", id).Update("status", status).Error
}

func (r *auctionItemRepository) UpdateBidInfo(id uint, highestBid float64, bidCount int) error {
	return r.db.Model(&model.AuctionItem{}).
		Where("item_id = ?", id).
		Updates(map[string]interface{}{
			"current_highest_bid": highestBid,
			"bid_count":           bidCount,
		}).Error
}

func (r *auctionItemRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.AuctionItem{}).
		Where("item_id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *auctionItemRepository) Delete(id uint) error {
	return r.db.Delete(&model.AuctionItem{}, id).Error
}

// ========== IMAGE REPOSITORY ==========

type ItemImageRepository interface {
	Create(image *model.ItemImage) error
	CreateBatch(images []model.ItemImage) error
	FindByItemID(itemID uint) ([]model.ItemImage, error)
	Update(image *model.ItemImage) error
	Delete(id uint) error
	DeleteByItemID(itemID uint) error
}

type itemImageRepository struct {
	db *gorm.DB
}

func NewItemImageRepository(db *gorm.DB) ItemImageRepository {
	return &itemImageRepository{db: db}
}

func (r *itemImageRepository) Create(image *model.ItemImage) error {
	return r.db.Create(image).Error
}

func (r *itemImageRepository) CreateBatch(images []model.ItemImage) error {
	return r.db.Create(&images).Error
}

func (r *itemImageRepository) FindByItemID(itemID uint) ([]model.ItemImage, error) {
	var images []model.ItemImage
	err := r.db.Where("item_id = ?", itemID).Order("display_order ASC").Find(&images).Error
	return images, err
}

func (r *itemImageRepository) Update(image *model.ItemImage) error {
	return r.db.Save(image).Error
}

func (r *itemImageRepository) Delete(id uint) error {
	return r.db.Delete(&model.ItemImage{}, id).Error
}

func (r *itemImageRepository) DeleteByItemID(itemID uint) error {
	return r.db.Where("item_id = ?", itemID).Delete(&model.ItemImage{}).Error
}

// ========== SCHEDULE REPOSITORY ==========

type AuctionScheduleRepository interface {
	Create(schedule *model.AuctionSchedule) error
	FindByItemID(itemID uint) (*model.AuctionSchedule, error)
	Update(schedule *model.AuctionSchedule) error
	Delete(id uint) error
}

type auctionScheduleRepository struct {
	db *gorm.DB
}

func NewAuctionScheduleRepository(db *gorm.DB) AuctionScheduleRepository {
	return &auctionScheduleRepository{db: db}
}

func (r *auctionScheduleRepository) Create(schedule *model.AuctionSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *auctionScheduleRepository) FindByItemID(itemID uint) (*model.AuctionSchedule, error) {
	var schedule model.AuctionSchedule
	err := r.db.Where("item_id = ?", itemID).First(&schedule).Error
	return &schedule, err
}

func (r *auctionScheduleRepository) Update(schedule *model.AuctionSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *auctionScheduleRepository) Delete(id uint) error {
	return r.db.Delete(&model.AuctionSchedule{}, id).Error
}

// ========== BID REPOSITORY ==========

type BidRepository interface {
	Create(bid *model.Bid) error
	FindByID(id uint) (*model.Bid, error)
	FindByItemID(itemID uint) ([]model.Bid, error)
	FindByUserID(userID string) ([]model.Bid, error)
	FindHighestBid(itemID uint) (*model.Bid, error)
	FindByItemAndUser(itemID uint, userID string) ([]model.Bid, error)
	Update(bid *model.Bid) error
	UpdateStatus(id uint, status model.BidStatus) error
	MarkAllAsOutbid(itemID uint, exceptBidID uint) error
}

type bidRepository struct {
	db *gorm.DB
}

func NewBidRepository(db *gorm.DB) BidRepository {
	return &bidRepository{db: db}
}

func (r *bidRepository) Create(bid *model.Bid) error {
	return r.db.Create(bid).Error
}

func (r *bidRepository) FindByID(id uint) (*model.Bid, error) {
	var bid model.Bid
	err := r.db.Preload("User").Preload("Item").First(&bid, id).Error
	return &bid, err
}

func (r *bidRepository) FindByItemID(itemID uint) ([]model.Bid, error) {
	var bids []model.Bid
	err := r.db.Where("item_id = ?", itemID).
		Preload("User").
		Order("bid_amount DESC").
		Find(&bids).Error
	return bids, err
}

func (r *bidRepository) FindByUserID(userID string) ([]model.Bid, error) {
	var bids []model.Bid
	err := r.db.Where("user_id = ?", userID).
		Preload("Item").
		Order("bid_time DESC").
		Find(&bids).Error
	return bids, err
}

func (r *bidRepository) FindHighestBid(itemID uint) (*model.Bid, error) {
	var bid model.Bid
	err := r.db.Where("item_id = ? AND bid_status IN ?", itemID, []model.BidStatus{model.BidStatusActive, model.BidStatusWinning}).
		Order("bid_amount DESC").
		First(&bid).Error
	return &bid, err
}

func (r *bidRepository) FindByItemAndUser(itemID uint, userID string) ([]model.Bid, error) {
	var bids []model.Bid
	err := r.db.Where("item_id = ? AND user_id = ?", itemID, userID).
		Order("bid_time DESC").
		Find(&bids).Error
	return bids, err
}

func (r *bidRepository) Update(bid *model.Bid) error {
	return r.db.Save(bid).Error
}

func (r *bidRepository) UpdateStatus(id uint, status model.BidStatus) error {
	return r.db.Model(&model.Bid{}).Where("bid_id = ?", id).Update("bid_status", status).Error
}

func (r *bidRepository) MarkAllAsOutbid(itemID uint, exceptBidID uint) error {
	return r.db.Model(&model.Bid{}).
		Where("item_id = ? AND bid_id != ? AND bid_status = ?", itemID, exceptBidID, model.BidStatusActive).
		Updates(map[string]interface{}{
			"bid_status": model.BidStatusOutbid,
			"is_highest": false,
		}).Error
}
