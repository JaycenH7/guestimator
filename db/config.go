package db

import (
	"gopkg.in/pg.v4"
	"os"
)

const dbname = "guestimator"

// dbAddr - database address
var dbAddr = os.Getenv("DB_ADDR")

// DevDB - development database
var DevDB = pg.Connect(Options("dev", dbAddr))

// Options - this is a test
func Options(env string, addr string) *pg.Options {
	if env == "" {
		env = "dev"
	}
	if addr == "" {
		addr = "localhost"
	}

	var opts pg.Options

	opts.Database = dbname + "_" + env
	opts.User = opts.Database
	opts.Password = opts.Database
	opts.SSL = false
	ops.Addr = addr

	return &opts
}
