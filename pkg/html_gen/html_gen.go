package html_gen

import (
	"bytes"
	"log/slog"
	"os"
	"text/template"
)

func Generate(tmplPath string, tmplData any) (bytes.Buffer, error) {
	var buf bytes.Buffer
	htmlTemplate, err := os.ReadFile(tmplPath)
	if err != nil {
		slog.Error("Unable to read template file")
		return buf, err
	}
	tpl, err := template.New("htmlTemplate").Funcs(template.FuncMap{}).Parse(string(htmlTemplate))
	if err != nil {
		slog.Error("Failed to parse template", err)
		return buf, err
	}
	if err := tpl.Execute(&buf, tmplData); err != nil {
		slog.Error("Failed to execute template", "err", err)
		panic(1)
	}
	return buf, nil
}

func Save(buf bytes.Buffer, fileName string) error {
	if err := os.WriteFile(fileName, buf.Bytes(), 0o644); err != nil {
		slog.Error("Error writing generated HTML to file", "err", err)
		return err
	}
	return nil
}
