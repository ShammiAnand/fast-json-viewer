package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shammianand/fast-json-viewer/internal/services"
)

func GetStructureHandler(w http.ResponseWriter, r *http.Request, sm *services.SessionManager) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sessionID := parts[2]
	trie, err := sm.GetSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	children, err := trie.GetChildren("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(children)
}
