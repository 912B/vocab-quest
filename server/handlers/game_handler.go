package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"vocab-quest/server/services"
)

type GameHandler struct {
	DB     *sql.DB
	Engine *services.LearningEngine
}

func (h *GameHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	dictIDStr := r.URL.Query().Get("dictionary_id")
	var dictID int
	if dictIDStr != "" {
		fmt.Sscanf(dictIDStr, "%d", &dictID)
	}

	words, err := h.Engine.GenerateSession(userID, dictID, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(words)
}

type ResultRequest struct {
	WordID  int  `json:"word_id"`
	Success bool `json:"success"`
}

func (h *GameHandler) SubmitResult(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.Engine.SubmitResult(userID, req.WordID, req.Success); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *GameHandler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		userIDStr = r.Header.Get("X-User-ID")
	}

	if userIDStr == "" {
		http.Error(w, "Unauthorized: Missing User ID", http.StatusUnauthorized)
		return
	}

	var userID int
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Parse Dictionary ID
	dictIDStr := r.URL.Query().Get("dictionary_id")
	var dictID int
	if dictIDStr != "" {
		fmt.Sscanf(dictIDStr, "%d", &dictID)
	}

	stats, err := services.GetUserStats(h.DB, userID, dictID)
	if err != nil {
		log.Println("Error fetching stats:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
