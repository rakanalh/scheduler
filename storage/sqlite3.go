package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rakanalh/scheduler/task"
)

const (
	StatusPending = iota
	StatusDone
)

type Sqlite3Config struct {
	DbName string
}

type Sqlite3Storage struct {
	config Sqlite3Config
	db     *sql.DB
}

func NewSqlite3Storage(config Sqlite3Config) Sqlite3Storage {
	return Sqlite3Storage{
		config: config,
	}
}

func (sqlite *Sqlite3Storage) Connect() error {
	db, err := sql.Open("sqlite3", sqlite.config.DbName)
	if err != nil {
		return err
	}
	sqlite.db = db
	return nil
}

func (sqlite *Sqlite3Storage) Close() error {
	return sqlite.db.Close()
}

func (sqlite *Sqlite3Storage) Initialize() error {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS task_store (
        id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
        name text,
        duration integer,
        next_run text,
        is_recurring integer,
        status integer
    );
	`
	_, err := sqlite.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

func (sqlite Sqlite3Storage) Store(task *task.ScheduledTask) error {
	tx, err := sqlite.db.Begin()
	if err != nil {
		return err
	}
	// TODO: Check if task already exists and update it if it does.
	stmt, err := tx.Prepare(`
        INSERT INTO task_store(name, duration, next_run, is_recurring, status)
        VALUES(?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Name,
		task.Duration,
		task.NextRun,
		task.IsRecurring,
		StatusPending,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (sqlite Sqlite3Storage) Fetch() ([]*task.ScheduledTask, error) {
	rows, err := sqlite.db.Query(`
        SELECT name, duration, next_run, is_recurring, status
        FROM task_store`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var duration, isRecurring, status int
		var name, nextRun string
		err = rows.Scan(&name, &duration, nextRun, isRecurring)
		if err != nil {
			log.Fatal(err)
		}
		if status == StatusDone {
			continue
		}
		fmt.Println(name, duration, nextRun, isRecurring)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return nil, nil
}
