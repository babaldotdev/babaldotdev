package main

import (
	"amantuladhar/amantuladhar/pkg/html_to_image"
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log/slog"
	"os"
	"strings"
	"text/template"
)

func main() {
	htmlTemplate, err := os.ReadFile("./cmd/animated_text/animated-text.tmpl")
	if err != nil {
		slog.Error("Unable to read template file")
		panic(1)
	}

	tpl, err := template.New("htmlTemplate").Funcs(template.FuncMap{}).Parse(string(htmlTemplate))
	if err != nil {
		slog.Error("Failed to parse template: %v", err)
		panic(1)
	}

	nFrame := 75
	var buf bytes.Buffer
	err = tpl.Execute(&buf, struct {
		FrameDetails map[int]string
		Text         string
	}{
		FrameDetails: getFrameDetails(nFrame),
		Text:         strings.ReplaceAll("Hey, This is Aman", " ", "&nbsp;"),
	})

	if err != nil {
		slog.Error("Failed to execute template", "err", err)
		panic(1)
	}

	// DEBUG
	//	if err := os.WriteFile("test.html", buf.Bytes(), 0o644); err != nil {
	//		slog.Error("Error writing generated HTML to file", "err", err)
	//		panic(1)
	//	}
	// END DEBUG

	framesSelector := make([]string, nFrame)
	for i := 0; i < nFrame; i++ {
		framesSelector[i] = fmt.Sprintf("#frame-%d", i)
	}
	imgDataList := html_to_image.ImageByteMulti(string(buf.Bytes()), framesSelector)

	// DEBUG
	//	for i, bb := range imgDataList {
	//		if err := os.WriteFile(fmt.Sprintf("frame-%d.png", i), bb, 0o644); err != nil {
	//			slog.Error("Error occurred", "err", err)
	//			panic(1)
	//		}
	//	}
	// END DEBUG

	// GENERATE GIF
	outGif := &gif.GIF{}
	for _, imgData := range imgDataList {
		// Decode PNG image
		inImg, err := png.Decode(bytes.NewReader(imgData))
		if err != nil {
			slog.Error("Unable to decode image from bytes %v \n", err)
			panic(1)
		}

		// Convert to Paletted image
		bounds := inImg.Bounds()
		palettedImg := image.NewPaletted(bounds, palette.Plan9) // Use Plan9 palette (or another palette as needed)
		draw.Draw(palettedImg, bounds, inImg, bounds.Min, draw.Src)

		outGif.Image = append(outGif.Image, palettedImg)
		outGif.Delay = append(outGif.Delay, 0)
	}

	f, _ := os.OpenFile("out.gif", os.O_WRONLY|os.O_CREATE, 0o644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			slog.Error("Error unable to close the file")
			panic(1)
		}
	}(f)

	err = gif.EncodeAll(f, outGif)
	if err != nil {
		slog.Error("Unable to encode GIF", "err", err)
		panic(1)
	}
}

func getFrameDetails(n int) map[int]string {
	x := make(map[int]string)
	for i := 0; i < n; i++ {
		percentage := float32(i+1) / float32(n) * 100.0

		if percentage <= 40 {
			percentage = percentage / 40.0 * 100
		} else if percentage >= 70 {
			percentage = (100 - percentage) / 30.0 * 100
		} else {
			percentage = 100.0
		}
		x[i] = fmt.Sprintf("%0.1f", percentage)
	}
	return x
}
