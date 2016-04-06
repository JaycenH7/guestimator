package db

import (
	"gopkg.in/pg.v4"
)

const dbname = "guestimator"

func Options(env string) *pg.Options {
	if env == "" {
		env = "development"
	}

	var opts pg.Options

	opts.Database = dbname + "_" + env
	opts.User = opts.Database
	opts.Password = opts.Database
	opts.SSL = false

	return &opts
}
