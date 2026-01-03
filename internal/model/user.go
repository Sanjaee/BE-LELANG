package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBlocked   UserStatus = "blocked"
)

type IDCardType string

const (
	IDCardTypeKTP      IDCardType = "KTP"
	IDCardTypeSIM      IDCardType = "SIM"
	IDCardTypePassport IDCardType = "PASSPORT"
)

type User struct {
	ID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Username     *string    `gorm:"type:varchar(100);uniqueIndex" json:"username,omitempty"`
	Phone        *string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	FullName     string     `gorm:"type:varchar(255);not null" json:"full_name"`
	PasswordHash string     `gorm:"type:varchar(255)" json:"-"`
	UserType     string     `gorm:"type:varchar(50);default:'member'" json:"user_type"`
	ProfilePhoto *string    `gorm:"type:text" json:"profile_photo,omitempty"`
	DateOfBirth  *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Gender       *string    `gorm:"type:varchar(20)" json:"gender,omitempty"`

	// ID Card Information
	IDCardNumber *string     `gorm:"type:varchar(50);uniqueIndex" json:"id_card_number,omitempty"`
	IDCardType   *IDCardType `gorm:"type:varchar(20)" json:"id_card_type,omitempty"`

	// Address Information
	Address    *string `gorm:"type:text" json:"address,omitempty"`
	City       *string `gorm:"type:varchar(100)" json:"city,omitempty"`
	Province   *string `gorm:"type:varchar(100)" json:"province,omitempty"`
	PostalCode *string `gorm:"type:varchar(10)" json:"postal_code,omitempty"`

	// Balance for auction bidding
	Balance decimal.Decimal `gorm:"type:decimal(15,2);default:0" json:"balance"`

	// Status and Verification
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	IsVerified        bool       `gorm:"default:false" json:"is_verified"`
	Status            UserStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	VerificationToken *string    `gorm:"type:varchar(255)" json:"-"`
	VerificationDate  *time.Time `gorm:"type:timestamp" json:"verification_date,omitempty"`

	// Login Information
	LastLogin *time.Time `gorm:"type:timestamp" json:"last_login,omitempty"`
	LoginType string     `gorm:"type:varchar(50);default:'credential'" json:"login_type"`
	GoogleID  *string    `gorm:"type:varchar(255);uniqueIndex" json:"-"`

	// OTP for verification
	OTPCode      *string    `gorm:"type:varchar(6)" json:"-"`
	OTPExpiresAt *time.Time `gorm:"type:timestamp" json:"-"`

	// Password Reset
	ResetToken     *string    `gorm:"type:text" json:"-"`
	ResetExpiresAt *time.Time `gorm:"type:timestamp" json:"-"`

	// Timestamps
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}
