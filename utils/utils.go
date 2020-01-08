package utils

import (
	"net/http"
	"encoding/json"
	"go-rest-apis/models"
)

// SendError - notifies client about error
func SendError(w http.ResponseWriter, status int, err models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

// SendSuccess - responds with the data for the request
func SendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
