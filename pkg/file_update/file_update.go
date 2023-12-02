package file_update

import (
	"log"
	"os"
	"strings"
)

func Update(filePath, contentToAdd, startMarker, endMarker string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	startIndex := strings.Index(string(content), startMarker)
	endIndex := strings.Index(string(content), endMarker)

	// Update the content between the markers
	if startIndex != -1 && endIndex != -1 && startIndex < endIndex {
		newContent := string(content[:startIndex+len(startMarker)]) + "\n\n" +
			contentToAdd + "\n\n" +
			string(content[endIndex:])
		err = os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
	} else {
		log.Fatalf("Markers not found or in incorrect order")
	}
}
