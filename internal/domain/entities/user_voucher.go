package entities

import (
	"nusantara_service/internal/data/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserVoucherEntity struct {
	ID         uuid.UUID               `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID               `gorm:"type:uuid;index" json:"-"`
	User       model.User              `gorm:"foreignKey:UserID" json:"user"`
	VoucherID  uuid.UUID               `gorm:"type:uuid;index" json:"-"`
	Voucher    model.Voucher           `gorm:"foreignKey:VoucherID" json:"voucher"`
	DetailID   uuid.UUID               `gorm:"type:uuid;index" json:"-"` // Relasi ke snapshot
	Detail     model.UserVoucherDetail `gorm:"foreignKey:DetailID" json:"voucher_detail"`
	IsUsed     bool                    `gorm:"default:false" json:"is_used"`
	RedeemedAt *time.Time              `gorm:"type:timestamp" json:"redeemed_at"`
	ClaimedAt  time.Time               `gorm:"autoCreateTime" json:"claimed_at"`
	CreatedAt  time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt          `gorm:"index" json:"deleted_at"`
}

type UserVoucherDetailEntity struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	VoucherCode       string         `gorm:"type:varchar(50);not null" json:"voucher_code"`
	DiscountType      string         `gorm:"type:varchar(100);not null" json:"discount_type"`
	DiscountAmount    int            `gorm:"not null" json:"discount_amount"`
	DiscountPercent   int            `gorm:"not null" json:"discount_percent"`
	MinPurchaseAmount int            `gorm:"not null" json:"min_purchase_amount"`
	ValidFrom         time.Time      `json:"valid_from"`
	ValidUntil        time.Time      `json:"valid_until"`
	Description       string         `gorm:"type:text" json:"description"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (UserVoucherEntity) TableName() string {
	return "user_vouchers"
}
func (UserVoucherDetailEntity) TableName() string {
	return "user_voucher_details"
}
