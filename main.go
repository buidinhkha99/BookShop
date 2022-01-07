package main

import (
	// "ProjectBookShop/app/migrate"
	// "ProjectBookShop/app/model"
	"ProjectBookShop/app/pkg"
	"ProjectBookShop/app/router"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Start..")
	err := pkg.LoadConfig()
	if err != nil {
		log.Error("Error loading cofig ")

	}

	// auto create table in database
	// migrate.CreateTable()

	// insert data in DB to Elastich search
	// err = model.InsertDataELT()
	// if err != nil {
	// 	log.Error("Error update data to Elasticsearch ")
	// }

	router.Run()

}
