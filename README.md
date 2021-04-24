# hb migrations - A Better migration engine for [go-pg/pg](https://github.com/go-pg/pg)

### Basic terminology:
- **Batch**:
  - Is a set of migration go files that were run in a single transaction/batch.
- **Migrate**:
  - Compare the migration go files against the contents of the `migrationTableName` and runs the missing migration go files' `up` function in a batch.
- **Rollback**:
  - Get the last batch's migrations from the `migrationTableName` table and run each ones `down` function in reverse insert order order.
- **Corrupt Migration Table**:
  - An entry exists in the `migrationTableName` that does not exists in any of the migration go files.
- **One-by-one**:
  - Special migration run mode that runs each migration go file in their own batch
- **Init**:
  - If you already have an existing set of tables and things that you want "seed" a test or other db with, you create a special migration called `000000000000_init.go` see the one in the example for more details [[example/000000000000_init.go](https://github.com/dbreedt/hb_migrations/blob/master/example/000000000000_init.go)]
- **Template**:
  - If you don't like the default template created for a migration you can create your own, the only gotcha is that the package expects it to be in a `templates` folder
## Basic Commands
- init
  - runs the specified initial migration as a batch on it's own.
- migrate
  - runs all available migrations that have not been run inside a batch
- rollback
  - reverts the last batch of migrations.
- create **name**
  - creates a migration file using the name provided.

## Usage
The best way to understand how this library works is to look at the [example/migration.go](https://github.com/dbreedt/hb_migrations/blob/master/example/migration.go)
It creates a cli to make the running, rollback and creation of migrations using the library easier.

### Compiling
```bash
$> go build -o migration *.go
```

### Run the init migration
```bash
$> ./migration -command=init
Batch 1 run: 1 migrations
Completed 000000000000_init
```

#### DB state
```psql
example=# \dt
            List of relations
 Schema |      Name       | Type  | Owner
--------+-----------------+-------+-------
 public | example_data    | table | dev
 public | migrations_home | table | dev
(2 rows)

example=# \d example_data
                            Table "public.example_data"
 Column |  Type  | Collation | Nullable |                 Default
--------+--------+-----------+----------+------------------------------------------
 id     | bigint |           | not null | nextval('example_data_id_seq'::regclass)
 name   | text   |           |          |

example=# select * from migrations_home;
 id |       name        | batch |        migration_time
----+-------------------+-------+-------------------------------
  1 | 000000000000_init |     1 | 2021-04-24 06:52:40.820578+02
(1 row)

```

### Create a new migration
```bash
$> ./migration -command create -name 'expand_example_data'
Created migration /data/proj/hb_migrations/example/20210424065345_expand_example_data.go
```

### Migrate
```bash
-> ./migration -command migrate
Batch 2 run: 1 migrations
Completed 20210424065345_expand_example_data
```

#### Db State
```psql
example=# \d example_data
                             Table "public.example_data"
  Column  |  Type  | Collation | Nullable |                 Default
----------+--------+-----------+----------+------------------------------------------
 id       | bigint |           | not null | nextval('example_data_id_seq'::regclass)
 name     | text   |           |          |
 question | text   |           |          |

example=# select * from migrations_home;
 id |                name                | batch |        migration_time
----+------------------------------------+-------+-------------------------------
  1 | 000000000000_init                  |     1 | 2021-04-24 07:07:40.10731+02
  2 | 20210424065345_expand_example_data |     2 | 2021-04-24 07:07:44.959022+02
(2 rows)
```

### Rollback
```bash
$> ./migration -command rollback
Batch 2 rollback: 1 migrations
Rolledback 20210424065345_expand_example_data
```

#### Db State
```psql

example=# \d example_data
                            Table "public.example_data"
 Column |  Type  | Collation | Nullable |                 Default
--------+--------+-----------+----------+------------------------------------------
 id     | bigint |           | not null | nextval('example_data_id_seq'::regclass)
 name   | text   |           |          |

example=# select * from migrations_home;
 id |       name        | batch |        migration_time
----+-------------------+-------+------------------------------
  1 | 000000000000_init |     1 | 2021-04-24 07:07:40.10731+02
(1 row)
```

### One-By-One Migrations
```bash
./migration -command migrate -extra one-by-one
Batch 2 run: 1 migration - 20210424065345_expand_example_data
Completed 20210424065345_expand_example_data
Batch 3 run: 1 migration - 20210424070924_expand_example_data_more
Completed 20210424070924_expand_example_data_more
```

#### Db State
```psql
example=# select * from migrations_home;
 id |                  name                   | batch |        migration_time
----+-----------------------------------------+-------+-------------------------------
  1 | 000000000000_init                       |     1 | 2021-04-24 07:07:40.10731+02
  4 | 20210424065345_expand_example_data      |     2 | 2021-04-24 07:10:59.790563+02
  5 | 20210424070924_expand_example_data_more |     3 | 2021-04-24 07:10:59.792458+02
(3 rows)
```

### Create a migration using a custom template
```bash
$> ./migration -command create -name 'expand_example_data_again' -template custom
Created migration /data/proj/hb_migrations/example/20210424072843_expand_example_data_again.go
```

### What do errors look ike
```bash
$> ./migration -command rollback
Batch 2 rollback: 1 migrations
2021/04/24 07:01:28 run failed 20210424065345_expand_example_data failed to rollback: ERROR #42703 column "question_text" of relation "example_data" does not exist
```

## Notes on generated file names
```bash
$> ./migrations -command create -name new_index
```
Creates a file in the current folder called `20181031230738_new_index.go` with the following contents:

```golang
package main

import (
	migrations "github.com/dbreedt/hb_migrations"
	"github.com/go-pg/pg/v10"
)

func init() {
	migrations.Register(
		"20181031230738_new_index",
		up20181031230738NewIndex,
		down20181031230738NewIndex,
	)
}

func up20181031230738NewIndex(tx *pg.Tx) error {
	_, err := tx.Exec(``)
	return err
}

func down20181031230738NewIndex(tx *pg.Tx) error {
	_, err := tx.Exec(``)
	return err
}
```

Forward migration sql commands go in `up*` and Rollback migrations sql commands go in `down*`

