package miniofs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
)

//const (
//	maxWriteSize int = 1e4
//)

type readerAtCloser interface {
	io.ReadCloser
	io.ReaderAt
}

type minioFileResource struct {
	ctx context.Context
	fs  *Fs

	name     string
	fileMode os.FileMode

	currentIoSize int64
	offset        int64
	reader        readerAtCloser
	writer        io.WriteCloser

	closed bool
}

func (o *minioFileResource) Close() error {
	o.closed = true
	// TODO rawGcsObjectsMap ?
	return o.maybeCloseIo()
}

func (o *minioFileResource) maybeCloseIo() error {
	if err := o.maybeCloseReader(); err != nil {
		return fmt.Errorf("error closing reader: %v", err)
	}
	if err := o.maybeCloseWriter(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	return nil
}

func (o *minioFileResource) maybeCloseReader() error {
	if o.reader == nil {
		return nil
	}
	if err := o.reader.Close(); err != nil {
		return err
	}
	o.reader = nil
	return nil
}

func (o *minioFileResource) maybeCloseWriter() error {
	if o.writer == nil {
		return nil
	}

	// In cases of partial writes (e.g. to the middle of a file stream), we need to
	// append any remaining data from the original file before we close the reader (and
	// commit the results.)
	// For small writes it can be more efficient
	// to keep the original reader but that is for another iteration
	//if o.currentIoSize > o.offset {
	//
	//}

	if err := o.writer.Close(); err != nil {
		return err
	}
	o.writer = nil
	return nil
}

func (o *minioFileResource) ReadAt(p []byte, off int64) (n int, err error) {
	if cap(p) == 0 {
		return 0, nil
	}

	// Assume that if the reader is open; it is at the correct offset
	// a good performance assumption that we must ensure holds
	if off == o.offset && o.reader != nil {
		n, err = o.reader.ReadAt(p, off)
		o.offset += int64(n)
		return n, err
	}

	// If any writers have written anything; commit it first so we can read it back.
	if err = o.maybeCloseIo(); err != nil {
		return 0, err
	}

	opts := minio.GetObjectOptions{}
	r, err := o.fs.client.GetObject(o.ctx, o.fs.bucket, o.name, opts)
	if err != nil {
		return 0, err
	}
	o.reader = r
	o.offset = off

	read, err := r.ReadAt(p, off)
	o.offset += int64(read)
	return read, err
}

func (o *minioFileResource) WriteAt(b []byte, off int64) (n int, err error) {
	// If the writer is opened and at the correct offset we're good!
	if off == o.offset && o.writer != nil {
		n, err = o.writer.Write(b)
		o.offset += int64(n)
		return n, err
	}

	// Ensure readers must be re-opened and that if a writer is active at another
	// offset it is first committed before we do a "seek" below
	if err = o.maybeCloseIo(); err != nil {
		return 0, err
	}

	// WriteAt to a non existing file
	if off > o.currentIoSize {
		return 0, ErrOutOfRange
	}
	o.offset = off
	//o.writer =

	// byt 写入 buffer
	buffer := bytes.NewReader(b)
	// 写入 minio
	opts := minio.PutObjectOptions{
		ContentType: http.DetectContentType(b),
	}
	if off > 0 {
		opts.PartSize = uint64(off)
		opts.NumThreads = 8
		opts.ConcurrentStreamParts = false
		opts.DisableMultipart = true
	}
	_, err = o.fs.client.PutObject(o.ctx, o.fs.bucket, o.name, buffer, buffer.Size(), opts)
	if err != nil {
		return 0, err
	}

	o.offset += int64(buffer.Len())
	return buffer.Len(), nil
}

func (o *minioFileResource) Truncate(_ int64) error {
	return ErrNotSupported
}
