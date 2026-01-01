package srs

import (
	"math"
	"time"
	"vocab-quest/server/core/types"
)

// CalculateReview updates the progress based on the performance (success/fail).
// This is an implementation of a simplified SM-2 Algorithm.
func CalculateReview(current types.UserProgress, success bool) types.UserProgress {
	updated := current
	updated.LastPlayedAt = time.Now()
	updated.Attempts++

	if success {
		updated.Successes++

		if updated.SRSStage == 0 {
			// First success: Interval = 1 day
			updated.Interval = 1
			updated.SRSStage = 1
		} else if updated.SRSStage == 1 {
			// Second success: Interval = 6 days
			updated.Interval = 6
			updated.SRSStage = 2
		} else {
			// Subsequent: Interval * EaseFactor
			if updated.EaseFactor < 1.3 {
				updated.EaseFactor = 1.3
			}
			updated.Interval = math.Ceil(updated.Interval * updated.EaseFactor)
			updated.SRSStage++
		}

		// Bonus: Increase Ease Factor slightly on success?
		// Standard SM-2 adjusts EF based on strict "quality" (0-5 rating).
		// Since we only have Boolean (Pass/Fail), we keep EF stable or slight bump.
		updated.EaseFactor += 0.1

	} else {
		// Failure: Reset Interval
		updated.SRSStage = 0   // Back to Learning
		updated.Interval = 0.5 // Review in 12 hours (Next Session)
		updated.EaseFactor -= 0.2
		if updated.EaseFactor < 1.3 {
			updated.EaseFactor = 1.3
		}
	}

	// Calculate Next Review Time
	// Interval is in Days
	duration := time.Duration(updated.Interval * 24 * float64(time.Hour))
	updated.NextReviewAt = time.Now().Add(duration)

	return updated
}
