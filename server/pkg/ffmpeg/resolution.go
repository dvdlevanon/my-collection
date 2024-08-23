package ffmpeg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

func ResolutionFromString(resString string) (Resolution, error) {
	parts := strings.Split(resString, ":")
	if len(parts) != 2 {
		return Resolution{}, errors.Errorf("invalid resolution format %s", resString)
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return Resolution{}, err
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return Resolution{}, err
	}

	return Resolution{
		Width:  width,
		Height: height,
	}, nil
}

func NewResolution(width int, height int) Resolution {
	return Resolution{
		Width:  width,
		Height: height,
	}
}

type Resolution struct {
	Width  int
	Height int
}

func (r Resolution) String() string {
	return fmt.Sprintf("%d:%d", r.Width, r.Height)
}
