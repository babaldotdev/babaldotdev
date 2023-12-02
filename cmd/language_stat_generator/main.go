package main

import (
	"amantuladhar/amantuladhar/pkg/github_client"
	"amantuladhar/amantuladhar/pkg/html_to_image"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}
	ghc := github_client.New().FetchAllRepository()
	html := generateHtml(ghc.LangMap, ghc.SortedLang[:5])
	html_to_image.Save(html, "div.lang-stat", "lang-stat.png")
	// No need to update readme as image URL will be same
}
