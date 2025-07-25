package main

import (
	"net/http"

	"github.com/rs/cors"
)

var (
	// put the git repos you want to extract keywords from here
	RepoURLs = [3]string{
		"https://github.com/KiranMahn/Kavi-s-meme-SoundBoard",
		"https://github.com/KiranMahn/rustpad",
		"https://github.com/KiranMahn/journal",
	}

	// add specific stopwords here like org names to prevent those being keywords (optional)
	customStopwords = [4]string{"org", "company", "inc", "llc"}

	// specify the maximum number of keywords per file here (optional. Must be greater than 0)
	numKeywords = 5
)

func main() {
	createJSONdata()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHttpRequest)

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8081", handler)
}
