package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/shammianand/fast-json-viewer/internal/services"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, sm *services.SessionManager, parser *services.Parser) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("jsonFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > parser.MaxFileSize {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	trie, err := parser.ParseJSON(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionID := sm.CreateSession(trie)

	w.Header().Set("Content-Type", "text/html")
	numBytesWritten, err := fmt.Fprintf(w, `<div id="json-viewer" hx-get="/structure/%s" hx-trigger="load"></div>`, sessionID)
	if err != nil {
		log.Println("error while writing to rw", err)
	} else {
		log.Println("wriiten: ", numBytesWritten, " bytes")
	}
}
