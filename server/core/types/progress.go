package types

import "time"

// LearningStage represents the mastery level in SRS
type LearningStage int

const (
	StageNew      LearningStage = 0
	StageLearning LearningStage = 1
	StageReview   LearningStage = 2
	StageMastered LearningStage = 3 // Simplified for UI
)

type UserProgress struct {
	UserID       int       `json:"user_id"`
	WordID       int       `json:"word_id"`
	Attempts     int       `json:"attempts"`
	Successes    int       `json:"successes"`
	LastPlayedAt time.Time `json:"last_played_at"`

	// SRS Fields (The "Brain" Upgrade)
	NextReviewAt time.Time `json:"next_review_at"`
	Interval     float64   `json:"interval"`    // Days
	EaseFactor   float64   `json:"ease_factor"` // Default 2.5
	SRSStage     int       `json:"srs_stage"`   // 0-5 (SM-2 stages)
}
