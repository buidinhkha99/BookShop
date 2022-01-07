package model

import "gorm.io/gorm"

type Images struct {
	gorm.Model
	Image string `gorm:"type:varchar(1000);" json:"image,omitempty"`
	BookID uint64 `json:"book_id,omitempty" `
}
