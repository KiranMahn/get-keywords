package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// File represents information about a file
type File struct {
	Name          string
	Title         string
	Path          string
	LastModified  time.Time
	Created       time.Time // New field for creation time
	Keywords      []string
	ContentLength int
	Content       string // Add Content field to store file content
	IsOrphan      bool
	FileType      string // Add FileType field to store file type
	Org           string // organisation
	Upvotes       int
	Downvotes     int
	Comments      []string
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Function to retrieve file content
func ReadFileContent(file *File) error {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return err
	}
	file.Content = string(content) // Assign file content to Content field
	file.ContentLength = len(content)
	return nil
}

// Function to retrieve keywords from file content
func ExtractKeywords(file *File) {
	// Implement your logic to extract keywords from content
	// This could involve using regular expressions or natural language processing
	// For simplicity, let's assume we're extracting words longer than 3 characters
	words := strings.Fields(file.Content) // Use Content field instead of content variable
	for _, word := range words {
		if len(word) > 3 && len(word) < 15 {
			file.Keywords = append(file.Keywords, word)
		}
	}

	// if there are more than 10 keywords, limit to numKeywords
	if (len(file.Keywords) >= 10) && (numKeywords > 0) && (len(file.Keywords) >= numKeywords) {
		file.Keywords = file.Keywords[0:numKeywords]
	}
}

// Function to clone or update a GitHub repository
func CloneRepository(repoURL, cloneDir string) error {
	print("\ncloning repo: " + repoURL + "\n")
	// Check if the clone directory already exists
	// Directory does not exist, perform git clone
	fmt.Printf("Cloning %s into %s\n", repoURL, cloneDir+repoURL[28:])
	cmd := exec.Command("git", "clone", repoURL, cloneDir+repoURL[28:])
	return cmd.Run()
}

// Helper function to run a command and return its output
func runCommand(cmd *exec.Cmd) (string, error) {
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Function to remove the cloned repository directory
func removeRepository(cloneDir string) error {
	cmd := exec.Command("rm", "-rf", cloneDir)
	return cmd.Run()
}

// Function to walk through a directory and filter files
func WalkAndFilterDirectory(dir string, filterFunc func(path string, info os.FileInfo) bool) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// fmt.Println("WAlking filepath: ", path)
		if err != nil {
			return err
		}
		if !info.IsDir() && filterFunc(path, info) {
			// fmt.Println("appending dir: ", dir, ", path: ", path)
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// Function to determine if a file is an orphan
func isOrphan(file *File) bool {
	// Calculate the duration since last modification
	durationSinceLastModified := time.Since(file.LastModified)

	// Check if the file is older than 3 years and has not been modified in the past 2 years
	if durationSinceLastModified > (3*365*24*time.Hour) && durationSinceLastModified > (2*365*24*time.Hour) {
		return true
	}
	wordCount := len(strings.Fields(file.Content))
	if wordCount < 200 {
		return true
	} else {
		return false
	}
}

// getFileType takes a file path and returns the file type as a string.
func getFileType(path string) (string, error) {
	// First, check the file extension.
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md":
		return "markdown", nil
	case ".txt":
		return "text", nil
	case ".json":
		return "json", nil
	case ".html":
		return "html", nil
	case ".xml":
		return "xml", nil
	case ".go":
		return "go", nil
	case ".py":
		return "python", nil
	case ".java":
		return "java", nil
	case ".js":
		return "javascript", nil
	case ".ts":
		return "typescript", nil
	case ".yaml", ".yml":
		return "yaml", nil
	case ".toml":
		return "toml", nil
	// Add more extensions as needed
	default:
		// If the extension is not recognized, fall back to content-based detection.
	}
	return "unknown", nil
}

// getFileDetails retrieves details of a file, including type detection.
func getFileDetails(path string) (*File, error) {

	files, err := getExisitingFiles()
	if err != nil {
		fmt.Println("error getting existing files :( : %s", err)
	}

	// check if already exists
	fileExists := alreadyExists(files, path)

	if fileExists {
		return updateFile(files, path)
	} else {
		return mkNewFile(path)
	}

}

// findTitle parses a markdown file and returns the title
func findTitle(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return "none"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Remove leading '#' characters and any leading/trailing whitespace
			title := strings.TrimSpace(strings.TrimLeft(line, "#"))
			// println("found title: ", title)
			return title
		}
		if strings.HasPrefix(line, "title:") {
			// Remove leading '#' characters and any leading/trailing whitespace
			title := strings.TrimSpace(strings.TrimLeft(line, "title:"))
			// println("found title: ", title)

			return title
		}
	}
	return "none"
}

func getExisitingFiles() ([]File, error) {

	// get exisitng files
	var files []File
	fileData, err := os.ReadFile("./data/file_data.json")
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	err = json.Unmarshal(fileData, &files)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	//         usage:
	// for i, file := range files {
	// 	if file.Path == req.FilePath {
	// 		files[i].Comments = append(files[i].Comments, req.Comment)
	// 		break
	// 	}
	// }

	return files, nil
}

func alreadyExists(files []File, path string) bool {
	for _, file := range files {
		if file.Path == path {
			return true
		}
	}
	return false
}

func getFile(files []File, path string) *File {
	for _, file := range files {
		if file.Path == path {
			return &file
		}
	}
	fmt.Errorf("File was said to exists but didnt! returning empty file")
	return nil
}

func mkNewFile(path string) (*File, error) {
	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Get last modified time using custom function (getLastModified)
	lastModified, err := getLastModified(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get last modified time: %w", err)
	}

	// Get creation time using git blame
	creationTime, err := getDocumentCreationDate(path)
	if err != nil {
		// Handle error or ignore if creation time is not critical
		fmt.Printf("Warning: Failed to get creation time for %s: %v\n", path, err)
	}

	parent := getParentRepo(path)

	title := findTitle(path)

	// Create File struct
	file := &File{
		Name:         filepath.Base(path),
		Title:        title,
		Path:         path,
		LastModified: lastModified,
		Created:      creationTime,
		IsOrphan:     isOrphan(&File{LastModified: fileInfo.ModTime(), Path: path}),
		Org:          parent,
		Upvotes:      0,
	}

	// Read file content
	if err := ReadFileContent(file); err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Extract keywords
	ExtractKeywords(file)

	// Detect file type
	fileType, err := getFileType(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file type: %w", err)
	}
	file.FileType = fileType

	return file, nil

}

func updateFile(files []File, path string) (*File, error) {
	// Get file info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Get last modified time using custom function (getLastModified)
	lastModified, err := getLastModified(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get last modified time: %w", err)
	}

	title := findTitle(path)

	file := getFile(files, path)

	file.Name = filepath.Base(path)
	file.Title = title
	file.LastModified = lastModified
	file.IsOrphan = isOrphan(&File{LastModified: fileInfo.ModTime(), Path: path})
	// Read file content
	if err := ReadFileContent(file); err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	// Extract keywords
	ExtractKeywords(file)

	// Detect file type
	fileType, err := getFileType(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file type: %w", err)
	}
	file.FileType = fileType

	return file, nil

}

func CloneRepositoryBranches(repoURL, cloneDir string) error {
	repoName := repoURL[strings.LastIndex(repoURL, "/")+1:] // Extract repo name from URL
	repoPath := cloneDir + repoName

	// Check if the clone directory already exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// Directory does not exist, perform git clone
		fmt.Printf("Cloning %s into %s\n", repoURL, repoPath)
		cmd := exec.Command("git", "clone", repoURL, repoPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
	}

	// Change directory to the cloned repository
	if err := os.Chdir(repoPath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	// Fetch all branches
	fmt.Println("Fetching all branches")
	cmd := exec.Command("git", "fetch", "--all")
	if _, err := runCommand(cmd); err != nil {
		return fmt.Errorf("failed to fetch all branches: %w", err)
	}

	// Get the list of branches
	cmd = exec.Command("git", "branch", "-r")
	output, err := runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to list branches: %w", err)
	}

	// Get the list of merged branches
	cmd = exec.Command("git", "branch", "--merged")
	mergedBranchesOutput, err := runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to list merged branches: %w", err)
	}
	mergedBranches := make(map[string]bool)
	for _, branch := range strings.Split(mergedBranchesOutput, "\n") {
		branch = strings.TrimSpace(branch)
		if branch != "" {
			mergedBranches[branch] = true
		}
	}

	// Loop through each branch and clone it if it is not merged
	branches := strings.Split(output, "\n")
	for _, branch := range branches {
		branch = strings.TrimSpace(branch)
		if branch == "" || strings.HasPrefix(branch, "origin/HEAD") {
			continue // Skip invalid or non-branch references
		}
		if strings.HasPrefix(branch, "origin/") {
			branchName := branch[7:] // Remove 'origin/' prefix

			// Check if the branch is not merged (i.e., it's active)
			if _, ok := mergedBranches[branchName]; !ok {
				// Create a directory for the branch
				branchPath := "." + cloneDir + branchName
				if err := os.MkdirAll(branchPath, 0755); err != nil {
					return fmt.Errorf("failed to create directory for branch %s: %w", branchName, err)
				}

				// Clone the branch into the new directory
				fmt.Printf("Cloning branch %s into %s\n", branchName, branchPath)
				cmd = exec.Command("git", "clone", "--branch", branchName, repoURL, branchPath)
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to clone branch %s: %w", branchName, err)
				}
			}
		}
	}

	return nil
}
