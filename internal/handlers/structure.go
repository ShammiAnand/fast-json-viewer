package handlers

import (
	"encoding/json"
	"log"
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
	log.Printf("sessionID: %v", sessionID)
	trie, err := sm.GetSession(sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	children, err := trie.GetChildren("")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(children)
	if err != nil {
		log.Println("error while encoding children", err)
	}
}
