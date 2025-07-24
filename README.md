# Repo Keyword Extractor

A Go package that extracts keywords from files in a given code repository. For each file, it outputs a JSON entry containing metadata and content-based keywords.

## Features

- Accepts a repository URL (e.g., GitHub, GitLab)
- Supports custom stopwords (optional)
- Returns a specified number of keywords per file (optional; default: 5)
- Outputs metadata for each file, including:
  - Filename
  - List of extracted keywords
  - Creation date
  - Last modified date
  - File path
  - Title (derived from the file name or contents)
  - Content length

## Notes
Currently supports repositories hosted on GitHub. Must be ran locally for access to private repos

Only markdown files are analyzed.

Keywords are extracted using TF-IDF and stopword filtering logic. Future versions may include NLP enhancements.

## Usage
```go

var (
	RepoURLs    = [3]string{
        "https://github.com/KiranMahn/Kavi-s-meme-SoundBoard", 
        "https://github.com/KiranMahn/rustpad", 
        "https://github.com/KiranMahn/journal"
    }                               // List of repositorys for source data
	stopwords   = []string{}        // Optional
	numKeywords = 5                 // Optional (default is 5)
)

func main() {
	createJSONdata()
}
```

## Output
A file_data.json file will be created in the /data directory with the following structure:
```json
[
  {
    "file_name": "example.go",
    "keywords": ["parser", "token", "lexer", "input", "syntax"],
    "created_at": "2023-10-01T14:52:03Z",
    "last_modified": "2025-07-01T10:14:55Z",
    "file_path": "src/example.go",
    "title": "Example",
    "content_length": 1542
  },
  ...
]
```

## License
MIT License