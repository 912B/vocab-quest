package repositories

import (
	"database/sql"
	"time"
	"vocab-quest/server/core/types"
)

type ProgressRepository struct {
	DB *sql.DB
}

func (r *ProgressRepository) Get(userID, wordID int) (*types.UserProgress, error) {
	query := `
		SELECT attempts, successes, last_played_at, 
		       next_review_at, 
			   COALESCE(interval, 1), 
			   COALESCE(ease_factor, 2.5), 
			   COALESCE(srs_stage, 0)
		FROM user_progress 
		WHERE user_id = ? AND word_id = ?
	`
	row := r.DB.QueryRow(query, userID, wordID)

	p := &types.UserProgress{UserID: userID, WordID: wordID}

	var nextReview sql.NullTime

	err := row.Scan(
		&p.Attempts,
		&p.Successes,
		&p.LastPlayedAt,
		&nextReview,
		&p.Interval,
		&p.EaseFactor,
		&p.SRSStage,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}

	if nextReview.Valid {
		p.NextReviewAt = nextReview.Time
	} else {
		p.NextReviewAt = p.LastPlayedAt
	}

	return p, nil
}

func (r *ProgressRepository) Save(p types.UserProgress) error {
	query := `
		INSERT INTO user_progress (
			user_id, word_id, attempts, successes, last_played_at,
			next_review_at, interval, ease_factor, srs_stage
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, word_id) DO UPDATE SET
			attempts = excluded.attempts,
			successes = excluded.successes,
			last_played_at = excluded.last_played_at,
			next_review_at = excluded.next_review_at,
			interval = excluded.interval,
			ease_factor = excluded.ease_factor,
			srs_stage = excluded.srs_stage
	`
	_, err := r.DB.Exec(query,
		p.UserID, p.WordID, p.Attempts, p.Successes, p.LastPlayedAt,
		p.NextReviewAt, p.Interval, p.EaseFactor, p.SRSStage,
	)
	return err
}

// GetDue gets words that are due for review (NextReviewAt <= Now)
func (r *ProgressRepository) GetDue(userID int, dictionaryID int, limit int) ([]types.Word, error) {
	query := `
		SELECT w.id, w.dictionary_id, w.text, w.definition, w.difficulty, COALESCE(w.pronunciation_url, ''),
		       up.srs_stage
		FROM words w
		JOIN user_progress up ON w.id = up.word_id
		WHERE up.user_id = ? 
		AND (up.next_review_at <= ? OR up.next_review_at IS NULL)
	`
	args := []interface{}{userID, time.Now()}

	if dictionaryID > 0 {
		query += " AND w.dictionary_id = ?"
		args = append(args, dictionaryID)
	}

	query += " ORDER BY up.next_review_at ASC LIMIT ?"
	args = append(args, limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []types.Word
	for rows.Next() {
		var w types.Word
		rows.Scan(&w.ID, &w.DictionaryID, &w.Text, &w.Definition, &w.Difficulty, &w.PronunciationURL, &w.Proficiency)
		words = append(words, w)
	}
	return words, nil
}

// GetNew gets words that have no progress yet
func (r *ProgressRepository) GetNew(userID int, dictionaryID int, limit int) ([]types.Word, error) {
	query := `
		SELECT w.id, w.dictionary_id, w.text, w.definition, w.difficulty, COALESCE(w.pronunciation_url, ''),
		       0 -- Proficiency
		FROM words w
		LEFT JOIN user_progress up ON w.id = up.word_id AND up.user_id = ?
		WHERE up.word_id IS NULL
	`
	args := []interface{}{userID}

	if dictionaryID > 0 {
		query += " AND w.dictionary_id = ?"
		args = append(args, dictionaryID)
	}

	query += " ORDER BY RANDOM() LIMIT ?"
	args = append(args, limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []types.Word
	for rows.Next() {
		var w types.Word
		rows.Scan(&w.ID, &w.DictionaryID, &w.Text, &w.Definition, &w.Difficulty, &w.PronunciationURL, &w.Proficiency)
		words = append(words, w)
	}
	return words, nil
}

// GetReviewAhead gets words that are scheduled for the future (Cramming)
func (r *ProgressRepository) GetReviewAhead(userID int, dictionaryID int, limit int) ([]types.Word, error) {
	query := `
		SELECT w.id, w.dictionary_id, w.text, w.definition, w.difficulty, COALESCE(w.pronunciation_url, ''),
		       up.srs_stage
		FROM words w
		JOIN user_progress up ON w.id = up.word_id
		WHERE up.user_id = ? 
		AND up.next_review_at > ?
	`
	args := []interface{}{userID, time.Now()}

	if dictionaryID > 0 {
		query += " AND w.dictionary_id = ?"
		args = append(args, dictionaryID)
	}

	query += " ORDER BY up.next_review_at ASC LIMIT ?"
	args = append(args, limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []types.Word
	for rows.Next() {
		var w types.Word
		rows.Scan(&w.ID, &w.DictionaryID, &w.Text, &w.Definition, &w.Difficulty, &w.PronunciationURL, &w.Proficiency)
		words = append(words, w)
	}
	return words, nil
}
