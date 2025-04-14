package miniofs

import (
	"context"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/afero"
)

type MinioFs struct {
	source *Fs
}

func NewMinio(ctx context.Context, dsn string) afero.Fs {
	url, _ := url.Parse(dsn)
	minioOpts, _ := ParseURL(dsn)

	client, _ := minio.New(url.Host, minioOpts)
	fs := NewFs(ctx, client, url.Path[1:])

	return &MinioFs{
		source: fs,
	}
}

func (fs *MinioFs) Name() string {
	return fs.source.Name()
}

func (fs *MinioFs) Create(name string) (afero.File, error) {
	return fs.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
}

func (fs *MinioFs) Mkdir(name string, perm os.FileMode) error {
	return fs.source.Mkdir(name, perm)
}

func (fs *MinioFs) MkdirAll(path string, perm os.FileMode) error {
	return fs.source.MkdirAll(path, perm)
}

func (fs *MinioFs) Open(name string) (afero.File, error) {
	return fs.OpenFile(name, os.O_RDONLY, 0)
}

func (fs *MinioFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return fs.source.OpenFile(name, flag, perm)
}

func (fs *MinioFs) Remove(name string) error {
	return fs.source.Remove(name)
}

func (fs *MinioFs) RemoveAll(path string) error {
	return fs.source.RemoveAll(path)
}

func (fs *MinioFs) Rename(oldname, newname string) error {
	return fs.source.Rename(oldname, newname)
}

func (fs *MinioFs) Stat(name string) (os.FileInfo, error) {
	return fs.source.Stat(name)
}

func (fs *MinioFs) Chmod(name string, mode os.FileMode) error {
	return fs.source.Chmod(name, mode)
}

func (fs *MinioFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return fs.source.Chtimes(name, atime, mtime)
}

func (fs *MinioFs) Chown(name string, uid, gid int) error {
	return fs.source.Chown(name, uid, gid)
}

func ParseURL(minioURL string) (*minio.Options, error) {
	u, err := url.Parse(minioURL)
	if err != nil {
		return nil, err
	}

	o := &minio.Options{
		Region: "us-east-1",
	}
	// credentials
	username, password := getUserPassword(u)
	token := u.Query().Get("token")
	o.Creds = credentials.NewStaticV4(username, password, token)
	//
	if u.Scheme == "https" {
		o.Secure = true
	}

	//
	if u.Query().Has("region") {
		o.Region = u.Query().Get("region")
	}

	return o, nil
}

func getUserPassword(u *url.URL) (string, string) {
	var user, password string
	if u.User != nil {
		user = u.User.Username()
		if p, ok := u.User.Password(); ok {
			password = p
		}
	}
	return user, password
}

func getHostPortWithDefaults(u *url.URL) (string, string) {
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "9000"
	}
	return host, port
}
