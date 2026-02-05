package health

import (
	"net/http"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
