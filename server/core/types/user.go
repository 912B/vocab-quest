package types

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}

type UserStats struct {
	UserID        int `json:"user_id"`
	TotalWords    int `json:"total_words"`
	MasteredWords int `json:"mastered_words"`
	LearningWords int `json:"learning_words"`
	NewWords      int `json:"new_words"`
	TotalReviews  int `json:"total_reviews"`
}
