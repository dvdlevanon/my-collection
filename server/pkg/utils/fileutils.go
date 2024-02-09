package utils

import (
	"image"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/draw"

	"github.com/h2non/filetype"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("utils")

func IsVideo(trustFileExtenssion bool, path string) bool {
	if trustFileExtenssion {
		path = strings.ToLower(path)
		return strings.HasSuffix(path, ".avi") ||
			strings.HasSuffix(path, ".mkv") ||
			strings.HasSuffix(path, ".mpg") ||
			strings.HasSuffix(path, ".mpeg") ||
			strings.HasSuffix(path, ".wmv") ||
			strings.HasSuffix(path, ".mp4")
	}

	file, err := os.Open(path)
	if err != nil {
		logger.Errorf("Error opening file for reading %s - %t", file, err)
		return false
	}

	stat, err := file.Stat()
	if err != nil {
		logger.Errorf("Error getting stats of file %s - %t", path, err)
		return false
	}

	header := make([]byte, int(math.Max(float64(stat.Size())-1, 1024)))
	_, err = file.Read(header)
	if err != nil {
		logger.Errorf("Error reading from file %s - %t", path, err)
		return false
	}

	return filetype.IsVideo(header)
}

func DownloadFile(url string, targetFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func ExtractImage(imageFile, outputFile string, x, y, width, height, destWidth, destHeight int) error {
	inFile, err := os.Open(imageFile)
	if err != nil {
		return err
	}
	defer inFile.Close()

	img, _, err := image.Decode(inFile)
	if err != nil {
		return err
	}

	rect := image.Rect(x, y, x+width, y+height)
	subImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)

	resizedImg := image.NewRGBA(image.Rect(0, 0, destWidth, destHeight))
	draw.ApproxBiLinear.Scale(resizedImg, resizedImg.Bounds(), subImg, subImg.Bounds(), draw.Over, nil)

	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return png.Encode(outFile, resizedImg)
}
