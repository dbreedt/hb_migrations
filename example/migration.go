package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"

	migrations "github.com/dbreedt/hb_migrations"
	"github.com/go-pg/pg/v10"
)

const usageText = `This program runs command on the database.
Supported commands are:
  - init - runs the specified intial migration as a batch on it's own.
  - migrate - runs all available migrations.
  - rollback - reverts the last batch of migration.
  - create - creates a migration file.
Usage:
  migrations -command=[command] -name=[name of migration] -template=[filename of the template]
`

func main() {
	flag.Usage = usage
	cmd := flag.String("command", "migrate", "Command that should be executed on the migration engine. Supported commands are: init, create, migrate and rollback")
	name := flag.String("name", "", "Name of the migration to be created")
	templateName := flag.String("template", "", "The filename of the template that should be used to create the migration")
	extra := flag.String("extra", "", "Extra parameters to pass to the command. Currently only migrate has an extra parameter called one-by-one which runs the migrations in batches of one")
	flag.Parse()

	var db *pg.DB

	ctx := context.Background()

	migrations.SetMigrationTableName("migrations_home")
	migrations.SetInitialMigration("000000000000_init")
	migrations.SetMigrationNameConvention(migrations.SnakeCase)

	template := ""

	if len(*templateName) > 0 {
		pwd, _ := os.Getwd()
		tmpPath := path.Join(pwd, "templates", *templateName)
		buf, err := ioutil.ReadFile(tmpPath)
		if err != nil {
			log.Fatalln("file not found", tmpPath)
		}

		template = string(buf)
	}

	if *cmd != "create" {
		db = pg.Connect(&pg.Options{
			Addr:     "127.0.0.1:5432",
			User:     "dev",
			Database: "example",
			Password: "12345",
		})
	}

	var err error

	switch *cmd {
	case "create":
		err = migrations.Run(ctx, db, *cmd, *name, template)
	case "migrate":
		err = migrations.Run(ctx, db, *cmd, *extra)
	default:
		err = migrations.Run(ctx, db, *cmd)
	}

	if err != nil {
		log.Fatalln("run failed", err)
	}
}

func usage() {
	log.Printf(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}
