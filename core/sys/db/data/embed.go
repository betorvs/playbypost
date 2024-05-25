package data

import (
	"embed"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var fs embed.FS

func Assets() (source.Driver, error) {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return d, err
	}
	return d, nil
}
