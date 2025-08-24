package fssync

import (
	"my-collection/server/pkg/bl/directories"
	"my-collection/server/pkg/directorytree"
	"my-collection/server/pkg/model"
	"my-collection/server/pkg/relativasor"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"k8s.io/utils/pointer"
)

func TestAddMissingDirectoryTags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dr := model.NewMockDirectoryReader(ctrl)

	dr.EXPECT().GetAllDirectories().Return(&[]model.Directory{
		{Path: "dir1"},
		{Path: "dir-2"},
		{Path: "dir_3"},
		{Path: "some/inner/path"},
		{Path: "tag/exists"},
		{Path: "excluded", Excluded: pointer.Bool(true)},
	}, nil)

	trw := model.NewMockTagReaderWriter(ctrl)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "dir1", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "dir-2", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "dir_3", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "some/inner/path", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}).Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Title: "tag/exists", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}, nil)
	errs := addMissingDirectoryTags(dr, trw)
	assert.Equal(t, 0, len(errs))
}

func TestAddMissingDirs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	drw := model.NewMockDirectoryReaderWriter(ctrl)

	relativasor.Init("/root/dir")
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "path1", Excluded: pointer.Bool(true)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "some/path - to dir", Excluded: pointer.Bool(true)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(&model.Directory{Path: "/absolute/exists/path3"}, nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "/absolute/error/inner/path4", Excluded: pointer.Bool(true)}).Return(errors.Errorf("test"))
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: model.ROOT_DIRECTORY_PATH, Excluded: pointer.Bool(true)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(&model.Directory{Path: model.ROOT_DIRECTORY_PATH}, nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "some/relative/path", Excluded: pointer.Bool(true)}).Return(nil)

	errs := addMissingDirs(drw, []directorytree.Change{
		{Path1: "path1"},
		{Path1: "some/path - to dir"},
		{Path1: "/absolute/exists/path3"},
		{Path1: "/absolute/error/inner/path4"},
		{Path1: ""},
		{Path1: ""},
		{Path1: "/root/dir/some/relative/path"},
	})
	assert.Equal(t, 1, len(errs))
}

func TestAddNewFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	iw := model.NewMockItemWriter(ctrl)
	digs := model.NewMockDirectoryItemsGetterSetter(ctrl)
	dctg := model.NewMockDirectoryConcreteTagsGetter(ctrl)
	fmg := model.NewMockFileMetadataGetter(ctrl)

	relativasor.Init("/root/dir")

	tags := []*model.Tag{{Id: 3, Title: "old"}}
	digs.EXPECT().GetBelongingItem("new", "file").Return(&model.Item{Title: "file", Origin: "new", Url: "new/file", LastModified: 4234234, Tags: tags}, nil)
	concreteTags1 := []*model.Tag{{Id: 1, Title: "concrete1"}, {Id: 2, Title: "concrete2"}}
	dctg.EXPECT().GetConcreteTags("new").Return(concreteTags1, nil)
	iw.EXPECT().UpdateItem(&model.Item{Title: "file", Origin: "new", Url: "new/file", LastModified: 4234234, Tags: append(tags, concreteTags1...)})

	digs.EXPECT().GetBelongingItem("some/deep/deep/deep", "file").Return(nil, nil)
	concreteTags2 := []*model.Tag{{Title: "concrete1"}, {Title: "concrete2"}}
	dctg.EXPECT().GetConcreteTags("some/deep/deep/deep").Return(concreteTags2, nil)
	fmg.EXPECT().GetFileMetadata("/root/dir/some/deep/deep/deep/file").Return(int64(7657657), int64(0), nil)
	digs.EXPECT().AddBelongingItem(&model.Item{Title: "file", Origin: "some/deep/deep/deep", Url: "some/deep/deep/deep/file", LastModified: 7657657, Tags: concreteTags2}).Return(nil)

	digs.EXPECT().GetBelongingItem("/absolute", "path").Return(nil, nil)
	dctg.EXPECT().GetConcreteTags("/absolute").Return([]*model.Tag{}, nil)
	fmg.EXPECT().GetFileMetadata("/absolute/path").Return(int64(7567657), int64(0), nil)
	digs.EXPECT().AddBelongingItem(&model.Item{Title: "path", Origin: "/absolute", Url: "/absolute/path", LastModified: 7567657, Tags: make([]*model.Tag, 0)}).Return(nil)

	digs.EXPECT().GetBelongingItem("some", "file").Return(nil, nil)
	dctg.EXPECT().GetConcreteTags("some").Return([]*model.Tag{}, nil)
	fmg.EXPECT().GetFileMetadata("/root/dir/some/file").Return(int64(9876532), int64(0), nil)
	digs.EXPECT().AddBelongingItem(&model.Item{Title: "file", Origin: "some", Url: "some/file", LastModified: 9876532, Tags: make([]*model.Tag, 0)}).Return(nil)

	errs := addNewFiles(iw, digs, dctg, fmg, []directorytree.Change{
		{Path1: "new/file"},
		{Path1: "some/deep/deep/deep/file"},
		{Path1: "/absolute/path"},
		{Path1: "/root/dir/some/file"},
	})

	assert.Equal(t, 0, len(errs))
}

func TestRemoveStaleDirs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trw := model.NewMockTagReaderWriter(ctrl)
	dw := model.NewMockDirectoryWriter(ctrl)

	relativasor.Init("/root/dir")
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Id: 1, Title: "tag1"}, nil)
	trw.EXPECT().RemoveTag(uint64(1)).Return(nil)
	dw.EXPECT().RemoveDirectory("/existing/dir").Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, nil)
	dw.EXPECT().RemoveDirectory("/not/exists/tag").Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Id: 2, Title: "relative"}, nil)
	trw.EXPECT().RemoveTag(uint64(2)).Return(nil)
	dw.EXPECT().RemoveDirectory("relative/path").Return(nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(&model.Tag{Id: 3, Title: model.ROOT_DIRECTORY_PATH}, nil)
	trw.EXPECT().RemoveTag(uint64(3)).Return(nil)
	dw.EXPECT().RemoveDirectory(model.ROOT_DIRECTORY_PATH).Return(nil)

	errs := removeStaleDirs(trw, dw, []string{
		"/existing/dir",
		"/not/exists/tag",
		"/root/dir/relative/path",
		"",
	})
	assert.Equal(t, 0, len(errs))
}

func TestRemoveStaleItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dig := model.NewMockDirectoryItemsGetter(ctrl)
	iw := model.NewMockItemWriter(ctrl)

	relativasor.Init("/root/dir")
	dig.EXPECT().GetBelongingItem("/some", "not-exists").Return(nil, nil)
	dig.EXPECT().GetBelongingItem("/some/inner", "file").Return(&model.Item{Id: 1}, nil)
	iw.EXPECT().RemoveItem(uint64(1)).Return(nil)
	dig.EXPECT().GetBelongingItem("relative", "file").Return(&model.Item{Id: 2}, nil)
	iw.EXPECT().RemoveItem(uint64(2)).Return(nil)

	errs := removeStaleItems(dig, iw, []string{
		"/some/not-exists",
		"/some/inner/file",
		"/root/dir/relative/file",
	})
	assert.Equal(t, 0, len(errs))
}

func TestRenameFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	trw := model.NewMockTagReaderWriter(ctrl)
	irw := model.NewMockItemReaderWriter(ctrl)
	drw := model.NewMockDirectoryReaderWriter(ctrl)

	relativasor.Init("/root/dir")
	drw.EXPECT().GetDirectory(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	drw.EXPECT().CreateOrUpdateDirectory(&model.Directory{Path: "/new", Excluded: pointer.Bool(false)}).Return(nil)
	drw.EXPECT().GetDirectory(gomock.Any()).Return(&model.Directory{Path: "/new", Excluded: pointer.Bool(false)}, nil)
	trw.EXPECT().GetTag(gomock.Any()).Return(nil, gorm.ErrRecordNotFound)
	trw.EXPECT().CreateOrUpdateTag(&model.Tag{Title: "/new", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())})
	oldTag := &model.Tag{Id: 4321, Title: "/", ParentID: pointer.Uint64(directories.GetDirectoriesTagId()), Items: []*model.Item{{Id: 1}}}
	trw.EXPECT().GetTag(gomock.Any()).Return(oldTag, nil)
	irw.EXPECT().GetItems(gomock.Any()).Return(&[]model.Item{{Id: 1, Title: "file1"}}, nil)
	irw.EXPECT().UpdateItem(&model.Item{Id: 1, Title: "file1", Origin: "/new", Url: "/new/file1"}).Return(nil)
	newTag := &model.Tag{Title: "new", ParentID: pointer.Uint64(directories.GetDirectoriesTagId())}
	trw.EXPECT().GetTag(gomock.Any()).Return(newTag, nil)
	irw.EXPECT().CreateOrUpdateItem(&model.Item{Id: 1, Title: "file1", Origin: "/new", Url: "/new/file1", Tags: []*model.Tag{newTag}})
	trw.EXPECT().GetTag(gomock.Any()).Return(oldTag, nil)
	irw.EXPECT().RemoveTagFromItem(uint64(1), uint64(4321)).Return(nil)

	errs := renameFiles(trw, drw, irw, []directorytree.Change{
		{Path1: "/file1", Path2: "/new/file1"},
	})
	assert.Equal(t, 0, len(errs))
}
