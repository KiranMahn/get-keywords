package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Function to write file data to a JSON file
func WriteJSONFile(files []File, filename string) error {
	// Open the JSON file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serialize the file data to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ") // Indent JSON output for readability
	if err := encoder.Encode(files); err != nil {
		return err
	}

	// Get the absolute path of the created file
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	// Print the absolute path of the created file
	fmt.Println("File created at:", absPath)

	// fmt.Println(file.Stat())

	return nil
}

func createJSONdata() {

	for _, RepoURL := range RepoURLs {
		// Clone the repository locally
		if err := CloneRepository(RepoURL, CloneDir); err != nil {
			fmt.Println("Failed to clone repository:", err)
			return
		}
		// Walk through the directory and filter files
		filterFunc := func(path string, info os.FileInfo) bool { // TODO:: dont include files under a certain length
			return !info.IsDir() && strings.HasSuffix(info.Name(), ".md") && info.ModTime().After(startDate) && info.ModTime().Before(endDate)
		}
		localDir := "./repository" + RepoURL[28:] // ex ./repository/design-technology/design-docs

		// gets an array of all the filepaths in that repo
		filePaths, err := WalkAndFilterDirectory(localDir, filterFunc)
		if err != nil {
			fmt.Println("Failed to walk through directory:", err)
			return
		}

		// Read file content, extract keywords, and append to slice
		for _, path := range filePaths {

			file, err := getFileDetails(path)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			if file.ContentLength > 200 {
				Files = append(Files, *file)
			}
		}
	}

	tfi, idf = CreateTermFrequencyIndex(Files)

	tfidf := getBetterKeywords(tfi, idf)

	for i := range Files {
		path := Files[i].Path
		wordScores, exists := tfidf[path]
		if !exists {
			continue
		}

		var wordList []WordScore
		for word, score := range wordScores {
			if len(word) < 40 {
				wordList = append(wordList, WordScore{Word: word, Score: score})

			}
		}

		// Sort words by score in descending order
		sort.Slice(wordList, func(i, j int) bool {
			return wordList[i].Score > wordList[j].Score
		})

		// Select top 5 words
		var topKeywords []string
		for j := 0; j < len(wordList) && j < numKeywords; j++ {
			topKeywords = append(topKeywords, wordList[j].Word)
		}

		// Set the keywords for the file
		Files[i].Keywords = topKeywords
		fmt.Println("topKeywords for path: ", path, ": ", topKeywords)
	}

	if err := WriteJSONFile(Files, DataFile); err != nil {
		fmt.Println("Failed to write data to JSON file:", err)
		return
	}

	// end of handle server //////////////////////////////////////
	fmt.Println("Data has been successfully written to", DataFile)

	// handle server //////////////////////////////////////

	// if err := store(); err != nil {
	// 	fmt.Println("Failed to store data to s3 bucket:", err)
	// }
}
