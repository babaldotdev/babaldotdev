package gif_creator

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
)

func Create(htmlBuf bytes.Buffer, nFrame int, fileName string) error {
	imgDataList, err := getImageBytes(htmlBuf, nFrame)
	if err != nil {
		slog.Error("Unable to generate image bytes", "err", err)
		return err
	}
	outGif, err := prepareGif(imgDataList)
	if err != nil {
		slog.Error("Unable to prepare GIF", "err", err)
		return err
	}
	return writeGifToFile(outGif, fileName)
}

func writeGifToFile(outGif *gif.GIF, fileName string) error {
	f, _ := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0o644)

	err := gif.EncodeAll(f, outGif)
	if err != nil {
		slog.Error("Unable to encode GIF", "err", err)
		return err
	}

	err = f.Close()
	if err != nil {
		slog.Error("Error unable to close the file")
		return err
	}

	return nil
}

func prepareGif(imgDataList [][]byte) (*gif.GIF, error) {
	outGif := &gif.GIF{}
	for _, imgData := range imgDataList {
		// Decode PNG image
		inImg, err := png.Decode(bytes.NewReader(imgData))
		if err != nil {
			slog.Error("Unable to decode image from bytes %v \n", err)
			return nil, err
		}

		// Convert to Paletted image
		bounds := inImg.Bounds()
		palettedImg := image.NewPaletted(bounds, palette.Plan9) // Use Plan9 palette (or another palette as needed)
		draw.Draw(palettedImg, bounds, inImg, bounds.Min, draw.Src)

		outGif.Image = append(outGif.Image, palettedImg)
		outGif.Delay = append(outGif.Delay, 0)
	}
	return outGif, nil
}

func getImageBytes(htmlBuf bytes.Buffer, nFrame int) ([][]byte, error) {
	framesSelector := make([]string, nFrame)
	for i := 0; i < nFrame; i++ {
		framesSelector[i] = fmt.Sprintf("#frame-%d", i)
	}
	return html_to_image.ImageByteMulti(
		string(htmlBuf.Bytes()),
		framesSelector)
}
