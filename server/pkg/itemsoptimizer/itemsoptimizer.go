package itemsoptimizer

import (
	"my-collection/server/pkg/ffmpeg"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/processor"
	"my-collection/server/pkg/utils"

	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("itemsoptimizer")

func New(ir model.ItemReader, processor processor.Processor, maxResolution int) *ItemsOptimizer {
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
	processor      processor.Processor
	triggerChannel chan bool
}

func (d *ItemsOptimizer) Trigger() {
	d.triggerChannel <- true
}

func (d *ItemsOptimizer) Run() {
	for range d.triggerChannel {
		d.runItemsOptimizer()
	}
}

func (d *ItemsOptimizer) runItemsOptimizer() {
	logger.Infof("ItemsOptimizer started")
	if err := d.optimizeItems(); err != nil {
		utils.LogError(err)
	}
	logger.Infof("ItemsOptimizer finished")
}

func (d *ItemsOptimizer) HandleItem(item *model.Item) {
	if item.Height <= d.maxResolution {
		return
	}

	logger.Infof("Item with high resolution %v", item.Title)
	d.processor.EnqueueChangeResolution(item.Id, ffmpeg.NewResolution(-1, d.maxResolution).String())
}

func (d *ItemsOptimizer) optimizeItems() error {
	allItems, err := d.ir.GetAllItems()
	if err != nil {
		return err
	}

	for _, item := range *allItems {
		d.HandleItem(&item)
	}

	return nil
}
