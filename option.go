package miniofs

import (
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

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
