package model

import "fmt"

type TaskType int

const (
	REFRESH_COVER_TASK = iota
	REFRESH_PREVIEW_TASK
	REFRESH_METADATA_TASK
	SET_MAIN_COVER
)

func (t TaskType) ToDescription(title string) string {
	switch t {
	case REFRESH_COVER_TASK:
		return fmt.Sprintf("Extracting covers for %s", title)
	case REFRESH_PREVIEW_TASK:
		return fmt.Sprintf("Generating preview for %s", title)
	case REFRESH_METADATA_TASK:
		return fmt.Sprintf("Reading metadata for %s", title)
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
	default:
		return "unknown"
	}
}

type PushMessageType int

const (
	PUSH_PING = iota
	PUSH_QUEUE_METADATA
)
