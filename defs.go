package main

import (
	"os"
	"time"
)

var (
	RepoURLs         = [3]string{"https://github.com/KiranMahn/Kavi-s-meme-SoundBoard", "https://github.com/KiranMahn/rustpad", "https://github.com/KiranMahn/journal"}
	CloneDir         = "repository"
	DataFile         = "./data/file_data.json"
	proposalRepoURL  = "https://github.com/repo1/design-docs"
	proposalCloneDir = "./branches/"

	Files         []File
	proposalFiles []File
	tfi           TermFrequencyIndex
	idf           []WordData

	// Define AWS credentials
	awsAccessKeyID     = os.Getenv("awsAccessKeyID")
	awsSecretAccessKey = os.Getenv("awsSecretAccessKey")
	region             = "US"
	endpoint           = "s3.endpoint"
	bucket             = "imp-articles"
	key                = "file_data.json"

	// define if cloning should be done
	doClone = os.Getenv("doClone")
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
