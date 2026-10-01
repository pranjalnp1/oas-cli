// Package handler holds Vercel Go serverless functions. Each file here
// becomes a separate deployed function at /api/<filename-without-ext>.
package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yourusername/oas-cli/internal/parser"
	"github.com/yourusername/oas-cli/internal/service"
)

// Handler implements POST /api/inspect.
//
// TODO: once auth + Supabase storage are wired up, this should accept a
// spec ID (looked up via the authenticated user) instead of taking the
// raw spec body directly on every request.
func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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

	result, err := service.InspectSpec(spec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
