package models

type Quote struct {
	Text string `json:"text"`
}

type LanguageQuotes struct {
	Language Language `json:"language"`
	Quotes   []Quote  `json:"quotes"`
}
