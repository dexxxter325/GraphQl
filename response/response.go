package response

import (
	"encoding/json"
	"net/http"
)

func ErrHandler(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  err.Error(),
		"status": statusCode,
	})
}
