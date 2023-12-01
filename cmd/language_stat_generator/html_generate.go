package main

import (
	"bytes"
	"fmt"
	"log"
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
	t, err := template.New("htmlTemplate").Funcs(template.FuncMap{}).Parse(templateString)
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

const templateString = `<div class="lang-stat">
		<div class="lang-diagram">
			{{range $idx, $lang := .Langs}} <div>&nbsp;</div> {{end}}
		</div>
		<div class="lang-detail">
			{{range $idx, $lang := .Langs}} <div class="lang"> <span class="lang-color">&nbsp;</span> <span>{{ $lang }} {{ index $.LanguageStats $lang }}</span> </div> {{end}}
		</div>
	</div>
	<style>
    .lang-diagram {
        display: grid;
        grid-auto-flow: column;
        grid-template-columns: {{ .GridTemplateColumns }} ;
    }
    .lang-stat {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }
    .lang-detail {
        font-size: .8rem;
    }
    .lang-diagram>*:first-child {
        border-top-left-radius: 10px;
        border-bottom-left-radius: 10px;
    }
    .lang-diagram>*:last-child {
        border-top-right-radius: 10px;
        border-bottom-right-radius: 10px;
    }
    .lang-detail {
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
    }
    .lang {
        display: flex;
        gap: .5rem;
        align-items: center;
    }
    .lang-color {
        aspect-ratio: 1/1;
        display: inline-block;
        min-width: 0.8rem;
        max-width: 0.8rem;
        border-radius: 100%;
    }
    .lang-diagram>*:nth-child(1),
    .lang:nth-child(1)>.lang-color {
        background-color: rgb(238, 225, 112);
    }
    .lang-diagram>*:nth-child(2),
    .lang:nth-child(2)>.lang-color {
        background-color: rgb(210, 87, 53);
    }
    .lang-diagram>*:nth-child(3),
    .lang:nth-child(3)>.lang-color {
        background-color: rgb(68, 118, 192);
    }
    .lang-diagram>*:nth-child(4),
    .lang:nth-child(4)>.lang-color {
        background-color: rgb(82, 62, 120);
    }
    .lang-diagram>*:nth-child(5),
    .lang:nth-child(5)>.lang-color {
        background-color: rgb(97, 133, 162);
    }
</style>
`
