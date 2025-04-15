package miniofs

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

const (
	minioDsn = "https://Q3AM3UQ867SPQQA43P2F:zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG@play.min.io/my-bucket?region=us-east-1"
)

func TestNew(t *testing.T) {
	appFs := NewMinioFs(context.Background(), minioDsn)

	t.Run("create", func(t *testing.T) {
		name := uuid.New().String() + ".create"
		f, err := appFs.Create(name)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()
	})

	t.Run("open", func(t *testing.T) {
		name := uuid.New().String() + ".open"
		f, err := appFs.Open(name)
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		if _, err := f.WriteString("open_test"); err != nil {
			t.Error(err)
		}
	})

}
