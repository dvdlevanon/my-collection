package fswatch

// func newFsSync() *fsSync {

// }

// type fsSync struct {
// 	stales directorytree
// }

// type DirectoryItemsGetterImpl struct {
// 	db *db.Database
// }

// func (g DirectoryItemsGetterImpl) GetBelongingItems(path string) (*[]model.Item, error) {
// 	name := filepath.Base(path)
// 	title := directories.DirectoryNameToTag(name)
// 	tag, err := tags.GetChildTag(g.db, directories.DIRECTORIES_TAG_ID, title)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			empty := make([]model.Item, 0)
// 			return &empty, nil
// 		}
// 		return nil, err
// 	}

// 	items, err := tags.GetItems(g.db, tag)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return items, nil
// }

// func TestTemp(t *testing.T) {
// 	rootfs, _ := BuildFromPath("/mnt/usb1", func(path string) bool {
// 		return utils.IsVideo(true, path)
// 	})
// 	db, err := db.New("/mnt/usb1", "test.sqlite")
// 	assert.NoError(t, err)
// 	assert.NoError(t, err)
// 	rootdb, _ := BuildFromDb(db, DirectoryItemsGetterImpl{db: db})
// 	diff := Compare(rootfs, rootdb)

// 	os.WriteFile("/tmp/db", []byte(rootdb.String(0)), 0755)
// 	os.WriteFile("/tmp/fs", []byte(rootfs.String(0)), 0755)
// 	os.WriteFile("/tmp/diff", []byte(diff.String()), 0755)
// }
