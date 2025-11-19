package db

import (
	"context"
	"fmt"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/utils"
	"time"

	"github.com/op/go-logging"
)

var l = logging.MustGetLogger("db-logger")

type dbLogger struct {
	db *databaseImpl
}

func (d *dbLogger) log(ctx context.Context, operation string, start time.Time, err error, result interface{}) {
	duration := time.Since(start)
	ms := duration.Milliseconds()

	// Format duration with fixed width (6 chars: "1200ms")
	durationStr := fmt.Sprintf("%4dms", ms)

	// Format operation name with fixed width (20 chars)
	opStr := fmt.Sprintf("%-20s", operation)

	subject := utils.GetSubject(ctx)
	subjectStr := fmt.Sprintf("%-20s", subject)

	if err != nil {
		l.Errorf("[%s] %s %s error: %v", subjectStr, durationStr, opStr, err)
		return
	}

	// Format result description
	resultStr := formatResult(result)
	l.Debugf("[%s] %s %s %s", subjectStr, durationStr, opStr, resultStr)
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
func (d *dbLogger) CreateOrUpdateDirectory(ctx context.Context, directory *model.Directory) error {
	start := time.Now()
	err := d.db.CreateOrUpdateDirectory(ctx, directory)
	d.log(ctx, "CreateUpdateDir", start, err, directory)
	return err
}

func (d *dbLogger) UpdateDirectory(ctx context.Context, directory *model.Directory) error {
	start := time.Now()
	err := d.db.UpdateDirectory(ctx, directory)
	d.log(ctx, "UpdateDir", start, err, directory)
	return err
}

func (d *dbLogger) RemoveDirectory(ctx context.Context, path string) error {
	start := time.Now()
	err := d.db.RemoveDirectory(ctx, path)
	d.log(ctx, "RemoveDir", start, err, fmt.Sprintf("path=%s", path))
	return err
}

func (d *dbLogger) RemoveTagFromDirectory(ctx context.Context, directoryPath string, tagId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagFromDirectory(ctx, directoryPath, tagId)
	d.log(ctx, "RemoveTagFromDir", start, err, fmt.Sprintf("tag=%d", tagId))
	return err
}

func (d *dbLogger) GetDirectory(ctx context.Context, conds ...interface{}) (*model.Directory, error) {
	start := time.Now()
	result, err := d.db.GetDirectory(ctx, conds...)
	d.log(ctx, "GetDir", start, err, result)
	return result, err
}

func (d *dbLogger) GetDirectories(ctx context.Context, conds ...interface{}) (*[]model.Directory, error) {
	start := time.Now()
	result, err := d.db.GetDirectories(ctx, conds...)
	d.log(ctx, "GetDirs", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllDirectories(ctx context.Context) (*[]model.Directory, error) {
	start := time.Now()
	result, err := d.db.GetAllDirectories(ctx)
	d.log(ctx, "GetAllDirs", start, err, result)
	return result, err
}

// Item operations
func (d *dbLogger) CreateOrUpdateItem(ctx context.Context, item *model.Item) error {
	start := time.Now()
	err := d.db.CreateOrUpdateItem(ctx, item)
	d.log(ctx, "CreateUpdateItem", start, err, item)
	return err
}

func (d *dbLogger) UpdateItem(ctx context.Context, item *model.Item) error {
	start := time.Now()
	err := d.db.UpdateItem(ctx, item)
	d.log(ctx, "UpdateItem", start, err, item)
	return err
}

func (d *dbLogger) RemoveItem(ctx context.Context, itemId uint64) error {
	start := time.Now()
	err := d.db.RemoveItem(ctx, itemId)
	d.log(ctx, "RemoveItem", start, err, fmt.Sprintf("id=%d", itemId))
	return err
}

func (d *dbLogger) RemoveTagFromItem(ctx context.Context, itemId uint64, tagId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagFromItem(ctx, itemId, tagId)
	d.log(ctx, "RemoveTagFromItem", start, err, fmt.Sprintf("item=%d tag=%d", itemId, tagId))
	return err
}

func (d *dbLogger) GetItem(ctx context.Context, conds ...interface{}) (*model.Item, error) {
	start := time.Now()
	result, err := d.db.GetItem(ctx, conds...)
	d.log(ctx, "GetItem", start, err, result)
	return result, err
}

func (d *dbLogger) GetItems(ctx context.Context, conds ...interface{}) (*[]model.Item, error) {
	start := time.Now()
	result, err := d.db.GetItems(ctx, conds...)
	d.log(ctx, "GetItems", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllItems(ctx context.Context) (*[]model.Item, error) {
	start := time.Now()
	result, err := d.db.GetAllItems(ctx)
	d.log(ctx, "GetAllItems", start, err, result)
	return result, err
}

func (d *dbLogger) GetItemsCount(ctx context.Context) (int64, error) {
	start := time.Now()
	result, err := d.db.GetItemsCount(ctx)
	d.log(ctx, "GetItemsCount", start, err, result)
	return result, err
}

func (d *dbLogger) GetTotalDurationSeconds(ctx context.Context) (float64, error) {
	start := time.Now()
	result, err := d.db.GetTotalDurationSeconds(ctx)
	d.log(ctx, "GetTotalDuration", start, err, result)
	return result, err
}

// Tag Annotation operations
func (d *dbLogger) CreateTagAnnotation(ctx context.Context, tagAnnotation *model.TagAnnotation) error {
	start := time.Now()
	err := d.db.CreateTagAnnotation(ctx, tagAnnotation)
	d.log(ctx, "CreateAnnotation", start, err, tagAnnotation)
	return err
}

func (d *dbLogger) RemoveTag(ctx context.Context, tagId uint64) error {
	start := time.Now()
	err := d.db.RemoveTag(ctx, tagId)
	d.log(ctx, "RemoveTag", start, err, fmt.Sprintf("id=%d", tagId))
	return err
}

func (d *dbLogger) RemoveTagAnnotationFromTag(ctx context.Context, tagId uint64, annotationId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagAnnotationFromTag(ctx, tagId, annotationId)
	d.log(ctx, "RemoveAnnotation", start, err, fmt.Sprintf("tag=%d ann=%d", tagId, annotationId))
	return err
}

func (d *dbLogger) GetTagAnnotation(ctx context.Context, conds ...interface{}) (*model.TagAnnotation, error) {
	start := time.Now()
	result, err := d.db.GetTagAnnotation(ctx, conds...)
	d.log(ctx, "GetAnnotation", start, err, result)
	return result, err
}

func (d *dbLogger) GetTagAnnotations(ctx context.Context, tagId uint64) ([]model.TagAnnotation, error) {
	start := time.Now()
	result, err := d.db.GetTagAnnotations(ctx, tagId)
	d.log(ctx, "GetAnnotations", start, err, result)
	return result, err
}

// TagImageType operations
func (d *dbLogger) CreateOrUpdateTagImageType(ctx context.Context, tit *model.TagImageType) error {
	start := time.Now()
	err := d.db.CreateOrUpdateTagImageType(ctx, tit)
	d.log(ctx, "CreateUpdateImgType", start, err, tit)
	return err
}

func (d *dbLogger) GetTagImageType(ctx context.Context, conds ...interface{}) (*model.TagImageType, error) {
	start := time.Now()
	result, err := d.db.GetTagImageType(ctx, conds...)
	d.log(ctx, "GetImgType", start, err, result)
	return result, err
}

func (d *dbLogger) GetTagImageTypes(ctx context.Context, conds ...interface{}) (*[]model.TagImageType, error) {
	start := time.Now()
	result, err := d.db.GetTagImageTypes(ctx, conds...)
	d.log(ctx, "GetImgTypes", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllTagImageTypes(ctx context.Context) (*[]model.TagImageType, error) {
	start := time.Now()
	result, err := d.db.GetAllTagImageTypes(ctx)
	d.log(ctx, "GetAllImgTypes", start, err, result)
	return result, err
}

// TagCustomCommand operations
func (d *dbLogger) CreateOrUpdateTagCustomCommand(ctx context.Context, command *model.TagCustomCommand) error {
	start := time.Now()
	err := d.db.CreateOrUpdateTagCustomCommand(ctx, command)
	d.log(ctx, "CreateUpdateCmd", start, err, command)
	return err
}

func (d *dbLogger) GetTagCustomCommand(ctx context.Context, conds ...interface{}) (*[]model.TagCustomCommand, error) {
	start := time.Now()
	result, err := d.db.GetTagCustomCommand(ctx, conds...)
	d.log(ctx, "GetCmd", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllTagCustomCommands(ctx context.Context) (*[]model.TagCustomCommand, error) {
	start := time.Now()
	result, err := d.db.GetAllTagCustomCommands(ctx)
	d.log(ctx, "GetAllCmds", start, err, result)
	return result, err
}

// Tag operations
func (d *dbLogger) CreateOrUpdateTag(ctx context.Context, tag *model.Tag) error {
	start := time.Now()
	err := d.db.CreateOrUpdateTag(ctx, tag)
	d.log(ctx, "CreateUpdateTag", start, err, tag)
	return err
}

func (d *dbLogger) UpdateTag(ctx context.Context, tag *model.Tag) error {
	start := time.Now()
	err := d.db.UpdateTag(ctx, tag)
	d.log(ctx, "UpdateTag", start, err, tag)
	return err
}

func (d *dbLogger) GetTag(ctx context.Context, conds ...interface{}) (*model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetTag(ctx, conds...)
	d.log(ctx, "GetTag", start, err, result)
	return result, err
}

func (d *dbLogger) GetTagsWithoutChildren(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetTagsWithoutChildren(ctx, conds...)
	d.log(ctx, "GetTagsNoChildren", start, err, result)
	return result, err
}

func (d *dbLogger) GetTags(ctx context.Context, conds ...interface{}) (*[]model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetTags(ctx, conds...)
	d.log(ctx, "GetTags", start, err, result)
	return result, err
}

func (d *dbLogger) GetAllTags(ctx context.Context) (*[]model.Tag, error) {
	start := time.Now()
	result, err := d.db.GetAllTags(ctx)
	d.log(ctx, "GetAllTags", start, err, result)
	return result, err
}

func (d *dbLogger) RemoveTagImageFromTag(ctx context.Context, tagId uint64, imageId uint64) error {
	start := time.Now()
	err := d.db.RemoveTagImageFromTag(ctx, tagId, imageId)
	d.log(ctx, "RemoveTagImage", start, err, fmt.Sprintf("tag=%d img=%d", tagId, imageId))
	return err
}

func (d *dbLogger) UpdateTagImage(ctx context.Context, image *model.TagImage) error {
	start := time.Now()
	err := d.db.UpdateTagImage(ctx, image)
	d.log(ctx, "UpdateTagImage", start, err, image)
	return err
}

func (d *dbLogger) GetTagsCount(ctx context.Context) (int64, error) {
	start := time.Now()
	result, err := d.db.GetTagsCount(ctx)
	d.log(ctx, "GetTagsCount", start, err, result)
	return result, err
}

// Task operations
func (d *dbLogger) CreateTask(ctx context.Context, task *model.Task) error {
	start := time.Now()
	err := d.db.CreateTask(ctx, task)
	d.log(ctx, "CreateTask", start, err, task)
	return err
}

func (d *dbLogger) UpdateTask(ctx context.Context, task *model.Task) error {
	start := time.Now()
	err := d.db.UpdateTask(ctx, task)
	d.log(ctx, "UpdateTask", start, err, task)
	return err
}

func (d *dbLogger) RemoveTasks(ctx context.Context, conds ...interface{}) error {
	start := time.Now()
	err := d.db.RemoveTasks(ctx, conds...)
	d.log(ctx, "RemoveTasks", start, err, nil)
	return err
}

func (d *dbLogger) TasksCount(ctx context.Context, query interface{}, conds ...interface{}) (int64, error) {
	start := time.Now()
	result, err := d.db.TasksCount(ctx, query, conds...)
	d.log(ctx, "TasksCount", start, err, result)
	return result, err
}

func (d *dbLogger) GetTasks(ctx context.Context, offset int, limit int) (*[]model.Task, error) {
	start := time.Now()
	result, err := d.db.GetTasks(ctx, offset, limit)
	d.log(ctx, "GetTasks", start, err, result)
	return result, err
}
