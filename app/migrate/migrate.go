package migrate

import (
	"ProjectBookShop/app/connect"
	"ProjectBookShop/app/model"
)

func CreateTable() {
	db := connect.Connect()
	db.AutoMigrate(&model.Books{}, &model.Category{}, &model.GroupBooks{}, &model.Images{})
}
