package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"vocab-quest/server/models"
)

type WordHandler struct {
	DB *sql.DB
}

func (h *WordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// List Words
		// Optional filter: dictionary_id
		dictIdStr := r.URL.Query().Get("dictionary_id")

		var query string
		var args []interface{}

		if dictIdStr != "" {
			dictID, _ := strconv.Atoi(dictIdStr)
			query = "SELECT id, dictionary_id, text, definition, difficulty FROM words WHERE dictionary_id = ?"
			args = append(args, dictID)
		} else {
			// Default to all or first active? For admin list, maybe all?
			query = "SELECT id, dictionary_id, text, definition, difficulty FROM words"
		}

		rows, err := h.DB.Query(query, args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		words := []models.Word{}
		for rows.Next() {
			var w models.Word
			if err := rows.Scan(&w.ID, &w.DictionaryID, &w.Text, &w.Definition, &w.Difficulty); err != nil {
				continue
			}
			words = append(words, w)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(words)

	case "POST":
		// Create Word
		var wd models.Word
		if err := json.NewDecoder(r.Body).Decode(&wd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := h.DB.Exec("INSERT INTO words (dictionary_id, text, definition, difficulty) VALUES (?, ?, ?, ?)",
			wd.DictionaryID, wd.Text, wd.Definition, wd.Difficulty)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case "PUT":
		// Update Word
		var wd models.Word
		if err := json.NewDecoder(r.Body).Decode(&wd); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := h.DB.Exec("UPDATE words SET text=?, definition=?, difficulty=? WHERE id=?",
			wd.Text, wd.Definition, wd.Difficulty, wd.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case "DELETE":
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		_, err := h.DB.Exec("DELETE FROM words WHERE id=?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
