package types

type Dictionary struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type Word struct {
	ID               int    `json:"id"`
	DictionaryID     int    `json:"dictionary_id"`
	Text             string `json:"text"`
	Definition       string `json:"definition"`
	Difficulty       int    `json:"difficulty"`
	PronunciationURL string `json:"pronunciation_url"`

	// Contextual/Computed
	Proficiency int `json:"proficiency,omitempty"`
}
