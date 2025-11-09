package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to write JSON", http.StatusInternalServerError)
	}
}

func readJson(body io.ReadCloser, v interface{}) error {
	defer body.Close()
	decoder := json.NewDecoder(body)
	return decoder.Decode(v)
}
