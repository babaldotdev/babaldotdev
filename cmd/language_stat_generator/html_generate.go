package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

func generateHtml(languageStats map[string]int, langs []string) string {
	total := 0
	for _, v := range langs {
		total += languageStats[v]
	}
	languageStatsPercentage := make(map[string]string)
	for _, lang := range langs {
		v := languageStats[lang]
		languageStatsPercentage[lang] = fmt.Sprintf("%0.1f%%", float64(v)/float64(total)*100)
	}

	gridTemplateColumns := ""
	for _, lang := range langs {
		gridTemplateColumns += fmt.Sprintf("%s ", languageStatsPercentage[lang])
	}

	// Create a template from the string
	templateByte, err := os.ReadFile("./cmd/language_stat_generator/lang-stat.tmpl")
	if err != nil {
		log.Fatalf("Unable to read lang-stat.tmpl file")
	}
	htmlTemplate := string(templateByte)
	t, err := template.New("htmlTemplate").Funcs(template.FuncMap{}).Parse(htmlTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Execute the template with the languageStats map
	var buf bytes.Buffer
	err = t.Execute(&buf, struct {
		GridTemplateColumns string
		LanguageStats       map[string]string
		Langs               []string
	}{
		GridTemplateColumns: gridTemplateColumns,
		LanguageStats:       languageStatsPercentage,
		Langs:               langs,
	})
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	return buf.String()
}
