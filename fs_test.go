package miniofs

import (
	"context"
	"testing"
)

const (
	minioDsn = "https://Q3AM3UQ867SPQQA43P2F:zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG@play.min.io/12345?region=us-east-1"
)

func TestNew(t *testing.T) {
	appFs := NewMinio(context.Background(), minioDsn)

	f, err := appFs.Open("test")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	if _, err := f.WriteString("test"); err != nil {
		t.Error(err)
	}
}
