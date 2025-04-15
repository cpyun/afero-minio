package miniofs

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/spf13/afero"
)

const (
	defaultFileMode = 0o755
)

// Fs is a Fs implementation that uses functions provided by google cloud storage
type Fs struct {
	ctx       context.Context
	client    *minio.Client
	bucket    string
	separator string
}

func NewMinioFs(ctx context.Context, dsn string) afero.Fs {
	url, _ := url.Parse(dsn)
	minioOpts, _ := ParseURL(dsn)

	client, _ := minio.New(url.Host, minioOpts)
	return NewFs(ctx, client, url.Path[1:])
}

func NewFs(ctx context.Context, client *minio.Client, bucket string) *Fs {
	return &Fs{
		ctx:       ctx,
		client:    client,
		bucket:    bucket,
		separator: "/",
	}
}

// normSeparators will normalize all "\\" and "/" to the provided separator
func (fs *Fs) normSeparators(s string) string {
	return strings.Replace(strings.Replace(s, "\\", fs.separator, -1), "/", fs.separator, -1)
}

func (fs *Fs) ensureNoLeadingSeparator(s string) string {
	if len(s) > 0 && strings.HasPrefix(s, fs.separator) {
		s = s[len(fs.separator):]
	}

	return s
}

//func (fs *Fs) getObj(name string) (*minio.Object, error) {
//	bucketName, path := fs.splitName(name)
//	getObjectOptions := minio.GetObjectOptions{}
//
//	return fs.client.GetObject(fs.ctx, bucketName, path, getObjectOptions)
//}

func (fs *Fs) Name() string { return "MinioFs" }

func (fs *Fs) Create(name string) (afero.File, error) {
	return fs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
}

func (fs *Fs) Mkdir(name string, _ os.FileMode) error {
	return ErrNotSupported
}

func (fs *Fs) MkdirAll(_ string, _ os.FileMode) error {
	return ErrNotSupported
}

func (fs *Fs) Open(name string) (afero.File, error) {
	return fs.OpenFile(name, os.O_RDONLY, 0)
}

func (fs *Fs) OpenFile(name string, flag int, fileMode os.FileMode) (afero.File, error) {
	var err error
	if flag&os.O_APPEND != 0 {
		return nil, errors.New("appending files will lead to trouble")
	}

	name = fs.ensureNoLeadingSeparator(fs.normSeparators(name))
	file := NewMinioFile(fs.ctx, fs, flag, fileMode, name)
	//
	if flag&os.O_CREATE != 0 {
		_, err = file.WriteString("")
	}

	return file, err
}

func (fs *Fs) Remove(name string) error {
	name = fs.ensureNoLeadingSeparator(fs.normSeparators(name))
	return fs.client.RemoveObject(fs.ctx, fs.bucket, name, minio.RemoveObjectOptions{
		GovernanceBypass: true,
	})
}

func (fs *Fs) RemoveAll(path string) error {
	path = fs.ensureNoLeadingSeparator(fs.normSeparators(path))

	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		opts := minio.ListObjectsOptions{Prefix: path, Recursive: true}
		for object := range fs.client.ListObjects(fs.ctx, fs.bucket, opts) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object
		}
	}()

	errorCh := fs.client.RemoveObjects(fs.ctx, fs.bucket, objectsCh, minio.RemoveObjectsOptions{})
	for e := range errorCh {
		return errors.New("Failed to remove " + e.ObjectName + ", error: " + e.Err.Error())
	}

	return nil
}

func (fs *Fs) Rename(oldName, newName string) error {
	if oldName == newName {
		return nil
	}

	oldName = fs.ensureNoLeadingSeparator(fs.normSeparators(oldName))
	newName = fs.ensureNoLeadingSeparator(fs.normSeparators(newName))

	// Source object
	src := minio.CopySrcOptions{
		Bucket: fs.bucket,
		Object: oldName,
	}
	dst := minio.CopyDestOptions{
		Bucket: fs.bucket,
		Object: newName,
	}
	_, err := fs.client.CopyObject(fs.ctx, dst, src)
	if err != nil {
		return err
	}

	return fs.Remove(oldName)
}

func (fs *Fs) Stat(name string) (os.FileInfo, error) {
	name = fs.ensureNoLeadingSeparator(fs.normSeparators(name))

	file := NewMinioFile(fs.ctx, fs, os.O_RDWR, defaultFileMode, name)
	return file.Stat()
}

func (fs *Fs) Chmod(_ string, _ os.FileMode) error {
	return errors.New("method Chmod is not implemented in Minio")
}

func (fs *Fs) Chtimes(_ string, _, _ time.Time) error {
	return errors.New("method Chtimes is not implemented. Create, Delete, Updated times are read only fields in Minio and set implicitly")
}

func (fs *Fs) Chown(_ string, _, _ int) error {
	return errors.New("method Chown is not implemented for Minio")
}
