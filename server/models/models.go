package models

import "time"

type Dictionary struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	WordsCount  int    `json:"words_count"`
}

type Word struct {
	ID               int    `json:"id"`
	DictionaryID     int    `json:"dictionary_id"`
	Text             string `json:"text"`
	Definition       string `json:"definition"`
	Difficulty       int    `json:"difficulty"`
	PronunciationURL string `json:"pronunciation_url"`
	Proficiency      int    `json:"proficiency"` // Computed for session
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`      // Don't send password to client
	Role     string `json:"role"`   // admin or user
	Avatar   string `json:"avatar"` // New field
}

type UserProgress struct {
	UserID       int       `json:"user_id"`
	WordID       int       `json:"word_id"`
	Attempts     int       `json:"attempts"`
	Successes    int       `json:"successes"`
	LastPlayedAt time.Time `json:"last_played_at"`
}

type UserStats struct {
	UserID        int `json:"user_id"`
	TotalWords    int `json:"total_words"`
	MasteredWords int `json:"mastered_words"` // Proficiency >= 4
	LearningWords int `json:"learning_words"` // Proficiency 1-3
	NewWords      int `json:"new_words"`      // Proficiency 0
	TotalReviews  int `json:"total_reviews"`  // Total sessions done? Or words reviewed?
	CurrentStreak int `json:"current_streak"` // Potentially calculated
}
