package main

var (
	RepoURLs    = [3]string{"https://github.com/KiranMahn/Kavi-s-meme-SoundBoard", "https://github.com/KiranMahn/rustpad", "https://github.com/KiranMahn/journal"}
	stopwords   = []string{}
	numKeywords = 5
)

func main() {
	createJSONdata()
}
