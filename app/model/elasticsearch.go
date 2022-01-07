package model

import (
	"ProjectBookShop/app/connect"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	elastic "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BookELT struct {
	ID                int     `json:"id,omitempty" `
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
	Image             string  `json:"images,omitempty" `
}
type ESClient struct {
	*elastic.Client
}
type BookManager struct {
	esClient *ESClient
}

func GetESClient() (*elastic.Client, error) {
	url := viper.GetString("elasticsearch.url")
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		log.Error("Error connect elasticsearch: ", err)
	}
	return client, err
}

func EsSearchByName(name string) (error, []*BookELT) {
	esclient, err := GetESClient()
	if err != nil {
		log.Error("Error initializing : ", err)
		return err, nil
	}
	// search
	bm := NewBookManager((*ESClient)(&ESClient{esclient}))
	productGotten := bm.SearchBooks(name)
	return nil, productGotten
}
func NewBookManager(es *ESClient) *BookManager {
	return &BookManager{esClient: es}
}
func (bm *BookManager) SearchBooks(name string) []*BookELT {
	ctx := context.Background()
	index := viper.GetString("elasticsearch.indext")
	if bm.esClient == nil {
		log.Warn("Nil es client")
		return nil
	}
	// build query to search for title
	query := elastic.NewSearchSource()
	query.Query(elastic.NewMatchQuery("name", name))
	// get search's service
	searchService := bm.esClient.
		Search().
		Index(index).
		SearchSource(query)
	// perform search query
	searchResult, err := searchService.Do(ctx)
	if err != nil {
		log.Error("Cannot perform search with ES", err)
		return nil
	}
	// get result
	var books []*BookELT
	for _, hit := range searchResult.Hits.Hits {
		var book BookELT
		err := json.Unmarshal(hit.Source, &book)
		if err != nil {
			log.Error("Get data error: ", err)
			continue
		}
		books = append(books, &book)
	}
	return books
}
func (bookElt BookELT) UpdateField() error {
	ctx := context.Background()
	index := viper.GetString("elasticsearch.indext")
	clientEs, err := GetESClient()
	if err != nil {
		log.Error("Error initializing : ", err)
		return err
	}
	data := map[string]interface{}{
		"id":              bookElt.ID,
		"name":            bookElt.Name,
		"supplier":        bookElt.Supplier,
		"quantily":        bookElt.Quantily,
		"description":     bookElt.Description,
		"age":             bookElt.Age,
		"price":           bookElt.Price,
		"publishing_year": bookElt.PublishingYear,
		"language":        bookElt.Language,
		"number_of_pages": bookElt.NumberOfPages,
		"rate":            bookElt.Rate,
		"images":          bookElt.Image,
	}
	update, err := clientEs.
		Update().
		Index(index).
		Id(strconv.Itoa(int(bookElt.ID))).
		Doc(data).
		Do(ctx)

	if err != nil {
		errC := errors.New("update index " + fmt.Sprint(err))
		return errC
	}
	log.Info("New version of domain ", update.Id)
	return nil
}

func (bookElt BookELT) DeleteDoc() error {
	index := viper.GetString("elasticsearch.indext")
	client, err := GetESClient()
	if err != nil {
		log.Error("Error initializing : ", err)
	}
	res, err := client.Delete().Index(index).Id(strconv.Itoa(int(bookElt.ID))).Do(context.TODO())
	if err != nil {
		return err
	}
	if want, have := "deleted", res.Result; want != have {
		errMsg := fmt.Errorf("expected Result = %q; got %q", want, have)
		return errMsg

	}
	_, err = client.Refresh().Index(index).Do(context.TODO())
	if err != nil {
		return err
	}
	return nil

}

func IsertDatabasetoELT() error {
	err := InsertDataELT()
	if err != nil {
		log.Error("Failed: %v", err)
		return err
	}
	log.Info("Insert data in DB to Elasticsearch successfull")
	return nil
}

func GetAllBooksDB(ch chan BookELT, wg *sync.WaitGroup) {
	db := connect.Connect()
	var books []Books
	result := db.Find(&books)
	if result.Error != nil {
		log.Error("Error when querry data: ", result.Error)
	} else {
		for _, book := range books {
			var image []Images
			db.Where("book_id = ?", book.ID).Find(&image)
			ch <- BookELT{
				ID:                int(book.ID),
				Name:              book.Name,
				Supplier:          book.Supplier,
				PublishingCompany: book.PublishingCompany,
				Quantily:          book.Quantily,
				Description:       book.Description,
				Age:               book.Age,
				Price:             book.Price,
				PublishingYear:    book.PublishingYear,
				Language:          book.Language,
				NumberOfPages:     book.NumberOfPages,
				Rate:              book.Rate,
				Image:             image[0].Image,
			}
		}
		close(ch)
		defer wg.Done()
	}
}

func InsertDataELT() error {
	ch := make(chan BookELT, 20)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go GetAllBooksDB(ch, wg)
	//creating book object
	for bookdata := range ch {
		err := bookdata.AddOne()
		if err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
}
func (bookdata BookELT) AddOne() error {
	esclient, _ := GetESClient()
	index := viper.GetString("elasticsearch.indext")
	ctx := context.Background()
	dataJSON, err := json.Marshal(bookdata)
	if err != nil {
		log.Error("Error marshal data id: ", bookdata.ID)
		return err
	}
	// insert data in ELT
	ind, err := esclient.Index().
		Index(index).
		Id(strconv.Itoa(int(bookdata.ID))).
		BodyJson(string(dataJSON)).
		Do(ctx)
	if err != nil {
		log.Error("Error when insert data in ELT.", string(dataJSON), ind)
		return err
	}
	return nil
}
