package handlers

import (
	"ProjectBookShop/app/connect"
	"ProjectBookShop/app/model"
	"ProjectBookShop/app/pkg/cache"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func GetAllProdcut(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var books []model.Books
	result := db.Find(&books)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", result.Error)
	} else {
		detailBooks := make(map[int]model.DetailBook)
		for k, book := range books {
			var detailBook model.DetailBook
			detailBook.Book = book
			db.Where("book_id = ?", book.ID).Find(&detailBook.Images)
			detailBooks[k] = detailBook
		}
		dataBooks, _ := json.Marshal(&detailBooks)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(dataBooks))
	}

}
func GetDetailBook(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var detailBook model.DetailBook
	var book model.Books
	var images []model.Images
	var groupBook []model.GroupBooks
	var category model.Category
	vars := mux.Vars(r)
	idBook, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}
	db.Where("id = ?", idBook).Find(&book)
	db.Where("book_id = ?", idBook).Find(&images)
	db.Where("book_id = ?", idBook).Find(&groupBook)
	for _, idGroupBook := range groupBook {
		db.Where("id = ?", idGroupBook).Find(&category)
		detailBook.Category = append(detailBook.Category, category)
	}
	detailBook.GroupBooks = groupBook
	detailBook.Book = book
	detailBook.Images = images
	dataBook, _ := json.Marshal(&detailBook)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(dataBook))

}
func GetTopBooks(w http.ResponseWriter, r *http.Request) {
	key := viper.GetString("redis.key")
	if data := cache.ServeJQueryWithRemoteCache(w, key); data == "" {
		db := connect.Connect()
		var books []model.Books
		db.Limit(9).Order("rate desc").Find(&books)
		b, _ := json.Marshal(books)
		err := cache.InsertData(key, string(b))
		if err != nil {
			log.Error("Error when add data in remote cache")
		}
		fmt.Fprintln(w, string(b))

	} else {
		fmt.Fprintln(w, data)
	}

}
func CreateBooks(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var deltaiBook []model.DetailBookCreate
	err := json.NewDecoder(r.Body).Decode(&deltaiBook)
	if err != nil {
		fmt.Fprint(w, err)
	}
	if len(deltaiBook) != 0 {
		for _, book := range deltaiBook {
			result := db.Create(&book.BookInf)
			if result.Error != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Create book have error: %v", result.Error)
			}
			for _, idCategory := range book.Category {
				result = db.Create(&model.GroupBooks{
					Model:      gorm.Model{},
					CategoryID: uint64(idCategory),
					BookID:     uint64(book.BookInf.ID),
				})
				if result.Error != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Create book have error: %v", result.Error)
				}
			}
			for _, image := range book.Images {
				result = db.Create(&model.Images{
					Model:  gorm.Model{},
					Image:  image,
					BookID: uint64(book.BookInf.ID),
				})
				if result.Error != nil {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Create Images of book have error: %v", result.Error)
				}
			}
			// create data in Elasticserch
			bookElt := model.BookELT{
				ID:                int(book.BookInf.ID),
				Name:              book.BookInf.Name,
				Supplier:          book.BookInf.Supplier,
				PublishingCompany: book.BookInf.PublishingCompany,
				Quantily:          book.BookInf.Quantily,
				Description:       book.BookInf.Description,
				Age:               book.BookInf.Age,
				Price:             book.BookInf.Price,
				PublishingYear:    book.BookInf.PublishingYear,
				Language:          book.BookInf.Language,
				NumberOfPages:     book.BookInf.NumberOfPages,
				Rate:              book.BookInf.Rate,
				Image:             book.Images[0],
			}
			err := bookElt.AddOne()
			if err != nil {
				log.Error("Error add data in Elasticsearch, ID: ", book.BookInf.ID)
			}
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Create book successfull")

	} else {
		fmt.Fprintf(w, "Wrong data !")
	}

}
func DeteleBook(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var book model.Books
	var image model.Images
	vars := mux.Vars(r)
	idBook, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}
	db.Delete(&book, idBook)
	db.Where("book_id =?", idBook).Delete(&image)
	bookElt := model.BookELT{
		ID: int(idBook),
	}
	// delete in Elasticsearch
	err = bookElt.DeleteDoc()
	if err != nil {
		log.Error("Error add data in Elasticsearch, ID: ", idBook)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Delete successfull !")
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var detailBook model.DetailBook
	vars := mux.Vars(r)
	idBook, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}
	err = json.NewDecoder(r.Body).Decode(&detailBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}

	for _, image := range detailBook.Images {
		var imageBook model.Images
		if len(image.Image) != 0 && image.ID != 0 {
			imageBook.BookID = uint64(idBook)
			imageBook.Image = image.Image
			db.Model(&imageBook).Where("id = ?", image.ID).Updates(imageBook)
		}
		// create new image
		if len(image.Image) != 0 && image.ID == 0 {
			imageBook.BookID = uint64(idBook)
			imageBook.Image = image.Image
			result := db.Create(&imageBook)
			if result.Error != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Create Images of book have error: %v", result.Error)
			}
		}
		// delete older image by len of Image == 0
		if len(image.Image) == 0 && image.ID != 0 {
			db.Where("id =?", image.ID).Delete(&imageBook)
		}
	}
	for _, groupBookCa := range detailBook.GroupBooks {
		var groupBook model.GroupBooks
		if groupBookCa.ID != 0 && groupBookCa.CategoryID != 0 {
			groupBook.BookID = uint64(idBook)
			groupBook.CategoryID = uint64(groupBookCa.CategoryID)
			db.Model(&groupBook).Where("id = ?", groupBookCa.ID).Updates(groupBook)
		}
		// create new group
		if groupBookCa.ID == 0 && groupBookCa.CategoryID != 0 {
			groupBook.BookID = uint64(idBook)
			groupBook.CategoryID = uint64(groupBookCa.CategoryID)
			result := db.Create(&groupBook)
			if result.Error != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Create group book have error: %v", result.Error)
			}
		}
		// delete older group by CategoryID == 0
		if groupBookCa.ID != 0 && groupBookCa.CategoryID == 0 {
			db.Where("id =?", groupBookCa.ID).Delete(&groupBook)
		}
	}
	db.Save(&detailBook.Book)
	// Update data in Elastic
	bookElt := model.BookELT{
		ID:                int(detailBook.Book.ID),
		Name:              detailBook.Book.Name,
		Supplier:          detailBook.Book.Supplier,
		PublishingCompany: detailBook.Book.PublishingCompany,
		Quantily:          detailBook.Book.Quantily,
		Description:       detailBook.Book.Description,
		Age:               detailBook.Book.Age,
		Price:             detailBook.Book.Price,
		PublishingYear:    detailBook.Book.PublishingYear,
		Language:          detailBook.Book.Language,
		NumberOfPages:     detailBook.Book.NumberOfPages,
		Rate:              detailBook.Book.Rate,
		Image:             detailBook.Images[0].Image,
	}
	err = bookElt.AddOne()
	if err != nil {
		log.Error("Error add data in Elasticsearch, ID: ", detailBook.Book.ID)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Update successfull !")
}

func ElasticSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := string(vars["name"])
	err, books := model.EsSearchByName(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}
	b, err := json.Marshal(books)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}
	fmt.Fprint(w, string(b))

}
