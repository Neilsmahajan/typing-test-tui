package models

type LanguageWords struct {
	Language Language `json:"name"`
	Words    []string `json:"words"`
}
