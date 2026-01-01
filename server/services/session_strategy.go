package services

import (
	"database/sql"
	"math/rand"
	"vocab-quest/server/models"
)

// GetSessionWords implements "Review Fixed, Remedial First" Strategy
// 1. Review (Fixed 3)
// 2. Weak (Fill Remainder) - Sacrifices New if many Weak
// 3. New (Fill Remainder)
func GetSessionWords(db *sql.DB, userID int, limit int) ([]models.Word, error) {
	sessionWords := make([]models.Word, 0)

	calcProf := "CASE WHEN attempts = 0 THEN 0 WHEN (CAST(successes AS FLOAT)/attempts) >= 0.9 THEN 5 WHEN (CAST(successes AS FLOAT)/attempts) >= 0.8 THEN 4 WHEN (CAST(successes AS FLOAT)/attempts) < 0.6 THEN 2 ELSE 3 END"

	// 1. Get REVIEW Words (Priority #1 - Fixed 3)
	// We want exactly 3 if possible to maintain "Review" slot
	reviewLimit := 3
	if limit < 3 {
		reviewLimit = limit
	}

	queryReview := `
		SELECT w.id, w.dictionary_id, w.text, w.definition, w.difficulty, COALESCE(w.pronunciation_url, ''),
		` + calcProf + ` as proficiency
		FROM words w
		JOIN user_progress up ON w.id = up.word_id
		WHERE up.user_id = ? AND up.attempts > 0
		ORDER BY up.last_played_at ASC LIMIT ? 
	`
	rowsReview, err := db.Query(queryReview, userID, reviewLimit)
	if err == nil {
		defer rowsReview.Close()
		for rowsReview.Next() {
			var w models.Word
			rowsReview.Scan(&w.ID, &w.DictionaryID, &w.Text, &w.Definition, &w.Difficulty, &w.PronunciationURL, &w.Proficiency)
			sessionWords = append(sessionWords, w)
		}
	}

	// 2. Get WEAK Words (Priority #2 - Fill All Remaining Slots)
	// If plenty of weak words, this will eat up the rest of the session (sacrificing New)
	slotsRemaining := limit - len(sessionWords)
	if slotsRemaining > 0 {
		queryWeak := `
			SELECT w.id, w.dictionary_id, w.text, w.definition, w.difficulty, COALESCE(w.pronunciation_url, ''), 
			` + calcProf + ` as proficiency
			FROM words w
			JOIN user_progress up ON w.id = up.word_id
			WHERE up.user_id = ? AND up.attempts > 0 AND (CAST(up.successes AS FLOAT) / up.attempts) < 0.6
			ORDER BY RANDOM() LIMIT ?
		`
		rowsWeak, err := db.Query(queryWeak, userID, slotsRemaining)
		if err == nil {
			defer rowsWeak.Close()
			for rowsWeak.Next() {
				var w models.Word
				rowsWeak.Scan(&w.ID, &w.DictionaryID, &w.Text, &w.Definition, &w.Difficulty, &w.PronunciationURL, &w.Proficiency)

				// Dedup with Review words
				exists := false
				for _, ew := range sessionWords {
					if ew.ID == w.ID {
						exists = true
						break
					}
				}
				if !exists {
					sessionWords = append(sessionWords, w)
				}
			}
		}
	}

	// 3. Get NEW Words (Priority #3 - Buffer)
	// Takes whatever is left
	slotsRemaining = limit - len(sessionWords)
	if slotsRemaining > 0 {
		queryNew := `
			SELECT w.id, w.dictionary_id, w.text, w.definition, w.difficulty, COALESCE(w.pronunciation_url, ''), 0 as proficiency
			FROM words w
			LEFT JOIN user_progress up ON w.id = up.word_id AND up.user_id = ?
			WHERE (up.word_id IS NULL OR up.attempts = 0)
			ORDER BY RANDOM() LIMIT ?
		`
		rowsNew, err := db.Query(queryNew, userID, slotsRemaining)
		if err == nil {
			defer rowsNew.Close()
			for rowsNew.Next() {
				var w models.Word
				rowsNew.Scan(&w.ID, &w.DictionaryID, &w.Text, &w.Definition, &w.Difficulty, &w.PronunciationURL, &w.Proficiency)

				exists := false
				for _, ew := range sessionWords {
					if ew.ID == w.ID {
						exists = true
						break
					}
				}
				if !exists {
					sessionWords = append(sessionWords, w)
				}
			}
		}
	}

	// Backfill validation check?
	// The new words query above handles filling empty slots, but if there are NO new words and NO weak words (e.g. perfect master),
	// we might fall short if we don't have enough review words.
	// For now, let's assume database is large enough or just return what we have.

	// Shuffle final result
	rand.Shuffle(len(sessionWords), func(i, j int) {
		sessionWords[i], sessionWords[j] = sessionWords[j], sessionWords[i]
	})

	return sessionWords, nil
}
