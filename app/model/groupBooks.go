package model

import "gorm.io/gorm"

type GroupBooks struct {
	gorm.Model
	CategoryID uint64 `json:"category_id,omitempty"`
	BookID     uint64 `json:"book_id,omitempty"`
}
