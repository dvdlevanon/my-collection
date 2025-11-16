package itemsoptimizer

import (
	"context"
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("itemsoptimizer")

type Processor interface {
	EnqueueChangeResolution(ctx context.Context, id uint64, newResolution string)
}

func New(ir model.ItemReader, processor Processor, maxResolution int) *ItemsOptimizer {
	return &ItemsOptimizer{
		ir:             ir,
		maxResolution:  maxResolution,
		processor:      processor,
		triggerChannel: make(chan bool),
	}
}

type ItemsOptimizer struct {
	ir             model.ItemReader
	maxResolution  int
	processor      Processor
	triggerChannel chan bool
}

func (d *ItemsOptimizer) EnqueueItemOptimizer() {
	d.triggerChannel <- true
}

func (d *ItemsOptimizer) Run(ctx context.Context) error {
	ctx = utils.ContextWithSubject(ctx, "items-optimizer")
	for {
		select {
		case <-d.triggerChannel:
			d.runItemsOptimizer(ctx)
		case <-ctx.Done():
			return nil
		}
	}
}

func (d *ItemsOptimizer) runItemsOptimizer(ctx context.Context) {
	logger.Infof("ItemsOptimizer started")
	if err := d.optimizeItems(ctx); err != nil {
		utils.LogError("Error in optimizeItems", err)
	}
	logger.Infof("ItemsOptimizer finished")
}

func (d *ItemsOptimizer) HandleItem(ctx context.Context, item *model.Item) {
	if item.Height <= d.maxResolution {
		return
	}

	logger.Infof("Item with high resolution %v", item.Title)
	d.processor.EnqueueChangeResolution(ctx, item.Id, ffmpeg.NewResolution(-1, d.maxResolution).String())
}

func (d *ItemsOptimizer) optimizeItems(ctx context.Context) error {
	allItems, err := d.ir.GetAllItems(ctx)
	if err != nil {
		return err
	}

	for _, item := range *allItems {
		d.HandleItem(ctx, &item)
	}

	return nil
}
