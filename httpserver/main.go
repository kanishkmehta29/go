package main
import (
	"fmt"
	"log"
	"net/http"
)

func landinghandler (w http.ResponseWriter,req *http.Request){
	name := req.FormValue("name")
	place := req.FormValue("place")
	fmt.Fprintf(w,"Hi %v from %v\n",name,place)
}

func hellohandler (w http.ResponseWriter,req *http.Request){
	fmt.Fprintf(w,"Hello world!\n")
}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/",fs)
	http.HandleFunc("/landing",landinghandler)
	http.HandleFunc("/hello",hellohandler)

    fmt.Println("Starting Server on port 8080")
	if err := http.ListenAndServe(":8080",nil);err != nil{
		log.Fatal(err)
	}

}