package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

// Inspired from: https://github.com/chromedp/examples/blob/master/screenshot/main.go
func GenerateDataImageUrl(htmlString string) string {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlDataUrl = htmlStringToDataUrl(htmlString)

	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(htmlDataUrl, `body`, &buf)); err != nil {
		log.Fatal(err)
	}
	dataURL := fmt.Sprintf("data:image/jpg;base64,%s", base64.StdEncoding.EncodeToString(buf))
	return dataURL
}

func htmlStringToDataUrl(html string) string {
	base64HTML := base64.StdEncoding.EncodeToString([]byte(html))
	dataURL := "data:text/html;base64," + base64HTML
	return dataURL
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(url, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
