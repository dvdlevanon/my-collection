package suggestions

import (
	"context"
	"math/rand"
	"my-collection/server/pkg/bl/items"
	"my-collection/server/pkg/bl/special_tags"
	"my-collection/server/pkg/bl/tags"
	"my-collection/server/pkg/model"
)

func GetSuggestionsForItem(ctx context.Context, ir model.ItemReader, tr model.TagReader, itemId uint64, count int) ([]*model.Item, error) {
	item, err := ir.GetItem(ctx, itemId)
	if err != nil {
		return nil, err
	}

	t, err := tags.GetFullTags(ctx, tr, item.Tags)
	if err != nil {
		return nil, err
	}

	return GetSuggestionsForTags(ctx, ir, tr, t, count)
}

func GetSuggestionsForTags(ctx context.Context, ir model.ItemReader, tr model.TagReader, tags *[]model.Tag, count int) ([]*model.Item, error) {
	relatedItems, err := getItemsOfTags(ctx, ir, tags)
	if err != nil {
		return nil, err
	}

	if len(relatedItems) < count {
		randomItems, err := items.GetRandomItems(ctx, ir, count-len(relatedItems), func(item *model.Item) bool {
			return !items.IsHighlight(item) && !items.IsSplittedItem(item)
		})
		if err != nil {
			return nil, err
		}

		relatedItems = append(relatedItems, randomItems...)
	}

	// Safety check: if we still don't have enough items, adjust count
	if len(relatedItems) == 0 {
		return []*model.Item{}, nil
	}

	if count > len(relatedItems) {
		count = len(relatedItems)
	}

	resultIndexes := getUniqueRandomNumbers(len(relatedItems), count)
	result := make([]*model.Item, len(resultIndexes))
	for i, n := range resultIndexes {
		result[i] = relatedItems[n-1]
	}

	return result, nil
}

func getItemsOfTags(ctx context.Context, ir model.ItemReader, t *[]model.Tag) ([]*model.Item, error) {
	relatedItems := make([]*model.Item, 0)

	for _, tag := range *t {
		if special_tags.IsSpecial(*tag.ParentID) {
			continue
		}

		tagItems, err := tags.GetItems(ctx, ir, &tag)
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
	// Safety check: can't generate more unique numbers than available
	if count > max {
		count = max
	}

	// If no numbers can be generated, return empty slice
	if max <= 0 || count <= 0 {
		return []int{}
	}

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
