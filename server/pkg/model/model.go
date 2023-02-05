package model

type ItemsAndTags struct {
	Items []Item `json:"items"`
	Tags  []Tag  `json:"tags"`
}

type Tag struct {
	Id          uint64           `json:"id,omitempty"`
	Title       string           `json:"title,omitempty" gorm:"uniqueIndex"`
	Items       []*Item          `json:"items,omitempty" gorm:"many2many:tag_items;"`
	Children    []*Tag           `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	ParentID    *uint64          `json:"parentId,omitempty"`
	Active      *bool            `json:"active,omitempty"`
	Selected    *bool            `json:"selected,omitempty"`
	Image       string           `json:"imageUrl,omitempty"`
	Annotations []*TagAnnotation `json:"tags_annotations,omitempty" gorm:"many2many:tags_annotations;"`
}

type Item struct {
	Id              uint64  `json:"id,omitempty"`
	Title           string  `json:"title,omitempty" gorm:"uniqueIndex:title_and_dir_idx;not null;default:null"`
	Origin          string  `json:"origin,omitempty" gorm:"uniqueIndex:title_and_dir_idx;not null;default:null"`
	DurationSeconds int     `json:"duration_seconds,omitempty"`
	Width           int     `json:"width,omitempty"`
	Height          int     `json:"height,omitempty"`
	CodecName       string  `json:"codec,omitempty"`
	Url             string  `json:"url,omitempty"`
	PreviewUrl      string  `json:"previewUrl,omitempty"`
	Covers          []Cover `json:"covers,omitempty"`
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
	Path       string `json:"path,omitempty" gorm:"primarykey"`
	Excluded   *bool  `json:"excluded,omitempty"`
	Tags       []*Tag `json:"tags,omitempty" gorm:"many2many:directory_tags;"`
	FilesCount int    `json:"filesCount,omitempty"`
	LastSynced int64  `json:"lastSynced,omitempty"`
}
