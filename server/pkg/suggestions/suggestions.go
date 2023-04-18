package suggestions

import (
	"math/rand"
	"my-collection/server/pkg/automix"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
)

func GetSuggestionsForItem(ir model.ItemReader, tr model.TagReader, itemId uint64, count int) ([]*model.Item, error) {
	item, err := ir.GetItem(itemId)
	if err != nil {
		return nil, err
	}

	t, err := tags.GetFullTags(tr, item.Tags)
	if err != nil {
		return nil, err
	}

	relatedItems, err := getItemsOfTags(ir, t)
	if err != nil {
		return nil, err
	}

	if len(relatedItems) < count {
		randomItems, err := items.GetRandomItems(ir, count-len(relatedItems))
		if err != nil {
			return nil, err
		}

		relatedItems = append(relatedItems, randomItems...)
	}

	resultIndexes := getUniqueRandomNumbers(len(relatedItems), count)
	result := make([]*model.Item, count)
	for i, n := range resultIndexes {
		result[i] = relatedItems[n-1]
	}

	return result, nil
}

func getItemsOfTags(ir model.ItemReader, t *[]model.Tag) ([]*model.Item, error) {
	relatedItems := make([]*model.Item, 0)

	for _, tag := range *t {
		if *tag.ParentID == automix.GetDailymixTagId() {
			continue
		}

		tagItems, err := tags.GetItems(ir, &tag)
		if err != nil {
			return nil, err
		}

		for _, item := range *tagItems {
			if !items.ItemExists(relatedItems, &item) {
				fake := item
				relatedItems = append(relatedItems, &fake)
			}
		}
	}

	return relatedItems, nil
}

func getUniqueRandomNumbers(max int, count int) []int {
	result := make([]int, count)

	for i := 0; i < count; i++ {
		num := rand.Intn(max) + 1
		for numberExists(result, num) {
			num = rand.Intn(max) + 1
		}

		result[i] = num
	}

	return result
}

func numberExists(nums []int, num int) bool {
	for _, n := range nums {
		if n == num {
			return true
		}
	}

	return false
}
