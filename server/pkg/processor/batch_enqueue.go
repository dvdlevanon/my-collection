package processor

import (
	"context"
	"my-collection/server/pkg/bl/items"
)

func (p *Processor) EnqueueAllItemsVideoMetadata(ctx context.Context, force bool) error {
	allItems, err := p.db.GetAllItems(ctx)
	if err != nil {
		return err
	}

	for _, item := range *allItems {
		if !force && item.DurationSeconds != 0 {
			modified, err := items.IsModified(&item, p)
			if !modified || err != nil {
				continue
			}
		}

		p.EnqueueItemVideoMetadata(ctx, item.Id, item.Title)
	}

	return nil
}

func (p *Processor) EnqueueAllItemsPreview(ctx context.Context, force bool) error {
	items, err := p.db.GetAllItems(ctx)
	if err != nil {
		return err
	}

	for _, item := range *items {
		if !force && item.PreviewUrl != "" {
			continue
		}

		p.EnqueueItemPreview(ctx, item.Id, item.Title)
	}

	return nil
}

func (p *Processor) EnqueueAllItemsCovers(ctx context.Context, force bool) error {
	items, err := p.db.GetAllItems(ctx)
	if err != nil {
		return err
	}

	for _, item := range *items {
		if !force && len(item.Covers) >= p.coversCount {
			continue
		}

		p.EnqueueItemCovers(ctx, item.Id, item.Title)
	}

	return nil
}

func (p *Processor) EnqueueAllItemsFileMetadata(ctx context.Context) error {
	allItems, err := p.db.GetAllItems(ctx)
	if err != nil {
		return err
	}

	for _, item := range *allItems {
		p.EnqueueItemFileMetadata(ctx, item.Id, item.Title)
	}

	return nil
}
