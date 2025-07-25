# Repo Keyword Extractor

A Go package that extracts keywords from files in a given Github code repository. For each file in the repos, it outputs a JSON entry containing metadata and content-based keywords.

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

## Hosting
The dockerfile attached with serve the file_data.json file on port 8081.
To build and run the dockerfile enter the following in your command line: 

```bash docker build -t get-keywords . && docker run -p 8081:8081 get-keywords ```

This will serve the code on port 8081. If running locally, you can view the output at http://localhost:8081/ 

See [this repo](https://github.com/KiranMahn/dtr) for an example on how to deploy this as a backend service with a frontend service reading the output and displaying it in an interactive react app. (TODO: will be made public soon)

## License
MIT License