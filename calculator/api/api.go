package api

import (
	"calculator/internal/handler"
	"net/http"
)

func HandleApi() {
	http.HandleFunc("/add", handler.HandleAdd)
	http.HandleFunc("/subtract", handler.HandleSubtract)
	http.HandleFunc("/multiply", handler.HandleMultiply)
	http.HandleFunc("/divide", handler.HandleDivide)
	http.HandleFunc("/sum", handler.HandleSum)
}
