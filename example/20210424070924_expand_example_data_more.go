package main

import (
	migrations "github.com/dbreedt/hb_migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20210424070924_expand_example_data_more",
		up20210424070924ExpandExampleDataMore,
		down20210424070924ExpandExampleDataMore,
	)
}

func up20210424070924ExpandExampleDataMore(tx *pg.Tx) error {
	_, err := tx.Exec(`
		alter table example_data
		add column answer text;
	`)
	return err
}

func down20210424070924ExpandExampleDataMore(tx *pg.Tx) error {
	_, err := tx.Exec(`
		alter table example_data
		drop column answer;
	`)
	return err
}
