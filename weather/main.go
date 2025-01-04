package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
)

type weatherData struct {
    Name string `json:"name"`
    Main struct {
        Kelvin float64 `json:"temp"`
    } `json:"main"`
}

func loadApiConfig() string {
    apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
    if apiKey == "" {
        log.Fatal("OPENWEATHERMAP_API_KEY environment variable is not set")
    }
    return apiKey
}

func query(city string) (weatherData, error) {
    key := loadApiConfig()
    res, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + key + "&q=" + city)
    if err != nil {
        return weatherData{}, err
    }
    defer res.Body.Close()

    var d weatherData
    err = json.NewDecoder(res.Body).Decode(&d)
    if err != nil {
        return weatherData{}, err
    }
    return d, nil
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "Hello to the program")
    })

    r.HandleFunc("/weather/{city}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        city := vars["city"]
        data, err := query(city)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            log.Printf("error in finding the weather: %v", err)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(data)
    })

    fmt.Println("Starting server on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}