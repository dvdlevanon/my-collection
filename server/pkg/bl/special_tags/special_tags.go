package special_tags

import (
	"my-collection/server/pkg/model"
)

var DailymixTag = &model.Tag{
	Title:    "DailyMix", // tags-utils.js
	ParentID: nil,
}

var MixOnDemandTag = &model.Tag{
	Title:    "Mod", // tags-utils.js
	ParentID: nil,
}

var SpecTag = &model.Tag{
	Title:          "Spec", // tags-utils.js
	ParentID:       nil,
	DisplayStyle:   "chip",
	DefaultSorting: "items-count",
}

func IsSpecial(tagId uint64) bool {
	return tagId == DailymixTag.Id || tagId == SpecTag.Id || tagId == MixOnDemandTag.Id
}
