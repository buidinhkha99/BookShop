package router

import (
	"ProjectBookShop/app/handlers"
	"ProjectBookShop/app/middlewares"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Run() {
	r := mux.NewRouter().StrictSlash(true)
	post := r.Methods(http.MethodPost).Subrouter()
	post.Path("/newbook").HandlerFunc(middlewares.SetMiddleware(handlers.CreateBooks))
	post.Path("/category").HandlerFunc(middlewares.SetMiddleware(handlers.CreateCategory))

	get := r.Methods(http.MethodGet).Subrouter()
	get.Path("/products").HandlerFunc(middlewares.SetMiddleware(handlers.GetAllProdcut))
	get.Path("/category").HandlerFunc(middlewares.SetMiddleware(handlers.GetCategory))
	get.Path("/book/{id}").HandlerFunc(middlewares.SetMiddleware(handlers.GetDetailBook))
	get.Path("/topbooks").HandlerFunc(middlewares.SetMiddleware(handlers.GetTopBooks))
	get.Path("/search-book/{name}").HandlerFunc(middlewares.SetMiddleware(handlers.ElasticSearch))
	delete := r.Methods(http.MethodDelete).Subrouter()
	delete.Path("/book/{id}").HandlerFunc(middlewares.SetMiddleware(handlers.DeteleBook))
	delete.Path("/category/{id}").HandlerFunc(middlewares.SetMiddleware(handlers.DeteleCategory))
	put := r.Methods(http.MethodPut).Subrouter()
	put.Path("/book/{id}").HandlerFunc(middlewares.SetMiddleware(handlers.UpdateBook))
	put.Path("/category/{id}").HandlerFunc(middlewares.SetMiddleware(handlers.UpdateCategory))

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS", "PUT"},
		AllowedHeaders: []string{"*"},
	}).Handler(r)
	http.ListenAndServe(":8080", handler)
}
