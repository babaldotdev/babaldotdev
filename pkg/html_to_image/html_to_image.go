package html_to_image

import (
	"context"
	"encoding/base64"
	"log/slog"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func Save(htmlString, elementSelector, fileName string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var htmlDataUrl = htmlStringToDataUrl(htmlString)

	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(htmlDataUrl, elementSelector, &buf)); err != nil {
		slog.Error("Unable to take screenshot", "selector", elementSelector, "err", err)
		panic(1)
	}
	if err := os.WriteFile(fileName, buf, 0o644); err != nil {
		slog.Error("Unable to save image to a file", "fileName", fileName, "err", err)
		panic(1)
	}
}

func ImageByteMulti(htmlString string, elementSelectors []string) ([][]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var htmlDataUrl = htmlStringToDataUrl(htmlString)

	// Load the HTML first
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(htmlDataUrl),
		chromedp.Sleep(500 * time.Millisecond), // Give chrome time to load font
	}); err != nil {
		slog.Error("Unable to navigate to URl")
		return nil, err
	}

	// Start to take screenshot
	allBytes := make([][]byte, len(elementSelectors))
	for i, selector := range elementSelectors {
		var imgByte []byte
		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Screenshot(selector, &imgByte, chromedp.NodeVisible),
		}); err != nil {
			slog.Error("Unable to take a screenshot of frame ", "i", i, "err", err)
			return nil, err
		}
		allBytes[i] = imgByte
	}
	return allBytes, nil
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
