package html_to_image

import (
	"context"
	"encoding/base64"
	"log"
	"os"

	"github.com/chromedp/chromedp"
)

// Inspired from: https://github.com/chromedp/examples/blob/master/screenshot/main.go
func Save(htmlString, elementSelector, fileName string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlDataUrl = htmlStringToDataUrl(htmlString)

	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(htmlDataUrl, elementSelector, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(fileName, buf, 0o644); err != nil {
		// if err := os.WriteFile("lang-stat.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}
}

func imageToDataUrl(image []byte) string {
	base64Image := base64.StdEncoding.EncodeToString(image)
	dataURL := "data:image/png;base64," + base64Image
	return dataURL
}

func htmlStringToDataUrl(html string) string {
	base64HTML := base64.StdEncoding.EncodeToString([]byte(html))
	dataURL := "data:text/html;base64," + base64HTML
	return dataURL
}

func elementScreenshot(url, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
