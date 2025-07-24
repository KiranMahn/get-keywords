package main

var (
	RepoURLs        = [3]string{"https://github.com/KiranMahn/Kavi-s-meme-SoundBoard", "https://github.com/KiranMahn/rustpad", "https://github.com/KiranMahn/journal"} // put the git repos you want to extract keywords from here
	customStopwords = [4]string{"org", "company", "inc", "llc"}                                                                                                        // add specific stopwords here like org names to prevent those being keywords (optional)
	numKeywords     = 5                                                                                                                                                // specify the number of keywords here (optional)
)

func main() {
	createJSONdata()
}
