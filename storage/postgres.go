package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresDBConfig struct {
	DbURL string
}

type postgresStorage struct {
	config PostgresDBConfig
	db     *sql.DB
}

// creates new instance of postgres DB
func NewPostgresStorage(config PostgresDBConfig) (postgres *postgresStorage, err error) {
	// TODO should connect and initialize as well.
	postgres = &postgresStorage{config: config}
	// tyr to connect to givenDB.
	err = postgres.connect()
	if err != nil {
		log.Printf("Unable to connect to DB : %s, error : %v", config.DbURL, err)
		return nil, err
	}
	// lets initialize the DB as needed.
	err = postgres.initialize()
	if err != nil {
		log.Printf("Couldn't initialize the DB, error : %v", err)
		return nil, err
	}
	return postgres, nil
}

// Connect creates a database connection to the given config URL, and assigns to the Storage fields `db`.
func (postgres *postgresStorage) connect() (err error) {
	db, err := sql.Open("postgres", postgres.config.DbURL)
	if err != nil {
		return err
	}
	postgres.db = db
	return nil
}

func (postgres *postgresStorage) initialize() (err error) {
	stmt := `
	CREATE TABLE IF NOT EXISTS task_store (
		id SERIAL NOT NULL PRIMARY KEY,
		name text,
		params text,
		duration text,
		last_run text,
		next_run text,
		is_recurring text,
		hash text
	);
	`
	_, err = postgres.db.Exec(stmt)
	if err != nil {
		log.Printf("Error while initializing: %q - %+v", stmt, err)
		return
	}
	return
}

func (postgres *postgresStorage) Close() error {
	return postgres.db.Close()
}

func (postgres *postgresStorage) Add(task TaskAttributes) error {
	// should add a task to the database `task_store` table
	var count int
	rows, err := postgres.db.Query("SELECT count(*) FROM task_store WHERE hash=($1) ;", task.Hash)
	defer rows.Close()
	if err == nil {
		rows.Next()
		_ = rows.Scan(&count)
	}

	if count == 0 {
		return postgres.insert(task)
	}
	return nil
}

func (postgres *postgresStorage) Fetch() ([]TaskAttributes, error) {
	// read all the rows task_store table.
	rows, err := postgres.db.Query(`
        SELECT name, params, duration, last_run, next_run, is_recurring
        FROM task_store ;`)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var tasks []TaskAttributes

	for rows.Next() {
		// var task TaskAttributes
		task := TaskAttributes{}
		err = rows.Scan(&task.Name, &task.Params, &task.Duration, &task.LastRun, &task.NextRun, &task.IsRecurring)
		if err != nil {
			return []TaskAttributes{}, err
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tasks, nil
}

func (postgres *postgresStorage) Remove(task TaskAttributes) error {
	// should delete the entry from `task_stor` table.
	stmt, err := postgres.db.Prepare(`DELETE FROM task_store WHERE hash=($1) ;`)

	if err != nil {
		return fmt.Errorf("Error while pareparing delete task statement: %s+v", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Hash,
	)
	if err != nil {
		return fmt.Errorf("Error while deleting task: %+v", err)
	}

	return nil
}

func (postgres *postgresStorage) insert(task TaskAttributes) (err error) {
	stmt, err := postgres.db.Prepare(`
        INSERT INTO task_store(name, params, duration, last_run, next_run, is_recurring, hash)
        VALUES(($1), ($2), ($3), ($4), ($5), ($6), ($7));`)

	if err != nil {
		return fmt.Errorf("Error while pareparing insert task statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Name,
		task.Params,
		task.Duration,
		task.LastRun,
		task.NextRun,
		task.IsRecurring,
		task.Hash,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting task: %s", err)
	}

	return nil
}
