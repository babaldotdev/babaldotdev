package main

import (
	"amantuladhar/amantuladhar/pkg/gif_creator"
	"amantuladhar/amantuladhar/pkg/html_gen"
	"fmt"
	"log/slog"
	"strings"
)

func main() {
	nFrame := 150
	buf, err := html_gen.Generate(
		"./cmd/animated_text/animated-text.tmpl",
		struct {
			FrameDetails map[int]string
			Text         string
		}{
			FrameDetails: getFrameDetails(nFrame),
			Text:         strings.ReplaceAll("(0__o) &#128591; I am Aman. Nice to meet you &#128076;", " ", "&nbsp;"),
		},
	)

	if err != nil {
		slog.Error("Unable to generate HTML", "err", err)
		panic(1)
	}

	err = gif_creator.Create(buf, nFrame, "animated-greeting.gif")
	if err != nil {
		slog.Error("Unable to generate GIF", "err", err)
		panic(1)
	}
}
func getFrameDetails(n int) map[int]string {
	x := make(map[int]string)
	for i := 0; i < n; i++ {
		percentage := float32(i+1) / float32(n) * 100.0

		if percentage <= 35 {
			percentage = percentage / 35 * 100
		} else if percentage >= 80 {
			percentage = (100 - percentage) / 20.0 * 100
		} else {
			percentage = 100.0
		}
		x[i] = fmt.Sprintf("%0.1f", percentage)
	}
	return x
}
