package main

import (
	"fmt"
	"log"
	"net/http"
	"postgres-golang/router"
)

func main(){

  r := router.Router()
  fmt.Println("Starting server on port 8080")
  log.Fatal(http.ListenAndServe(":8080",r))
  
}