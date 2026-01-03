package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ========== ENUMS ==========

type SellerType string

const (
	SellerTypeBank       SellerType = "bank"
	SellerTypeGovernment SellerType = "government"
	SellerTypeCompany    SellerType = "company"
	SellerTypeIndividual SellerType = "individual"
)

type OrganizerType string

const (
	OrganizerTypeKPKNL   OrganizerType = "KPKNL"
	OrganizerTypeBank    OrganizerType = "bank"
	OrganizerTypePrivate OrganizerType = "private"
)

type ItemType string

const (
	ItemTypeMovable   ItemType = "movable"
	ItemTypeImmovable ItemType = "immovable"
)

type AuctionMethod string

const (
	AuctionMethodOpenBidding   AuctionMethod = "open_bidding"
	AuctionMethodClosedBidding AuctionMethod = "closed_bidding"
	AuctionMethodTender        AuctionMethod = "tender"
)

type AuctionStatus string

const (
	AuctionStatusDraft     AuctionStatus = "draft"
	AuctionStatusPublished AuctionStatus = "published"
	AuctionStatusOngoing   AuctionStatus = "ongoing"
	AuctionStatusClosed    AuctionStatus = "closed"
	AuctionStatusCancelled AuctionStatus = "cancelled"
)

type ImageType string

const (
	ImageTypeMain     ImageType = "main"
	ImageTypeGallery  ImageType = "gallery"
	ImageTypeDocument ImageType = "document"
)

type BidType string

const (
	BidTypeManual BidType = "manual"
	BidTypeAuto   BidType = "auto"
	BidTypeProxy  BidType = "proxy"
)

type BidStatus string

const (
	BidStatusActive    BidStatus = "active"
	BidStatusOutbid    BidStatus = "outbid"
	BidStatusWinning   BidStatus = "winning"
	BidStatusWon       BidStatus = "won"
	BidStatusLost      BidStatus = "lost"
	BidStatusCancelled BidStatus = "cancelled"
)

// ========== MODELS ==========

// Seller represents the seller of auction items
type Seller struct {
	ID            string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SellerName    string         `gorm:"type:varchar(255);not null" json:"seller_name"`
	SellerType    SellerType     `gorm:"type:varchar(20)" json:"seller_type"`
	Address       *string        `gorm:"type:text" json:"address,omitempty"`
	Phone         *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email         *string        `gorm:"type:varchar(255)" json:"email,omitempty"`
	ContactPerson *string        `gorm:"type:varchar(255)" json:"contact_person,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID for Seller
func (s *Seller) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (Seller) TableName() string {
	return "sellers"
}

// Organizer represents the auction organizer
type Organizer struct {
	ID            uint           `gorm:"primaryKey;column:organizer_id" json:"id"`
	OrganizerName string         `gorm:"type:varchar(255);not null" json:"organizer_name"`
	OrganizerCode *string        `gorm:"type:varchar(50);uniqueIndex" json:"organizer_code,omitempty"`
	OrganizerType OrganizerType  `gorm:"type:varchar(20)" json:"organizer_type"`
	Address       *string        `gorm:"type:text" json:"address,omitempty"`
	City          *string        `gorm:"type:varchar(100)" json:"city,omitempty"`
	Province      *string        `gorm:"type:varchar(100)" json:"province,omitempty"`
	Phone         *string        `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Email         *string        `gorm:"type:varchar(255)" json:"email,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Organizer) TableName() string {
	return "organizers"
}

// ItemCategory represents auction item categories
type ItemCategory struct {
	ID               uint           `gorm:"primaryKey;column:category_id" json:"id"`
	CategoryName     string         `gorm:"type:varchar(100);not null" json:"category_name"`
	ParentCategoryID *uint          `gorm:"index" json:"parent_category_id,omitempty"`
	Description      *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	ParentCategory *ItemCategory  `gorm:"foreignKey:ParentCategoryID" json:"parent_category,omitempty"`
	SubCategories  []ItemCategory `gorm:"foreignKey:ParentCategoryID" json:"sub_categories,omitempty"`
}

func (ItemCategory) TableName() string {
	return "item_categories"
}

// AuctionItem represents an item being auctioned
type AuctionItem struct {
	ID                  uint            `gorm:"primaryKey;column:item_id" json:"id"`
	LotCode             string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"lot_code"`
	ItemName            string          `gorm:"type:varchar(255);not null" json:"item_name"`
	CategoryID          uint            `gorm:"not null;index" json:"category_id"`
	SellerID            string          `gorm:"type:uuid;not null;index" json:"seller_id"`
	OrganizerID         uint            `gorm:"not null;index" json:"organizer_id"`
	ItemType            ItemType        `gorm:"type:varchar(20);not null" json:"item_type"`
	SubType             *string         `gorm:"type:varchar(100)" json:"sub_type,omitempty"`
	Description         *string         `gorm:"type:text" json:"description,omitempty"`
	DetailedDescription *string         `gorm:"type:text" json:"detailed_description,omitempty"`
	OwnershipProof      *string         `gorm:"type:varchar(255)" json:"ownership_proof,omitempty"`
	OwnershipNumber     *string         `gorm:"type:varchar(255)" json:"ownership_number,omitempty"`
	OwnershipDate       *time.Time      `gorm:"type:date" json:"ownership_date,omitempty"`
	OwnershipHolderName *string         `gorm:"type:varchar(255)" json:"ownership_holder_name,omitempty"`
	LimitPrice          decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"limit_price"`
	DepositAmount       decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"deposit_amount"`
	StartingPrice       decimal.Decimal `gorm:"type:decimal(15,2)" json:"starting_price"`
	CurrentHighestBid   decimal.Decimal `gorm:"type:decimal(15,2)" json:"current_highest_bid"`
	IncrementAmount     decimal.Decimal `gorm:"type:decimal(15,2)" json:"increment_amount"`
	AuctionMethod       AuctionMethod   `gorm:"type:varchar(20)" json:"auction_method"`
	Status              AuctionStatus   `gorm:"type:varchar(20);default:'draft';index" json:"status"`
	ViewCount           int             `gorm:"default:0" json:"view_count"`
	BidCount            int             `gorm:"default:0" json:"bid_count"`
	CreatedAt           time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relations
	Category  *ItemCategory    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Seller    *Seller          `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Organizer *Organizer       `gorm:"foreignKey:OrganizerID" json:"organizer,omitempty"`
	Images    []ItemImage      `gorm:"foreignKey:ItemID" json:"images,omitempty"`
	Schedule  *AuctionSchedule `gorm:"foreignKey:ItemID" json:"schedule,omitempty"`
	Bids      []Bid            `gorm:"foreignKey:ItemID" json:"bids,omitempty"`
}

func (AuctionItem) TableName() string {
	return "auction_items"
}

// ItemImage represents images associated with an auction item
type ItemImage struct {
	ID           uint           `gorm:"primaryKey;column:image_id" json:"id"`
	ItemID       uint           `gorm:"not null;index" json:"item_id"`
	ImageURL     string         `gorm:"type:varchar(500);not null" json:"image_url"`
	ImageType    ImageType      `gorm:"type:varchar(20)" json:"image_type"`
	DisplayOrder int            `gorm:"default:0" json:"display_order"`
	Caption      *string        `gorm:"type:text" json:"caption,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ItemImage) TableName() string {
	return "item_images"
}

// AuctionSchedule represents the schedule for an auction item
type AuctionSchedule struct {
	ID                uint           `gorm:"primaryKey;column:schedule_id" json:"id"`
	ItemID            uint           `gorm:"not null;uniqueIndex" json:"item_id"`
	RegistrationStart *time.Time     `gorm:"type:timestamp" json:"registration_start,omitempty"`
	RegistrationEnd   *time.Time     `gorm:"type:timestamp" json:"registration_end,omitempty"`
	DepositDeadline   time.Time      `gorm:"type:timestamp;not null" json:"deposit_deadline"`
	AuctionStart      time.Time      `gorm:"type:timestamp;not null;index" json:"auction_start"`
	AuctionEnd        time.Time      `gorm:"type:timestamp;not null;index" json:"auction_end"`
	AnnouncementDate  *time.Time     `gorm:"type:timestamp" json:"announcement_date,omitempty"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AuctionSchedule) TableName() string {
	return "auction_schedules"
}

// Bid represents a bid on an auction item
type Bid struct {
	ID        uint            `gorm:"primaryKey;column:bid_id" json:"id"`
	ItemID    uint            `gorm:"not null;index" json:"item_id"`
	UserID    string          `gorm:"type:uuid;not null;index" json:"user_id"`
	BidAmount decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"bid_amount"`
	BidType   BidType         `gorm:"type:varchar(20)" json:"bid_type"`
	BidStatus BidStatus       `gorm:"type:varchar(20);default:'active';index" json:"bid_status"`
	IsHighest bool            `gorm:"default:false" json:"is_highest"`
	BidTime   time.Time       `gorm:"autoCreateTime;index" json:"bid_time"`
	IPAddress *string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent *string         `gorm:"type:text" json:"user_agent,omitempty"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	// Relations
	Item *AuctionItem `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	User *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Bid) TableName() string {
	return "bids"
}
