package services

import (
	"database/sql"
	"time"
	"vocab-quest/server/models"
)

// GetSessionWords logic has been moved to session_strategy.go

func SubmitResult(db *sql.DB, userID int, wordID int, success bool) error {
	// Simple Data Collection
	// 1. Check if exists
	var attempts, successes int
	err := db.QueryRow("SELECT attempts, successes FROM user_progress WHERE word_id = ? AND user_id = ?", wordID, userID).Scan(&attempts, &successes)

	if err == sql.ErrNoRows {
		attempts = 0
		successes = 0
	} else if err != nil {
		return err
	}

	attempts++
	if success {
		successes++
	}

	now := time.Now()

	// 2. Upsert
	query := `
		INSERT INTO user_progress (user_id, word_id, attempts, successes, last_played_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(user_id, word_id) DO UPDATE SET
		attempts = excluded.attempts,
		successes = excluded.successes,
		last_played_at = excluded.last_played_at
	`
	_, err = db.Exec(query, userID, wordID, attempts, successes, now)
	return err
}

func GetUserStats(db *sql.DB, userID int, dictionaryID int) (*models.UserStats, error) {
	// Re-implement stats based on raw data
	stats := &models.UserStats{UserID: userID}

	var err error
	var queryMastered, queryLearning, queryTotal string
	var argsMastered, argsLearning, argsTotal []interface{}

	// Base Args
	argsMastered = []interface{}{userID}
	argsLearning = []interface{}{userID}

	if dictionaryID > 0 {
		// Join with words table to filter by dictionary
		queryMastered = `
			SELECT COUNT(up.word_id) 
			FROM user_progress up 
			JOIN words w ON up.word_id = w.id
			WHERE up.user_id = ? AND w.dictionary_id = ? AND up.attempts >= 3 AND (CAST(up.successes AS FLOAT)/up.attempts) >= 0.8`

		queryLearning = `
			SELECT COUNT(up.word_id) 
			FROM user_progress up 
			JOIN words w ON up.word_id = w.id
			WHERE up.user_id = ? AND w.dictionary_id = ? AND (up.attempts < 3 OR (CAST(up.successes AS FLOAT)/up.attempts) < 0.8)`

		queryTotal = "SELECT COUNT(*) FROM words WHERE dictionary_id = ?"

		argsMastered = append(argsMastered, dictionaryID)
		argsLearning = append(argsLearning, dictionaryID)
		argsTotal = []interface{}{dictionaryID}
	} else {
		// Global Stats (No Join needed, faster)
		queryMastered = `SELECT COUNT(*) FROM user_progress WHERE user_id = ? AND attempts >= 3 AND (CAST(successes AS FLOAT)/attempts) >= 0.8`
		queryLearning = `SELECT COUNT(*) FROM user_progress WHERE user_id = ? AND (attempts < 3 OR (CAST(successes AS FLOAT)/attempts) < 0.8)`
		queryTotal = "SELECT COUNT(*) FROM words"
	}

	// Execute Queries
	if err = db.QueryRow(queryMastered, argsMastered...).Scan(&stats.MasteredWords); err != nil {
		return nil, err
	}
	if err = db.QueryRow(queryLearning, argsLearning...).Scan(&stats.LearningWords); err != nil {
		return nil, err
	}
	if err = db.QueryRow(queryTotal, argsTotal...).Scan(&stats.TotalWords); err != nil {
		return nil, err
	}

	if stats.TotalWords > 0 {
		stats.NewWords = stats.TotalWords - (stats.MasteredWords + stats.LearningWords)
		if stats.NewWords < 0 {
			stats.NewWords = 0
		}
	}

	return stats, nil
}
