package utils

import (
	"bufio"
	"bytes"
	"image"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/draw"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"github.com/h2non/filetype"
	"github.com/op/go-logging"
	"github.com/saintfish/chardet"
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

func BaseRemoveExtension(filename string) string {
	ext := filepath.Ext(filepath.Base(filename))
	return filename[:len(filename)-len(ext)]
}

func DetectEncodingAndRead(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(data)
	if err != nil {
		return "", err
	}

	if result.Charset != "UTF-8" {
		var decoder *encoding.Decoder
		switch result.Charset {
		case "ISO-8859-1", "windows-1252":
			decoder = charmap.Windows1252.NewDecoder()
		case "ISO-8859-15":
			decoder = charmap.ISO8859_15.NewDecoder()
		case "ISO-8859-8-I":
			decoder = charmap.ISO8859_8I.NewDecoder()
		case "ISO-8859-2":
			decoder = charmap.ISO8859_2.NewDecoder()
		case "ISO-8859-3":
			decoder = charmap.ISO8859_3.NewDecoder()
		case "ISO-8859-4":
			decoder = charmap.ISO8859_4.NewDecoder()
		case "ISO-8859-5":
			decoder = charmap.ISO8859_5.NewDecoder()
		case "ISO-8859-6":
			decoder = charmap.ISO8859_6.NewDecoder()
		case "ISO-8859-7":
			decoder = charmap.ISO8859_7.NewDecoder()
		case "ISO-8859-8":
			decoder = charmap.ISO8859_8.NewDecoder()
		case "ISO-8859-9":
			decoder = charmap.ISO8859_9.NewDecoder()
		case "ISO-8859-10":
			decoder = charmap.ISO8859_10.NewDecoder()
		case "ISO-8859-13":
			decoder = charmap.ISO8859_13.NewDecoder()
		case "ISO-8859-14":
			decoder = charmap.ISO8859_14.NewDecoder()
		case "ISO-8859-16":
			decoder = charmap.ISO8859_16.NewDecoder()
		case "windows-1250":
			decoder = charmap.Windows1250.NewDecoder()
		case "windows-1251":
			decoder = charmap.Windows1251.NewDecoder()
		case "windows-1253":
			decoder = charmap.Windows1253.NewDecoder()
		case "windows-1254":
			decoder = charmap.Windows1254.NewDecoder()
		case "windows-1255":
			decoder = charmap.Windows1255.NewDecoder()
		case "windows-1256":
			decoder = charmap.Windows1256.NewDecoder()
		case "windows-1257":
			decoder = charmap.Windows1257.NewDecoder()
		case "windows-1258":
			decoder = charmap.Windows1258.NewDecoder()
		case "windows-874":
			decoder = charmap.Windows874.NewDecoder()
		case "KOI8-R":
			decoder = charmap.KOI8R.NewDecoder()
		case "KOI8-U":
			decoder = charmap.KOI8U.NewDecoder()
		case "macintosh":
			decoder = charmap.Macintosh.NewDecoder()
		case "IBM437":
			decoder = charmap.CodePage437.NewDecoder()
		case "IBM850":
			decoder = charmap.CodePage850.NewDecoder()
		case "IBM852":
			decoder = charmap.CodePage852.NewDecoder()
		case "IBM855":
			decoder = charmap.CodePage855.NewDecoder()
		case "IBM858":
			decoder = charmap.CodePage858.NewDecoder()
		case "IBM860":
			decoder = charmap.CodePage860.NewDecoder()
		case "IBM862":
			decoder = charmap.CodePage862.NewDecoder()
		case "IBM863":
			decoder = charmap.CodePage863.NewDecoder()
		case "IBM865":
			decoder = charmap.CodePage865.NewDecoder()
		case "IBM866":
			decoder = charmap.CodePage866.NewDecoder()
		default:
			return string(data), nil // fallback
		}

		reader := transform.NewReader(bytes.NewReader(data), decoder)
		decoded, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}

	return string(data), nil
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
