# minio sdk for afero
[![Go Report Card](https://goreportcard.com/badge/github.com/cpyun/afero-minio)](https://goreportcard.com/report/github.com/cpyun/afero-minio)
[![GoDoc](https://godoc.org/github.com/cpyun/afero-minio?status.svg)](https://godoc.org/github.com/cpyun/afero-minio)

## About
It provides an afero filesystem implementation of a MinIO.

This was created to provide a backend to the MinIO server but can definitely be used in any other code.

I'm very opened to any improvement through issues or pull-request that might lead to a better implementation or even better testing.

## Key points
- Download & upload file streaming
- Very carefully linted

## Known limitations
- File appending / seeking for write is not supported because MinIO doesn't support it, it could be simulated by rewriting entire files.
- Chtimes is not supported because MinIO doesn't support it, it could be simulated through metadata.
- Chmod support is very limited

## How to use 
```go
import (
    "context"
    
    "github.com/minio/minio-go/v7"
	"github.com/cpyun/afero-minio"
)

func main() {
    minioClient := minio.New(endpoint, &opts)

	// Initialize the MinIOFS
    fs := miniofs.NewMinioFs(context.Background(), minioClient)
    
	// And use it
	file, _ := fs.Open("text.txt")
	defer file.Close()
	file.WriteString("Hello world.")
}
```

## Thanks
- [spf13/afero](https://github.com/spf13/afero)
- [fclairamb/afero-s3](https://github.com/fclairamb/afero-s3)

## License
Afero is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/cpyun/afero-minio/blob/master/LICENSE)