package handlers

import (
	"ProjectBookShop/app/connect"
	"ProjectBookShop/app/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCategory(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var category []model.Category
	result := db.Find(&category)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", result.Error)
	} else {
		dataCategory, _ := json.Marshal(&category)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(dataCategory))
	}
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var catergory model.Category
	vars := mux.Vars(r)
	idCategory, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}
	catergory.ID = uint(idCategory)
	err = json.NewDecoder(r.Body).Decode(&catergory)
	if err != nil {
		fmt.Fprint(w, err)
	}
	if len(catergory.Name) != 0 {
		db.Model(&catergory).Update("name", catergory.Name)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Update successfull !")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error")
	}

}
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var category model.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	}
	if len(category.Name) != 0 {
		result := db.Create(&category)
		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Create Images of book have error: %v", result.Error)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Create catergory successfull")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Wrong data !")
	}

}
func DeteleCategory(w http.ResponseWriter, r *http.Request) {
	db := connect.Connect()
	var category model.Category
	vars := mux.Vars(r)
	idCategory, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: %v", err)
	}
	db.Delete(&category, idCategory)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Delete successfull !")
}
