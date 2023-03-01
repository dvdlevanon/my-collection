package model

import "fmt"

type ItemsAndTags struct {
	Items []Item `json:"items"`
	Tags  []Tag  `json:"tags"`
}

type Tag struct {
	Id          uint64           `json:"id,omitempty"`
	Title       string           `json:"title,omitempty" gorm:"uniqueIndex:title_and_parent_idx"`
	ParentID    *uint64          `json:"parentId,omitempty" gorm:"uniqueIndex:title_and_parent_idx"`
	Items       []*Item          `json:"items,omitempty" gorm:"many2many:tag_items;"`
	Children    []*Tag           `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	Active      *bool            `json:"active,omitempty"`
	Selected    *bool            `json:"selected,omitempty"`
	Image       string           `json:"imageUrl,omitempty"`
	Annotations []*TagAnnotation `json:"tags_annotations,omitempty" gorm:"many2many:tags_annotations;"`
}

type Item struct {
	Id              uint64  `json:"id,omitempty"`
	Title           string  `json:"title,omitempty" gorm:"uniqueIndex:title_and_dir_idx"`
	Origin          string  `json:"origin,omitempty" gorm:"uniqueIndex:title_and_dir_idx"`
	DurationSeconds int     `json:"duration_seconds,omitempty"`
	Width           int     `json:"width,omitempty"`
	Height          int     `json:"height,omitempty"`
	VideoCodecName  string  `json:"videoCodec,omitempty"`
	AudioCodecName  string  `json:"audioCodec,omitempty"`
	Url             string  `json:"url,omitempty"`
	PreviewUrl      string  `json:"previewUrl,omitempty"`
	LastModified    int64   `json:"lastModified,omitempty"`
	Covers          []Cover `json:"covers,omitempty"`
	MainCoverUrl    *string `json:"mainCoverUrl,omitempty"`
	Tags            []*Tag  `json:"tags,omitempty" gorm:"many2many:tag_items;"`
}

type Cover struct {
	Id     uint64 `json:"id,omitempty"`
	Url    string `json:"url,omitempty"`
	ItemId uint64 `json:"itemId,omitempty"`
}

type FileUrl struct {
	Url string `json:"url,omitempty"`
}

type TagAnnotation struct {
	Id    uint64 `json:"id,omitempty" gorm:"primarykey"`
	Title string `json:"title,omitempty" gorm:"unique"`
	Tags  []*Tag `json:"tags,omitempty" gorm:"many2many:tags_annotations;"`
}

type Directory struct {
	Path            string `json:"path,omitempty" gorm:"primarykey"`
	Excluded        *bool  `json:"excluded,omitempty"`
	ProcessingStart *int64 `json:"processingStart,omitempty"`
	Tags            []*Tag `json:"tags,omitempty" gorm:"many2many:directory_tags;"`
	FilesCount      *int   `json:"filesCount,omitempty"`
	LastSynced      int64  `json:"lastSynced,omitempty"`
}

type TaskType int

const (
	REFRESH_COVER_TASK = iota
	REFRESH_PREVIEW_TASK
	REFRESH_METADATA_TASK
	SET_MAIN_COVER
)

type Task struct {
	Id              string   `json:"id,omitempty" gorm:"primarykey"`
	EnequeueTime    *int64   `json:"enqueueTime,omitempty"`
	ProcessingStart *int64   `json:"processingStart,omitempty"`
	ProcessingEnd   *int64   `json:"processingEnd,omitempty"`
	TaskType        TaskType `json:"type,omitempty"`
	IdParam         uint64   `json:"idParam,omitempty"`
	FloatParam      float64  `json:"floatParam,omitempty"`
	Description     string   `json:"description,omitempty" gorm:"-:all"`
}

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

type QueueMetadata struct {
	Size            *int64 `json:"size,omitempty"`
	Paused          *bool  `json:"paused,omitempty"`
	UnfinishedTasks *int64 `json:"unfinishedTasks,omitempty"`
}

type PushMessageType int

const (
	PUSH_PING = iota
	PUSH_QUEUE_METADATA
)

type PushMessage struct {
	MessageType PushMessageType `json:"type,omitempty"`
	Payload     interface{}     `json:"payload,omitempty"`
}
