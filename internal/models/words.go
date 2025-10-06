package models

type LanguageWords struct {
	Language Language `json:"language"`
	Words    []string `json:"words"`
}
