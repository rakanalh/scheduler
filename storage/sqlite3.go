package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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
	config    Sqlite3Config
	db        *sql.DB
	marshaler ParamMarshaler
}

func NewSqlite3Storage(config Sqlite3Config, marshaler ParamMarshaler) Sqlite3Storage {
	if marshaler == nil {
		marshaler = NewMarshaler()
	}
	return Sqlite3Storage{
		config:    config,
		marshaler: marshaler,
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
        params text,
        duration integer,
        last_run text,
        next_run text,
        is_recurring integer,
        hash text
    );
	`
	_, err := sqlite.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return nil
}

func (sqlite Sqlite3Storage) Add(task *task.Task) error {
	var count int
	rows, err := sqlite.db.Query("SELECT count(*) FROM task_store WHERE hash=?", task.Hash())
	if err == nil {
		rows.Next()
		_ = rows.Scan(&count)
	}
	_ = rows.Close()

	if count == 0 {
		return sqlite.insert(task)
	}
	return nil
}

func (sqlite Sqlite3Storage) Remove(task *task.Task) error {
	stmt, err := sqlite.db.Prepare(`DELETE FROM task_store WHERE hash=?`)

	if err != nil {
		return fmt.Errorf("Error while pareparing delete task statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Hash(),
	)
	if err != nil {
		return fmt.Errorf("Error while deleting task: %s", err)
	}

	return nil
}

func (sqlite Sqlite3Storage) Fetch() ([]*task.Task, error) {
	rows, err := sqlite.db.Query(`
        SELECT name, params, duration, last_run, next_run, is_recurring
        FROM task_store`)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var tasks []*task.Task

	for rows.Next() {
		var duration, isRecurring, status int
		var name, paramsStr, lastRunStr, nextRunStr string
		err = rows.Scan(&name, &paramsStr, &duration, &lastRunStr, &nextRunStr, &isRecurring)
		if err != nil {
			return []*task.Task{}, err
		}
		if status == StatusDone {
			continue
		}

		lastRun, err := time.Parse(time.RFC3339, lastRunStr)
		if err != nil {
			return []*task.Task{}, err
		}

		nextRun, err := time.Parse(time.RFC3339, nextRunStr)
		if err != nil {
			return []*task.Task{}, err
		}

		params, err := sqlite.marshaler.Unmarshal(paramsStr)
		if err != nil {
			return []*task.Task{}, err
		}

		tasks = append(tasks, &task.Task{
			FuncName: name,
			Params:   params,
			Schedule: task.Schedule{
				Duration:    time.Duration(duration),
				IsRecurring: isRecurring == 1,
				LastRun:     lastRun,
				NextRun:     nextRun,
			},
		})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tasks, nil
}

func (sqlite *Sqlite3Storage) insert(task *task.Task) error {
	stmt, err := sqlite.db.Prepare(`
        INSERT INTO task_store(name, params, duration, last_run, next_run, is_recurring, hash)
        VALUES(?, ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return fmt.Errorf("Error while pareparing insert task statement: %s", err)
	}

	defer stmt.Close()

	params, err := sqlite.marshaler.Marshal(task.Params)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		task.FuncName,
		params,
		task.Duration,
		task.LastRun.Format(time.RFC3339),
		task.NextRun.Format(time.RFC3339),
		task.IsRecurring,
		task.Hash(),
	)
	if err != nil {
		return fmt.Errorf("Error while inserting task: %s", err)
	}

	return nil
}
