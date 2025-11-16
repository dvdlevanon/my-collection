package db

import (
	"fmt"
	"my-collection/server/pkg/model"
	"time"

	"github.com/op/go-logging"
)

var l = logging.MustGetLogger("db-logger")

type dbLogger struct {
	db *databaseImpl
}

func (d *dbLogger) log(operation string, start time.Time, err error, result interface{}) {
	duration := time.Since(start)
	ms := duration.Milliseconds()

	// Format duration with fixed width (6 chars: "1200ms")
	durationStr := fmt.Sprintf("%4dms", ms)

	// Format operation name with fixed width (20 chars)
	opStr := fmt.Sprintf("%-20s", operation)

	if err != nil {
		l.Errorf("%s %s error: %v", durationStr, opStr, err)
		return
	}

	// Format result description
	resultStr := formatResult(result)
	l.Infof("%s %s %s", durationStr, opStr, resultStr)
}

func formatResult(result interface{}) string {
	if result == nil {
		return "ok"
	}

	switch v := result.(type) {
	case *model.Directory:
		if v == nil {
			return "not found"
		}
		return fmt.Sprintf("dir: %s", v.Path)
	case *[]model.Directory:
		if v == nil {
			return "empty"
		}
		return fmt.Sprintf("%d dirs", len(*v))
	case *model.Item:
		if v == nil {
			return "not found"
		}
		return fmt.Sprintf("item id=%d", v.Id)
	case *[]model.Item:
		if v == nil {
			return "empty"
		}
		return fmt.Sprintf("%d items", len(*v))
	case *model.Tag:
		if v == nil {
			return "not found"
		}
		return fmt.Sprintf("tag id=%d", v.Id)
	case *[]model.Tag:
		if v == nil {
			return "empty"
		}
		return fmt.Sprintf("%d tags", len(*v))
	case *model.TagAnnotation:
		if v == nil {
			return "not found"
		}
		return fmt.Sprintf("annotation id=%d", v.Id)
	case []model.TagAnnotation:
		return fmt.Sprintf("%d annotations", len(v))
	case *model.TagImageType:
		if v == nil {
			return "not found"
		}
		return fmt.Sprintf("image type id=%d", v.Id)
	case *[]model.TagImageType:
		if v == nil {
			return "empty"
		}
		return fmt.Sprintf("%d image types", len(*v))
	case *[]model.TagCustomCommand:
		if v == nil {
			return "empty"
		}
		return fmt.Sprintf("%d commands", len(*v))
	case *[]model.Task:
		if v == nil {
			return "empty"
		}
		return fmt.Sprintf("%d tasks", len(*v))
	case *model.Task:
		if v == nil {
			return "not found"
		}
		return fmt.Sprintf("task id=%s", v.Id)
	case int64:
		return fmt.Sprintf("count=%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	default:
		return "ok"
	}
}

// Directory operations
func (d *dbLogger) CreateOrUpdateDirectory(directory *model.Directory) error {
	start := time.Now()
	err := d.db.CreateOrUpdateDirectory(directory)
	d.log("CreateUpdateDir", start, err, directory)
	return err
}

func (d *dbLogger) UpdateDirectory(directory *model.Directory) error {
	start := time.Now()
	err := d.db.UpdateDirectory(directory)
	d.log("UpdateDir", start, err, directory)
	return err
}

func (d *dbLogger) RemoveDirectory(path string) error {
	start := time.Now()
	err := d.db.RemoveDirectory(path)
	d.log("RemoveDir", start, err, fmt.Sprintf("path=%s", path))
	return err
}

func (d *dbLogger) RemoveTagFromDirectory(directoryPath string, tagId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagFromDirectory(directoryPath, tagId)
	d.log("RemoveTagFromDir", start, err, fmt.Sprintf("tag=%d", tagId))
	return err
}

func (d *dbLogger) GetDirectory(conds ...interface{}) (*model.Directory, error) {
	start := time.Now()
	result, err := d.db.GetDirectory(conds...)
	d.log("GetDir", start, err, result)
	return result, err
}

func (d *dbLogger) GetDirectories(conds ...interface{}) (*[]model.Directory, error) {
	start := time.Now()
	result, err := d.db.GetDirectories(conds...)
	d.log("GetDirs", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllDirectories() (*[]model.Directory, error) {
	start := time.Now()
	result, err := d.db.GetAllDirectories()
	d.log("GetAllDirs", start, err, result)
	return result, err
}

// Item operations
func (d *dbLogger) CreateOrUpdateItem(item *model.Item) error {
	start := time.Now()
	err := d.db.CreateOrUpdateItem(item)
	d.log("CreateUpdateItem", start, err, item)
	return err
}

func (d *dbLogger) UpdateItem(item *model.Item) error {
	start := time.Now()
	err := d.db.UpdateItem(item)
	d.log("UpdateItem", start, err, item)
	return err
}

func (d *dbLogger) RemoveItem(itemId uint64) error {
	start := time.Now()
	err := d.db.RemoveItem(itemId)
	d.log("RemoveItem", start, err, fmt.Sprintf("id=%d", itemId))
	return err
}

func (d *dbLogger) RemoveTagFromItem(itemId uint64, tagId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagFromItem(itemId, tagId)
	d.log("RemoveTagFromItem", start, err, fmt.Sprintf("item=%d tag=%d", itemId, tagId))
	return err
}

func (d *dbLogger) GetItem(conds ...interface{}) (*model.Item, error) {
	start := time.Now()
	result, err := d.db.GetItem(conds...)
	d.log("GetItem", start, err, result)
	return result, err
}

func (d *dbLogger) GetItems(conds ...interface{}) (*[]model.Item, error) {
	start := time.Now()
	result, err := d.db.GetItems(conds...)
	d.log("GetItems", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllItems() (*[]model.Item, error) {
	start := time.Now()
	result, err := d.db.GetAllItems()
	d.log("GetAllItems", start, err, result)
	return result, err
}

func (d *dbLogger) GetItemsCount() (int64, error) {
	start := time.Now()
	result, err := d.db.GetItemsCount()
	d.log("GetItemsCount", start, err, result)
	return result, err
}

func (d *dbLogger) GetTotalDurationSeconds() (float64, error) {
	start := time.Now()
	result, err := d.db.GetTotalDurationSeconds()
	d.log("GetTotalDuration", start, err, result)
	return result, err
}

// Tag Annotation operations
func (d *dbLogger) CreateTagAnnotation(tagAnnotation *model.TagAnnotation) error {
	start := time.Now()
	err := d.db.CreateTagAnnotation(tagAnnotation)
	d.log("CreateAnnotation", start, err, tagAnnotation)
	return err
}

func (d *dbLogger) RemoveTag(tagId uint64) error {
	start := time.Now()
	err := d.db.RemoveTag(tagId)
	d.log("RemoveTag", start, err, fmt.Sprintf("id=%d", tagId))
	return err
}

func (d *dbLogger) RemoveTagAnnotationFromTag(tagId uint64, annotationId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagAnnotationFromTag(tagId, annotationId)
	d.log("RemoveAnnotation", start, err, fmt.Sprintf("tag=%d ann=%d", tagId, annotationId))
	return err
}

func (d *dbLogger) GetTagAnnotation(conds ...interface{}) (*model.TagAnnotation, error) {
	start := time.Now()
	result, err := d.db.GetTagAnnotation(conds...)
	d.log("GetAnnotation", start, err, result)
	return result, err
}

func (d *dbLogger) GetTagAnnotations(tagId uint64) ([]model.TagAnnotation, error) {
	start := time.Now()
	result, err := d.db.GetTagAnnotations(tagId)
	d.log("GetAnnotations", start, err, result)
	return result, err
}

// TagImageType operations
func (d *dbLogger) CreateOrUpdateTagImageType(tit *model.TagImageType) error {
	start := time.Now()
	err := d.db.CreateOrUpdateTagImageType(tit)
	d.log("CreateUpdateImgType", start, err, tit)
	return err
}

func (d *dbLogger) GetTagImageType(conds ...interface{}) (*model.TagImageType, error) {
	start := time.Now()
	result, err := d.db.GetTagImageType(conds...)
	d.log("GetImgType", start, err, result)
	return result, err
}

func (d *dbLogger) GetTagImageTypes(conds ...interface{}) (*[]model.TagImageType, error) {
	start := time.Now()
	result, err := d.db.GetTagImageTypes(conds...)
	d.log("GetImgTypes", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllTagImageTypes() (*[]model.TagImageType, error) {
	start := time.Now()
	result, err := d.db.GetAllTagImageTypes()
	d.log("GetAllImgTypes", start, err, result)
	return result, err
}

// TagCustomCommand operations
func (d *dbLogger) CreateOrUpdateTagCustomCommand(command *model.TagCustomCommand) error {
	start := time.Now()
	err := d.db.CreateOrUpdateTagCustomCommand(command)
	d.log("CreateUpdateCmd", start, err, command)
	return err
}

func (d *dbLogger) GetTagCustomCommand(conds ...interface{}) (*[]model.TagCustomCommand, error) {
	start := time.Now()
	result, err := d.db.GetTagCustomCommand(conds...)
	d.log("GetCmd", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllTagCustomCommands() (*[]model.TagCustomCommand, error) {
	start := time.Now()
	result, err := d.db.GetAllTagCustomCommands()
	d.log("GetAllCmds", start, err, result)
	return result, err
}

// Tag operations
func (d *dbLogger) CreateOrUpdateTag(tag *model.Tag) error {
	start := time.Now()
	err := d.db.CreateOrUpdateTag(tag)
	d.log("CreateUpdateTag", start, err, tag)
	return err
}

func (d *dbLogger) UpdateTag(tag *model.Tag) error {
	start := time.Now()
	err := d.db.UpdateTag(tag)
	d.log("UpdateTag", start, err, tag)
	return err
}

func (d *dbLogger) GetTag(conds ...interface{}) (*model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetTag(conds...)
	d.log("GetTag", start, err, result)
	return result, err
}

func (d *dbLogger) GetTagsWithoutChildren(conds ...interface{}) (*[]model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetTagsWithoutChildren(conds...)
	d.log("GetTagsNoChildren", start, err, result)
	return result, err
}

func (d *dbLogger) GetTags(conds ...interface{}) (*[]model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetTags(conds...)
	d.log("GetTags", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllTags() (*[]model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetAllTags()
	d.log("GetAllTags", start, err, result)
	return result, err
}

func (d *dbLogger) RemoveTagImageFromTag(tagId uint64, imageId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagImageFromTag(tagId, imageId)
	d.log("RemoveTagImage", start, err, fmt.Sprintf("tag=%d img=%d", tagId, imageId))
	return err
}

func (d *dbLogger) UpdateTagImage(image *model.TagImage) error {
	start := time.Now()
	err := d.db.UpdateTagImage(image)
	d.log("UpdateTagImage", start, err, image)
	return err
}

func (d *dbLogger) GetTagsCount() (int64, error) {
	start := time.Now()
	result, err := d.db.GetTagsCount()
	d.log("GetTagsCount", start, err, result)
	return result, err
}

// Task operations
func (d *dbLogger) CreateTask(task *model.Task) error {
	start := time.Now()
	err := d.db.CreateTask(task)
	d.log("CreateTask", start, err, task)
	return err
}

func (d *dbLogger) UpdateTask(task *model.Task) error {
	start := time.Now()
	err := d.db.UpdateTask(task)
	d.log("UpdateTask", start, err, task)
	return err
}

func (d *dbLogger) RemoveTasks(conds ...interface{}) error {
	start := time.Now()
	err := d.db.RemoveTasks(conds...)
	d.log("RemoveTasks", start, err, nil)
	return err
}

func (d *dbLogger) TasksCount(query interface{}, conds ...interface{}) (int64, error) {
	start := time.Now()
	result, err := d.db.TasksCount(query, conds...)
	d.log("TasksCount", start, err, result)
	return result, err
}

func (d *dbLogger) GetTasks(offset int, limit int) (*[]model.Task, error) {
	start := time.Now()
	result, err := d.db.GetTasks(offset, limit)
	d.log("GetTasks", start, err, result)
	return result, err
}
