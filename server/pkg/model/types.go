package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/go-errors/errors"
)

type TaskType int

const (
	REFRESH_COVER_TASK = iota
	REFRESH_PREVIEW_TASK
	REFRESH_METADATA_TASK
	REFRESH_FILE_TASK
	SET_MAIN_COVER
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
