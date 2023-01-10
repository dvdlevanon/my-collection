package model

type Tag struct {
	Id       uint64  `json:"id,omitempty"`
	Title    string  `json:"title,omitempty" gorm:"uniqueIndex"`
	Items    []*Item `json:"items,omitempty" gorm:"many2many:tag_items;"`
	Children []*Tag  `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	ParentID *uint64 `json:"parentId,omitempty"`
}

type Item struct {
	Id       uint64    `json:"id,omitempty"`
	Title    string    `json:"title,omitempty" gorm:"uniqueIndex"`
	Url      string    `json:"url,omitempty"`
	Previews []Preview `json:"previews,omitempty"`
	Tags     []*Tag    `json:"tags,omitempty" gorm:"many2many:tag_items;"`
}

type Preview struct {
	Id     uint64 `json:"id,omitempty"`
	Url    string `json:"url,omitempty"`
	ItemId uint64 `json:"itemId,omitempty"`
}
