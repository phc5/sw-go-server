package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprint(res, "This is home")
}

func main() {
	http.HandleFunc("/", homeHandler)
	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
