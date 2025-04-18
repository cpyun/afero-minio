package miniofs

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

const (
	folderSize = 42
)

type FileInfo struct {
	eTag     string
	name     string
	size     int64
	updated  time.Time
	isDir    bool
	fileMode os.FileMode
}

func newFileInfoFromAttrs(obj minio.ObjectInfo, fileMode os.FileMode) *FileInfo {
	res := &FileInfo{
		eTag:     obj.ETag,
		name:     obj.Key,
		size:     obj.Size,
		updated:  obj.LastModified,
		isDir:    false,
		fileMode: fileMode,
	}

	if res.name == "" {
		// deals with them at the moment
		//res.name = "folder"
		res.size = folderSize
		res.isDir = true
	}

	return res
}

func (fi *FileInfo) Name() string {
	return filepath.Base(filepath.FromSlash(fi.name))
}

func (fi *FileInfo) Size() int64 {
	return fi.size
}

func (fi *FileInfo) Mode() os.FileMode {
	if fi.IsDir() {
		return os.ModeDir | fi.fileMode
	}
	return fi.fileMode
}

func (fi *FileInfo) ModTime() time.Time {
	return fi.updated
}

func (fi *FileInfo) IsDir() bool {
	return fi.isDir
}

func (fi *FileInfo) Sys() interface{} {
	return nil
}

type ByName []*FileInfo

func (a ByName) Len() int { return len(a) }
func (a ByName) Swap(i, j int) {
	a[i].name, a[j].name = a[j].name, a[i].name
	a[i].size, a[j].size = a[j].size, a[i].size
	a[i].updated, a[j].updated = a[j].updated, a[i].updated
	a[i].isDir, a[j].isDir = a[j].isDir, a[i].isDir
}
func (a ByName) Less(i, j int) bool { return strings.Compare(a[i].Name(), a[j].Name()) == -1 }
