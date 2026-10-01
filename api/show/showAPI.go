// Package handler implements POST /api/show.
package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yourusername/oas-cli/internal/parser"
	"github.com/yourusername/oas-cli/internal/service"
)

// Handler expects the spec in the request body, with "method" and "path"
// as query parameters, e.g. POST /api/show?method=GET&path=/pets/{petId}.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	method := r.URL.Query().Get("method")
	path := r.URL.Query().Get("path")
	if method == "" || path == "" {
		http.Error(w, "method and path query parameters are required", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}

	spec, err := parser.Parse(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := service.ShowSpec(spec, method, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
