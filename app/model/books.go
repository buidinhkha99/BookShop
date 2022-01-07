package model

import (
	"gorm.io/gorm"
)

type Books struct {
	gorm.Model
	Name              string  `json:"name,omitempty" `
	Supplier          string  `json:"supplier,omitempty" `
	PublishingCompany string  `json:"publishing_company,omitempty" `
	Quantily          int     `json:"quantily,omitempty" `
	Description       string  `gorm:"type:varchar(250);" json:"description,omitempty" `
	Age               uint    `json:"age,omitempty" `
	Price             float32 `json:"price,omitempty" `
	PublishingYear    int     `json:"publishing_year,omitempty" `
	Language          string  `json:"language,omitempty" `
	NumberOfPages     int     `json:"number_of_pages,omitempty" `
	Rate              float32 `json:"rate,omitempty" `
}
type DetailBookCreate struct {
	BookInf  Books    `json:"book"`
	Category []int    `json:"category"`
	Images   []string `json:"images"`
}
type DetailBook struct {
	Book       Books        `json:"book"`
	Category   []Category   `json:"category"`
	Images     []Images     `json:"images"`
	GroupBooks []GroupBooks `json:"group"`
}
