package sqlstore_test

import (
	"os"
	"testing"
)

var databaseURL string

func TestMain(t *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "host=localhost user=postgres password=guest dbname=restapi_test sslmode=disable"
	}

	os.Exit(t.Run())
}
