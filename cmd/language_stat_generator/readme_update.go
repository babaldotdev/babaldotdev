package main

import (
	"log"
	"os"
	"strings"
)

func updateReadme(html string) {
	// Read the contents of the Readme.md file
	filePath := "Readme.md"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Find the markers in the content
	startMarker := "<!-- Start: LangStat  -->"
	endMarker := "<!-- End: LangStat  -->"
	startIndex := strings.Index(string(content), startMarker)
	endIndex := strings.Index(string(content), endMarker)

	// Update the content between the markers
	if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
		newContent := string(content[:startIndex+len(startMarker)]) + "\n" +
			html +
			string(content[endIndex:])
		err = os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
	} else {
		log.Fatalf("Markers not found or in incorrect order")
	}
}
