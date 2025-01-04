package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"math/rand"
	"net/http"
	"github.com/gorilla/mux"
)

type Movie struct{
	Id string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct{
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter,req *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter,req *http.Request){

	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(req)
	for _,item := range movies {

		if item.Id == params["id"]{
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func deleteMovie(w http.ResponseWriter,req *http.Request){

	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(req)
	for index,item := range movies {

		if item.Id == params["id"]{
			movies = append(movies[:index],movies[index+1:]...)
			break
		}
	}
}

func createMovie(w http.ResponseWriter,req *http.Request){

	w.Header().Set("Content-Type","application/json")
	var new_movie Movie
	_ = json.NewDecoder(req.Body).Decode(&new_movie)
	new_movie.Id = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies,new_movie)
	json.NewEncoder(w).Encode(new_movie)

}

func updateMovie(w http.ResponseWriter,req *http.Request){

	w.Header().Set("Content-Type","application/json")
	var new_movie Movie
	_ = json.NewDecoder(req.Body).Decode(&new_movie)
	params := mux.Vars(req)
	cid := params["id"]
	for index,item := range movies{
		if item.Id == cid{
			movies = append(movies[:index],movies[index+1:]...)
			movies = append(movies,new_movie)
			break
		}
	}

}
func main(){

	r := mux.NewRouter()

	movies = append(movies,Movie{Id:"1",Isbn:"438227",Title:"3 idiots",Director:&Director{Firstname: "Aamir",Lastname: "Khan"}})
    movies = append(movies,Movie{Id:"2",Isbn:"438228",Title:"iron man",Director:&Director{Firstname: "Robert",Lastname: "Downey"}})

	r.HandleFunc("/movies",getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}",getMovie).Methods("GET")
	r.HandleFunc("/movies",createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}",updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}",deleteMovie).Methods("DELETE")

	fmt.Println("Running Server on port 3000")
	log.Fatal(http.ListenAndServe(":3000",r))

}