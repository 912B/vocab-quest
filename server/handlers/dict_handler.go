package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"vocab-quest/server/models"

	"github.com/xuri/excelize/v2"
)

type DictHandler struct {
	DB *sql.DB
}

func (h *DictHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		query := `
			SELECT d.id, d.name, d.description, d.is_active,
			(SELECT COUNT(*) FROM words WHERE dictionary_id = d.id) as words_count
			FROM dictionaries d
		`
		rows, err := h.DB.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var dicts []models.Dictionary
		for rows.Next() {
			var d models.Dictionary
			if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.IsActive, &d.WordsCount); err != nil {
				continue
			}
			dicts = append(dicts, d)
		}

		if dicts == nil {
			dicts = make([]models.Dictionary, 0)
		}
		json.NewEncoder(w).Encode(dicts)
	} else if r.Method == "POST" {
		var d models.Dictionary
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := h.DB.Exec("INSERT INTO dictionaries (name, description, is_active) VALUES (?, ?, ?)", d.Name, d.Description, d.IsActive)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})

	} else if r.Method == "PUT" {
		var d models.Dictionary
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err := h.DB.Exec("UPDATE dictionaries SET name = ?, description = ?, is_active = ? WHERE id = ?", d.Name, d.Description, d.IsActive, d.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})

	} else if r.Method == "DELETE" {
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)

		// Manual Cascade (optional depending on DB constraints, but safer)
		h.DB.Exec("DELETE FROM words WHERE dictionary_id = ?", id)

		_, err := h.DB.Exec("DELETE FROM dictionaries WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

func (h *DictHandler) SetActive(w http.ResponseWriter, r *http.Request) {
	// Ideally use a router for path params, but for stdlib:
	// Expected path: /api/dictionaries/{id}/active
	// But we simpler API: POST /api/dictionaries/active with Body {id}
	// OR we use the path parsing manually?
	// Let's use simple Logic: POST /api/set-active-dictionary  Body: {id: 1}
	// Or stick to plan: PUT /api/dictionaries/:id/active.
	// Parsing ID from URL in stdlib is annoying without router.
	// Let's implement helper or use query param?
	// POST /api/dictionaries/activate {id: x} is easier.

	// Changing plan slightly for simplicity of stdlib:
	type ActivateReq struct {
		ID int `json:"id"`
	}
	var req ActivateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Deactivate all
	_, err = tx.Exec("UPDATE dictionaries SET is_active = 0")
	if err != nil {
		tx.Rollback()
		return
	}

	// Activate target
	_, err = tx.Exec("UPDATE dictionaries SET is_active = 1 WHERE id = ?", req.ID)
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Get Words of a Dictionary
func (h *DictHandler) ListWords(w http.ResponseWriter, r *http.Request) {
	// GET /api/dictionaries/words?id=1
	idStr := r.URL.Query().Get("id")
	dictID, _ := strconv.Atoi(idStr)

	query := "SELECT id, dictionary_id, text, definition, difficulty FROM words WHERE dictionary_id = ?"
	rows, err := h.DB.Query(query, dictID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var wd models.Word
		if err := rows.Scan(&wd.ID, &wd.DictionaryID, &wd.Text, &wd.Definition, &wd.Difficulty); err != nil {
			continue
		}
		words = append(words, wd)
	}
	if words == nil {
		words = make([]models.Word, 0)
	}
	json.NewEncoder(w).Encode(words)
}

// ImportWords handles Excel (.xlsx) file upload to import words
func (h *DictHandler) ImportWords(w http.ResponseWriter, r *http.Request) {
	// POST /api/dictionaries/import
	// Form: dict_id (int), file (multipart)

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse Multipart Form (10MB max)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// 2. Get Dict ID
	dictIDStr := r.FormValue("dictionary_id")
	dictID, err := strconv.Atoi(dictIDStr)
	if err != nil {
		http.Error(w, "Invalid dictionary ID", http.StatusBadRequest)
		return
	}

	// 3. Get File
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 4. Open Excel File
	f, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Invalid Excel file", http.StatusBadRequest)
		return
	}
	defer f.Close()

	// 5. Read Rows from first sheet
	sheetName := f.GetSheetList()[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		http.Error(w, "Failed to read rows", http.StatusInternalServerError)
		return
	}

	// 6. Insert Words
	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare("INSERT INTO words (dictionary_id, text, definition, difficulty) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Check duplicates (simple check within transaction for now, or INSERT OR IGNORE if supported?
	// SQLite doesn't have INSERT OR IGNORE in standard SQL without conflict clause.
	// We'll check existence manually to be safe/cross-db compatibleish)
	checkStmt, err := tx.Prepare("SELECT 1 FROM words WHERE dictionary_id = ? AND text = ?")
	if err != nil {
		tx.Rollback()
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer checkStmt.Close()

	count := 0
	// Skip header (row 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue // Skip incomplete rows
		}

		wordText := row[0]
		definition := row[1]
		difficulty := 1
		if len(row) > 2 {
			if d, err := strconv.Atoi(row[2]); err == nil {
				difficulty = d
			}
		}

		if wordText == "" {
			continue
		}

		// Check duplicate
		var exists int
		err := checkStmt.QueryRow(dictID, wordText).Scan(&exists)
		if err == sql.ErrNoRows {
			_, err = stmt.Exec(dictID, wordText, definition, difficulty)
			if err == nil {
				count++
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Return Success JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   count,
	})
}

// DownloadTemplate generates and serves an Excel template for importing words
func (h *DictHandler) DownloadTemplate(w http.ResponseWriter, r *http.Request) {
	f := excelize.NewFile()
	defer f.Close()

	// Create Sheet
	index, _ := f.NewSheet("Words")
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // Remove default sheet

	// Set Headers
	headers := []string{"Word", "Definition", "Difficulty (1-10)"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Words", cell, header)
	}

	// Add Sample Row
	f.SetCellValue("Words", "A2", "Example")
	f.SetCellValue("Words", "B2", "This is an example definition")
	f.SetCellValue("Words", "C2", 1)

	// Style Headers (Bold)
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	f.SetCellStyle("Words", "A1", "C1", style)

	// Set Column Widths
	f.SetColWidth("Words", "A", "A", 20)
	f.SetColWidth("Words", "B", "B", 50)
	f.SetColWidth("Words", "C", "C", 15)

	// Write Response
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=vocab_import_template.xlsx")

	if err := f.Write(w); err != nil {
		http.Error(w, "Failed to generate file", http.StatusInternalServerError)
	}
}
