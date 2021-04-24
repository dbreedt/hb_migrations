package main

import (
	migrations "github.com/dbreedt/hb_migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20210424065345_expand_example_data",
		up20210424065345ExpandExampleData,
		down20210424065345ExpandExampleData,
	)
}

func up20210424065345ExpandExampleData(tx *pg.Tx) error {
	_, err := tx.Exec(`
		alter table example_data
		add column question text;
	`)
	return err
}

func down20210424065345ExpandExampleData(tx *pg.Tx) error {
	_, err := tx.Exec(`
		alter table example_data
		drop column question;
	`)
	return err
}
