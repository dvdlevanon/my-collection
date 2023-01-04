package model

type Tag struct {
	Id       uint64  `json:"id,omitempty"`
	Title    string  `json:"title,omitempty" gorm:"uniqueIndex"`
	Items    []*Item `json:"items,omitempty" gorm:"many2many:tag_items;"`
	Children []*Tag  `json:"children,omitempty" gorm:"foreignkey:ParentID"`
	ParentID *uint64 `json:"parentId,omitempty"`
}

type Item struct {
	Id    uint64 `json:"id,omitempty"`
	Title string `json:"title,omitempty" gorm:"uniqueIndex"`
	Cover string `json:"cover,omitempty"`
	Tags  []*Tag `json:"tags,omitempty" gorm:"many2many:tag_items;"`
}
