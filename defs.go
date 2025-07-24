package main

import (
	"time"
)

var (
	CloneDir = "repository"
	DataFile = "./data/file_data.json"
	Files    []File
	tfi      TermFrequencyIndex
	idf      []WordData
)

// Date range for filtering files
var (
	startDate = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate   = time.Date(2026, 12, 31, 23, 59, 59, 999, time.UTC)
)

type WordScore struct {
	Word  string
	Score float64
}

type UpdateRequest struct {
	FilePath    string `json:"filePath"`
	ChangeValue int    `json:"changeValue"`
}

type UpdateComment struct {
	FilePath string `json:"filePath"`
	Comment  string `json:"comment"`
}

type UpdateDesignDoc struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
