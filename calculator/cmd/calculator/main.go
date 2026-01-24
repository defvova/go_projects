package main

import (
	"calculator/api"
	"fmt"
	"log"
	"net/http"
)

func main() {
	api.HandleApi()
	fmt.Println("Server starting on port 3000 ...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
