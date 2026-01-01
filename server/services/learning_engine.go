package services

import (
	"database/sql"
	"vocab-quest/server/core/srs"
	"vocab-quest/server/core/types"
	"vocab-quest/server/repositories"
)

type LearningEngine struct {
	Repo *repositories.ProgressRepository
}

func NewLearningEngine(db *sql.DB) *LearningEngine {
	return &LearningEngine{
		Repo: &repositories.ProgressRepository{DB: db},
	}
}

// GenerateSession creates a balanced session of Review + New words
func (s *LearningEngine) GenerateSession(userID int, dictionaryID int, totalSize int) ([]types.Word, error) {
	// Strategy: Prioritize Reviews. Fill rest with New.

	// 1. Fetch Reviews (Due Now)
	reviews, err := s.Repo.GetDue(userID, dictionaryID, totalSize)
	if err != nil {
		return nil, err
	}

	remainingSlots := totalSize - len(reviews)

	// 2. Fetch New Words
	var newWords []types.Word
	if remainingSlots > 0 {
		newWords, err = s.Repo.GetNew(userID, dictionaryID, remainingSlots)
		if err != nil {
			return nil, err
		}
	}

	// 3. Fallback: Review Ahead (Cramming)
	remainingSlots = totalSize - (len(reviews) + len(newWords))
	var aheadWords []types.Word
	if remainingSlots > 0 {
		aheadWords, err = s.Repo.GetReviewAhead(userID, dictionaryID, remainingSlots)
		if err != nil {
			return nil, err
		}
	}

	// 4. Combine
	session := append(reviews, newWords...)
	session = append(session, aheadWords...)

	// Note: Shuffle removed to enforce "Arranged/Priority" mode
	// Words will appear in order: Due -> New -> ReviewAhead

	return session, nil
}

// SubmitResult processes a gameplay result through the SRS engine
func (s *LearningEngine) SubmitResult(userID, wordID int, success bool) error {
	// 1. Get Current Progress
	progress, err := s.Repo.Get(userID, wordID)
	if err != nil {
		return err
	}

	// First time seeing it?
	if progress == nil {
		progress = &types.UserProgress{
			UserID: userID,
			WordID: wordID,
			// Defaults will be handled by SRS calculator logic handled as '0' stage
		}
	}

	// 2. Calculate SRS Update
	updatedProgress := srs.CalculateReview(*progress, success)

	// 3. Save
	return s.Repo.Save(updatedProgress)
}
