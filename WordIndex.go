package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type WordData struct {
	Word      string  `json:"word"`
	Frequency int     `json:"frequency"`
	IDF       float64 `json:"idf"`
}

var wordDataList []WordData

type TermFrequencyIndex map[string]map[string]int

type TfIdf map[string]map[string]float64

// LoadStopwords reads a file containing stopwords and returns a set of those words
func LoadStopwords(filePath string) (map[string]struct{}, error) {
	stopwords := make(map[string]struct{})
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("err:", err) // Debugging print
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		stopwords[word] = struct{}{}
		// fmt.Println("Loaded stopword:", word) // Debugging print
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return stopwords, nil
}

// returns a number to indicate if a word is a keyword or not
func getProbabilisticidf(word string, files []File, tfi TermFrequencyIndex) float64 {
	N := len(files)
	nt := len(tfi[word])
	if (N == 0) || (nt == 0) {
		return math.Inf(-1)
	}
	percent := float64((N - nt) / nt)
	// print("\n N: total num of docs: ", N)
	// print("\n nt: total num of docs containing word: ", nt)
	// fmt.Printf("\n percent: %f", percent)
	logfreq := math.Log(float64(percent))
	return logfreq
}

// Count the frequency of each word across all files
func getWordFrequencyAcrossAllFiles(files []File, stopwords map[string]struct{}, wordSplitter *regexp.Regexp) map[string]int {
	wordFrequency := make(map[string]int)

	// for each file
	for _, file := range files {

		// get all words in file
		words := wordSplitter.Split(strings.ToLower(file.Content), -1)

		// for each word in file
		for _, word := range words {

			// if word is not a space and not a digit and its length is longer than 2
			if word != "" && !unicode.IsDigit(rune(word[0])) && len(word) > 2 {

				// if the word is not a stopword
				if _, ok := stopwords[word]; !ok {

					// add to total word frquency
					wordFrequency[word]++
				}
			}
		}
	}
	return wordFrequency
}

func getWordFrequencyForEachFile(files []File, stopwords map[string]struct{}, wordSplitter *regexp.Regexp, wordFrequency map[string]int) TermFrequencyIndex {
	// make tfi
	tfi := make(TermFrequencyIndex)

	// for each file
	for _, file := range files {
		// get words
		words := wordSplitter.Split(strings.ToLower(file.Content), -1)
		// for each word
		for _, word := range words {
			// if word is not a space or digit and is longer than two and not a stopword
			if word != "" && !unicode.IsDigit(rune(word[0])) && len(word) > 2 {
				if _, ok := stopwords[word]; !ok {
					// if the words frquency is more than one and less than 50
					if wordFrequency[word] > 1 {

						// add word to tfi and increase count
						if _, found := tfi[word]; !found {
							tfi[word] = make(map[string]int)
						}
						tfi[word][file.Path]++
						getProbabilisticidf(word, files, tfi)
					}
				}
			}
		}
	}

	return tfi

}

// CreateTermFrequencyIndex creates a term frequency index from a list of files
func CreateTermFrequencyIndex(files []File) (TermFrequencyIndex, []WordData) {
	// get stopwords
	stopwords, err := LoadStopwords("./data/stopwords.txt")
	if err != nil {
		print("Error loading stopwords:", err)
	}

	// Regular expression to split words by non-alphanumeric characters
	wordSplitter := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	// Step 1: Count the frequency of each word across all files
	wordFrequency := getWordFrequencyAcrossAllFiles(files, stopwords, wordSplitter)

	// make tfi
	tfi := getWordFrequencyForEachFile(files, stopwords, wordSplitter, wordFrequency)

	for word := range tfi {
		idf := getProbabilisticidf(word, files, tfi)
		freq := len(tfi[word])
		// Create a WordData struct for the current word
		wordData := WordData{
			Word:      word,
			Frequency: freq,
			IDF:       idf,
		}
		// fmt.Printf("\n  {Word: \"%s\", frequency: %d, idf: %f},\n", word, freq, idf)
		wordDataList = append(wordDataList, wordData)

		// Iterate over the inner map
		// for _, frequency := range freqMap {
		// 	fmt.Printf("\n  Word: %s, frequency: %d, idf: %f\n", word, frequency, getProbabilisticidf(word, TestFiles, tfi))
		// }
	}
	// TestCreateTermFrequencyIndex(files)

	return tfi, wordDataList
}

func TestCreateTermFrequencyIndex(testFiles []File) {

	// get stopwords
	stopwords, err := LoadStopwords("./data/stopwords.txt")
	if err != nil {
		print("Error loading stopwords:", err)
	}

	// Regular expression to split words by non-alphanumeric characters
	wordSplitter := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	// Step 1: Count the frequency of each word across all files
	wordFrequency := getWordFrequencyAcrossAllFiles(testFiles, stopwords, wordSplitter)

	// make tfi
	tfi := getWordFrequencyForEachFile(testFiles, stopwords, wordSplitter, wordFrequency)

	for word := range tfi {
		idf := getProbabilisticidf(word, testFiles, tfi)
		freq := len(tfi[word])
		// Create a WordData struct for the current word
		wordData := WordData{
			Word:      word,
			Frequency: freq,
			IDF:       idf,
		}
		fmt.Printf("\n  {Word: \"%s\", frequency: %d, idf: %f},\n", word, freq, idf)
		wordDataList = append(wordDataList, wordData)
	}

	println(len(tfi))
}
func findWordIndex(wordDataList []WordData, targetWord string) int {
	for i, wordData := range wordDataList {
		if wordData.Word == targetWord {
			return i
		}
	}
	return -1 // Return -1 if the word is not found
}

func getBetterKeywords(tfi TermFrequencyIndex, idf []WordData) TfIdf {
	result := make(TfIdf)

	for word, paths := range tfi {
		for path, freq := range paths {
			// Get index of word
			thisIndex := findWordIndex(idf, word)
			if thisIndex == -1 {
				// Word not found in idf slice, skip it
				continue
			}

			// Get idf of that word
			thisIdf := idf[thisIndex].IDF

			// Multiply them
			total := float64(freq) * thisIdf

			// Add that path, word, and tf-idf to the result map
			if result[path] == nil {
				result[path] = make(map[string]float64)
			}
			result[path][word] = total
		}
	}
	// fmt.Println(result)

	return result
}
