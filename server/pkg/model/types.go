package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
)

const ROOT_DIRECTORY_PATH = "<root>"

type FsNodeType int

const (
	FS_NODE_DIR = iota + 1
	FS_NODE_FILE
)

type TaskType int

const (
	REFRESH_COVER_TASK = iota
	REFRESH_PREVIEW_TASK
	REFRESH_METADATA_TASK
	REFRESH_FILE_TASK
	SET_MAIN_COVER
	CROP_FRAME
	CHANGE_RESOLUTION
)

func (t TaskType) ToDescription(title string) string {
	switch t {
	case REFRESH_COVER_TASK:
		return fmt.Sprintf("Extracting covers for %s", title)
	case REFRESH_PREVIEW_TASK:
		return fmt.Sprintf("Generating preview for %s", title)
	case REFRESH_METADATA_TASK:
		return fmt.Sprintf("Reading metadata for %s", title)
	case REFRESH_FILE_TASK:
		return fmt.Sprintf("Reading file metadata for %s", title)
	case SET_MAIN_COVER:
		return fmt.Sprintf("Setting main cover for %s", title)
	case CHANGE_RESOLUTION:
		return fmt.Sprintf("Changing resolution %s", title)
	default:
		return "unknown"
	}
}

func (t TaskType) String() string {
	switch t {
	case REFRESH_COVER_TASK:
		return "cover"
	case REFRESH_PREVIEW_TASK:
		return "preview"
	case REFRESH_METADATA_TASK:
		return "metadata"
	case REFRESH_FILE_TASK:
		return "file-metadata"
	case SET_MAIN_COVER:
		return "main-cover"
	case CHANGE_RESOLUTION:
		return "change-resolution"
	default:
		return "unknown"
	}
}

type PushMessageType int

const (
	// values in ws.js
	PUSH_PING           = 1
	PUSH_QUEUE_METADATA = 2
)

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	H int `json:"height"`
	W int `json:"width"`
}

func (t Rect) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Rect) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("type assertion to []byte failed %v", value)
	}

	return json.Unmarshal(b, &t)
}

type RectFloat struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	H float64 `json:"height"`
	W float64 `json:"width"`
}

func (t RectFloat) String() string {
	return fmt.Sprintf("%f:%f %f:%f", t.X, t.Y, t.W, t.H)
}

func (t RectFloat) Serialize() string {
	return fmt.Sprintf("%f %f %f %f", t.X, t.Y, t.W, t.H)
}

func DeserializeRectFloat(str string) (RectFloat, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 4 {
		return RectFloat{}, fmt.Errorf("bad rect format %s", str)
	}

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return RectFloat{}, fmt.Errorf("bad rect format %s %w", str, err)
	}

	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return RectFloat{}, fmt.Errorf("bad rect format %s %w", str, err)
	}

	w, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return RectFloat{}, fmt.Errorf("bad rect format %s %w", str, err)
	}

	h, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return RectFloat{}, fmt.Errorf("bad rect format %s %w", str, err)
	}

	return RectFloat{
		X: x,
		Y: y,
		W: w,
		H: h,
	}, nil
}
