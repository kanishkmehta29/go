package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
    "math/rand"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Src struct {
    Original   string `json:"original"`
    Large2x    string `json:"large2x"`
    Large      string `json:"large"`
    Medium     string `json:"medium"`
    Small      string `json:"small"`
    Portrait   string `json:"portrait"`
    Landscape  string `json:"landscape"`
    Tiny       string `json:"tiny"`
}

type Photo struct {
    ID             int    `json:"id"`
    Width          int    `json:"width"`
    Height         int    `json:"height"`
    URL            string `json:"url"`
    Photographer   string `json:"photographer"`
    PhotographerURL string `json:"photographer_url"`
    PhotographerID int    `json:"photographer_id"`
    AvgColor       string `json:"avg_color"`
    Src            Src    `json:"src"`
    Liked          bool   `json:"liked"`
    Alt            string `json:"alt"`
}

type QueryData struct {
    TotalResults int     `json:"total_results"`
    Page         int     `json:"page"`
    PerPage      int     `json:"per_page"`
    Photos       []Photo `json:"photos"`
    NextPage     string  `json:"next_page"`
}

type Response struct{
	Link string `json:"link"`
}

func Handlequery(w http.ResponseWriter,r *http.Request){
	vars := mux.Vars(r)
	search_query := vars["query"]

	fmt.Printf("Query came for : %v\n",search_query)

	err := godotenv.Load()
	if err != nil{
		log.Fatalln("Error loading .env file")
	}
	key := os.Getenv("API_KEY")

	url := fmt.Sprintf("https://api.pexels.com/v1/search?query=%s", search_query)

	req, err := http.NewRequest("GET",url,nil)
	if err != nil{
		log.Fatalln(w,"error in making request")
	}

	req.Header = http.Header{
		"Authorization":{key},
	}
    
	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
	log.Fatalln("error in getting pics")
	}
	defer res.Body.Close()

	var d QueryData

	json.NewDecoder(res.Body).Decode(&d)
    
    final_link := Response{
		Link:d.Photos[rand.Intn(15)].Src.Original,
	}
	 
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(final_link)
	
}

func enableCors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            return
        }
        next.ServeHTTP(w, r)
    })
}

func main(){
    
	r := mux.NewRouter()
    
	r.HandleFunc("/",func(w http.ResponseWriter,req *http.Request){
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w,"Welcome to the Pexel-search")
	})

	r.HandleFunc("/search/{query}",Handlequery)



	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080",enableCors(r))

}