package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	// Read file content from the local JSON file
	filePath := "./data/file_data.json" // Adjust the path and filename as needed
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read the local JSON file: ", err)
		http.Error(w, fmt.Sprintf("Failed to read the local JSON file: %s", err), http.StatusInternalServerError)
		return
	}

	// Parse JSON data
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Println("Failed to parse JSON: ", err)
		http.Error(w, fmt.Sprintf("Failed to parse JSON: %s", err), http.StatusInternalServerError)
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Prepare JSON response
	var jsonResponse []byte
	if jsonData != nil {
		jsonResponse, err = json.Marshal(jsonData)
		if err != nil {
			fmt.Println("Failed to marshal JSON: ", err)
			http.Error(w, fmt.Sprintf("Failed to marshal JSON: %s", err), http.StatusInternalServerError)
			return
		}
	} else {
		jsonResponse, err = json.Marshal(map[string]string{"Message": "wompedy womp"})
		if err != nil {
			fmt.Println("Failed to marshal JSON: ", err)
			http.Error(w, fmt.Sprintf("Failed to marshal JSON: %s", err), http.StatusInternalServerError)
			return
		}
	}

	// Set Content-Length header
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(jsonResponse)))

	// Write JSON response
	w.Write(jsonResponse)
}
