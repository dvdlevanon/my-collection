package model

type ItemsAndTags struct {
	Items []Item `json:"items"`
	Tags  []Tag  `json:"tags"`
}

type Tag struct {
	Id       uint64  `json:"id,omitempty"`
	Title    string  `json:"title,omitempty" gorm:"uniqueIndex"`
	Items    []*Item `json:"items,omitempty" gorm:"many2many:tag_items;"`
	Children []*Tag  `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	ParentID *uint64 `json:"parentId,omitempty"`
	Active   *bool   `json:"active,omitempty"`
	Selected *bool   `json:"selected,omitempty"`
	Image    string  `json:"imageUrl,omitempty"`
}

type Item struct {
	Id         uint64  `json:"id,omitempty"`
	Title      string  `json:"title,omitempty" gorm:"uniqueIndex"`
	Url        string  `json:"url,omitempty"`
	PreviewUrl string  `json:"previewUrl,omitempty"`
	Covers     []Cover `json:"covers,omitempty"`
	Tags       []*Tag  `json:"tags,omitempty" gorm:"many2many:tag_items;"`
}

type Cover struct {
	Id     uint64 `json:"id,omitempty"`
	Url    string `json:"url,omitempty"`
	ItemId uint64 `json:"itemId,omitempty"`
}

type FileUrl struct {
	Url string `json:"url,omitempty"`
}
