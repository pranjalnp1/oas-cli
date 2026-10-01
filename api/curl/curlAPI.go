// Package handler implements POST /api/curl.
package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yourusername/oas-cli/internal/parser"
	"github.com/yourusername/oas-cli/internal/service"
)

// Handler expects the spec in the request body, with "method" and "path"
// as required query parameters and an optional "baseUrl" override for specs
// that declare a relative server URL, e.g.
// POST /api/curl?method=POST&path=/pets&baseUrl=https://example.com
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	method := query.Get("method")
	path := query.Get("path")
	baseURL := query.Get("baseUrl")
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

	result, err := service.CurlSpec(spec, method, path, baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
