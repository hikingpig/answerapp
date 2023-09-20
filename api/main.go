package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/api/questions/", handleQuestions)
	http.HandleFunc("/api/answers/", handleAnswers)
	http.HandleFunc("/api/votes/", handleVotes)
	appengine.Main()
}
