package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"postgres-golang/models"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct{
	Id int64 `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateStock(w http.ResponseWriter,req *http.Request){
    var stock models.Stock

	err := json.NewDecoder(req.Body).Decode(&stock)

	if err != nil{
		log.Fatalf("%v\n",err)
	}

	newid := models.InsertStock(stock)

	res := Response{
		Id : newid,
		Message : "stock created successfully",
	}
    
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}

func GetStock(w http.ResponseWriter,req *http.Request){
	vars := mux.Vars(req)
	required_id := vars["id"]

	id,err := strconv.Atoi(required_id)

	if err != nil{
		log.Fatalf("%v\n",err)
	}

	stock,err := models.GetStock(int64(id))

	if err != nil{
		log.Fatalf("%v\n",err)
	}

	json.NewEncoder(w).Encode(stock)
	w.WriteHeader(http.StatusOK)

}

func GetAllStock(w http.ResponseWriter,req *http.Request){
	stock,err := models.GetAllStock()

	if err != nil{
		log.Fatalf("%v\n",err)
	}

	json.NewEncoder(w).Encode(stock)
	w.WriteHeader(http.StatusOK)
	
}


func UpdateStock(w http.ResponseWriter,req *http.Request){
	vars := mux.Vars(req)
    
	id,err := strconv.Atoi(vars["id"])
	if err != nil{
		log.Fatalf("%v\n",err)
	}

	var new_stock models.Stock
	err = json.NewDecoder(req.Body).Decode(&new_stock)
	if err != nil{
		log.Fatalf("%v\n",err)
	}

	err = models.UpdateStock(int64(id),new_stock)

    if err != nil{
		log.Fatalf("%v\n",err)
	}
	
	msg := fmt.Sprintf("Stock updated successfully")

	res := Response{
		Id: int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)

}

func DeleteStock(w http.ResponseWriter,req *http.Request){

	vars := mux.Vars(req)

	id,err := strconv.Atoi(vars["id"])

	if err != nil{
		log.Fatalf("%v\n",err)
	}

	models.DeleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted Successfully")

	res := Response{
		Id:int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
	
}