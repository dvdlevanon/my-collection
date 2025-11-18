package model

type ItemsAndTags struct {
	Items []Item `json:"items"`
	Tags  []Tag  `json:"tags"`
}

type Tag struct {
	Id             uint64           `json:"id,omitempty"`
	Title          string           `json:"title,omitempty" gorm:"uniqueIndex:title_and_parent_idx"`
	ParentID       *uint64          `json:"parentId,omitempty" gorm:"uniqueIndex:title_and_parent_idx"`
	Items          []*Item          `json:"items,omitempty" gorm:"many2many:tag_items;"`
	Children       []*Tag           `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	Images         []*TagImage      `json:"images,omitempty"`
	Annotations    []*TagAnnotation `json:"tags_annotations,omitempty" gorm:"many2many:tags_annotations;"`
	DisplayStyle   string           `json:"display_style,omitempty"`
	DefaultSorting string           `json:"default_sorting,omitempty"`
	NoRandom       *bool            `json:"no_random,omitempty"`
}

type TagImageType struct {
	Id       uint64 `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty" gorm:"unique"`
	IconUrl  string `json:"iconUrl,omitempty"`
}

type TagImage struct {
	Id                uint64 `json:"id,omitempty"`
	Url               string `json:"url,omitempty"`
	TagId             uint64 `json:"tagId,omitempty"`
	ImageTypeId       uint64 `json:"imageType,omitempty"`
	ThumbnailRect     Rect   `json:"thumbnail_rect,omitempty" gorm:"type:json"`
	ImageNonce        int64  `json:"image_nonce,omitempty"`
	ThumbnailUrl      string `json:"thumbnail_url,omitempty"`
	ThumbnailUrlRect  Rect   `json:"thumbnail_url_rect,omitempty" gorm:"type:json"`
	ThumbnailUrlNonce int64  `json:"thumbnail_nonce,omitempty" gorm:"type:json"`
}

type Item struct {
	Id                    uint64  `json:"id,omitempty"`
	Title                 string  `json:"title,omitempty" gorm:"uniqueIndex:title_and_dir_idx"`
	Origin                string  `json:"origin,omitempty" gorm:"uniqueIndex:title_and_dir_idx"`
	DurationSeconds       float64 `json:"duration_seconds,omitempty"`
	FileSize              int64   `json:"file_size,omitempty"`
	Width                 int     `json:"width,omitempty"`
	Height                int     `json:"height,omitempty"`
	VideoCodecName        string  `json:"video_codec,omitempty"`
	AudioCodecName        string  `json:"audio_codec,omitempty"`
	Url                   string  `json:"url,omitempty"`
	PreviewUrl            string  `json:"preview_url,omitempty"`
	PreviewMode           string  `json:"preview_mode,omitempty"`
	LastModified          int64   `json:"last_modified,omitempty"`
	Covers                []Cover `json:"covers,omitempty"`
	MainCoverUrl          *string `json:"main_cover_url,omitempty"`
	MainCoverSecond       float64 `json:"main_cover_second,omitempty"`
	MainCoverNonce        int64   `json:"main_cover_nonce,omitempty"`
	Tags                  []*Tag  `json:"tags,omitempty" gorm:"many2many:tag_items;"`
	StartPosition         float64 `json:"start_position,omitempty"`
	EndPosition           float64 `json:"end_position,omitempty"`
	Highlights            []*Item `json:"highlights,omitempty" gorm:"foreignkey:HighlightParentItemId"`
	HighlightParentItemId *uint64 `json:"highlight_parent_id,omitempty"`
	SubItems              []*Item `json:"sub_items,omitempty" gorm:"foreignkey:MainItemId"`
	MainItemId            *uint64 `json:"main_item,omitempty"`
}

type Subtitle struct {
	Items []SubtitleItem `json:"items,omitempty"`
}

type SubtitleItem struct {
	StartMillis int64  `json:"start_millis"`
	EndMillis   int64  `json:"end_millis"`
	Text        string `json:"text"`
}

type SubtitleMetadata struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
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
	Path                 string `json:"path,omitempty" gorm:"primarykey"`
	Excluded             *bool  `json:"excluded,omitempty"`
	ProcessingStart      *int64 `json:"processingStart,omitempty"`
	AutoIncludeChildren  *bool  `json:"autoIncludeChildren,omitempty"`
	AutoIncludeHierarchy *bool  `json:"autoIncludeHierarchy,omitempty"`
	Tags                 []*Tag `json:"tags,omitempty" gorm:"many2many:directory_tags;"`
	FilesCount           *int   `json:"filesCount,omitempty"`
	LastSynced           int64  `json:"lastSynced,omitempty"`
}

type FsNode struct {
	Path     string     `json:"path,omitempty"`
	Type     FsNodeType `json:"type,omitempty"`
	Children []*FsNode  `json:"children,omitempty"`
	DirInfo  *Directory `json:"dirinfo,omitempty"`
}

type Task struct {
	Id              string   `json:"id,omitempty" gorm:"primarykey"`
	EnequeueTime    *int64   `json:"enqueueTime,omitempty"`
	ProcessingStart *int64   `json:"processingStart,omitempty"`
	ProcessingEnd   *int64   `json:"processingEnd,omitempty"`
	TaskType        TaskType `json:"type,omitempty"`
	IdParam         uint64   `json:"idParam,omitempty"`
	FloatParam      float64  `json:"floatParam,omitempty"`
	StringParam     string   `json:"stringParam,omitempty"`
	Description     string   `json:"description,omitempty" gorm:"-:all"`
}

type QueueMetadata struct {
	Size            *int64 `json:"size,omitempty"`
	Paused          *bool  `json:"paused,omitempty"`
	UnfinishedTasks *int64 `json:"unfinishedTasks,omitempty"`
}

type PushMessage struct {
	MessageType PushMessageType `json:"type,omitempty"`
	Payload     interface{}     `json:"payload,omitempty"`
}

type TagCustomCommand struct {
	Id      uint64 `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	TagId   uint64 `json:"tagId,omitempty"`
	Type    string `json:"type,omitempty"`
	Arg     string `json:"arg,omitempty"`
	Tooltip string `json:"tooltip,omitempty"`
	Icon    string `json:"icon,omitempty"`
}

type Stats struct {
	TagsCount            int64   `json:"tags_count,omitempty"`
	ItemsCount           int64   `json:"items_count,omitempty"`
	TotalDurationSeconds float64 `json:"total_duration_seconds,omitempty"`
}
