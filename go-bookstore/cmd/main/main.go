package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	_ "gorm.io/driver/mysql"
	"github.com/kanishkmehta29/go-bookstore/pkg/routes"
)

func main(){
	r := mux.NewRouter()
	routes.RegisterBookStoreRoutes(r)
	http.Handle("/",r)
	log.Fatal(http.ListenAndServe(":3000",r))
}

