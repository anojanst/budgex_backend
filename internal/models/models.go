package models

import "time"

type Base struct {
	ID        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	UserID    string     `gorm:"index;not null" json:"user_id"`
}

type Category struct {
	Base
	Name     string  `gorm:"not null" json:"name"`
	ParentID *string `gorm:"index" json:"parent_id,omitempty"`
}

type Transaction struct {
	Base
	Type       string    `gorm:"type:text;check:type IN ('income','expense');not null" json:"type"`
	Date       time.Time `gorm:"index" json:"date"`
	Amount     float64   `gorm:"not null" json:"amount"`
	Payee      *string   `json:"payee,omitempty"`
	Memo       *string   `json:"memo,omitempty"`
	CategoryID *string   `gorm:"index" json:"category_id,omitempty"`
	Source     string    `gorm:"default:'manual'" json:"source"`
	Tags       *string   `json:"tags,omitempty"`
}

type Budget struct {
	Base
	Month      string  `gorm:"type:char(7);index" json:"month"`
	CategoryID string  `gorm:"index" json:"category_id"`
	Amount     float64 `gorm:"not null" json:"amount"`
}
