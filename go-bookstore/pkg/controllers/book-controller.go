package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/kanishkmehta29/go-bookstore/pkg/models"
	"github.com/kanishkmehta29/go-bookstore/pkg/utils"
)

var NewBook models.Book

func GetBook(w http.ResponseWriter,req *http.Request){
	newBooks := models.GetAllBooks()
	res,_ := json.Marshal(newBooks)
	w.Header().Set("Content-Type","pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter,req *http.Request){
	vars := mux.Vars(req)
	bid := vars["bookId"]
	Id,err := strconv.ParseInt(bid,0,0)
	if err != nil{
		log.Fatal(err)
	}
	bookdetails,_ := models.GetBookById(Id)
	res,_ := json.Marshal(bookdetails)
	w.Header().Set("Content-Type","pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter,req *http.Request){
	CreateBook := &models.Book{}
	utils.ParseBody(req,CreateBook)
	b := CreateBook.CreateBook()
	res,_ := json.Marshal(b)
	w.Header().Set("Content-Type","pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteBook(w http.ResponseWriter,req *http.Request){
	vars := mux.Vars(req)
	bid := vars["bookId"]
	id,err := strconv.ParseInt(bid,0,0)
	if err != nil {
		log.Fatal(err)
	}
	models.DeleteBook(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted the book"))
}

func UpdateBook(w http.ResponseWriter,req *http.Request){
	vars := mux.Vars(req)
	bid := vars["bookId"]
	newbook := &models.Book{}
	utils.ParseBody(req,newbook)
	bid_int,err := strconv.ParseInt(bid,0,0)
	if err != nil{
		log.Fatal(err)
	}
	bookdetails,db := models.GetBookById(bid_int)
	if newbook.Author != ""{
		bookdetails.Author = newbook.Author
	}
	if newbook.Name != ""{
		bookdetails.Name = newbook.Name
	}
	if newbook.Publication != ""{
		bookdetails.Publication = newbook.Publication
	}
	db.Save(&bookdetails)
	res,_ := json.Marshal(bookdetails)

	w.Header().Set("Content-Type","pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}