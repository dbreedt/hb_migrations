package main

import (
	migrations "github.com/dbreedt/hb_migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"000000000000_init",
		up000000000000Init,
		down000000000000Init,
	)
}

func up000000000000Init(tx *pg.Tx) error {
	_, err := tx.Exec(`
		create table example_data(
			id bigserial,
			name text
		)
	`)
	return err
}

func down000000000000Init(tx *pg.Tx) error {
	_, err := tx.Exec(`
		drop table if exists example_data;
	`)
	return err
}
