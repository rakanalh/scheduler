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
        last_run text,
        next_run text,
        is_recurring integer
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
	var count int
	rows, err := sqlite.db.Query("SELECT count(*) FROM task_store WHERE name=?", task.Name)
	if err == nil {
		rows.Next()
		_ = rows.Scan(&count)
	}
	rows.Close()

	if count > 0 {
		return sqlite.update(task)
	}
	return sqlite.insert(task)
}

func (sqlite Sqlite3Storage) Fetch() ([]*task.ScheduledTask, error) {
	rows, err := sqlite.db.Query(`
        SELECT name, duration, last_run, next_run, is_recurring
        FROM task_store`)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tasks []*task.ScheduledTask

	for rows.Next() {
		var duration, isRecurring, status int
		var name, lastRunStr, nextRunStr string
		err = rows.Scan(&name, &duration, &lastRunStr, &nextRunStr, &isRecurring)
		if err != nil {
			return []*task.ScheduledTask{}, err
		}
		if status == StatusDone {
			continue
		}

		lastRun, err := time.Parse(time.RFC3339, lastRunStr)
		if err != nil {
			return []*task.ScheduledTask{}, err
		}

		nextRun, err := time.Parse(time.RFC3339, nextRunStr)
		if err != nil {
			return []*task.ScheduledTask{}, err
		}

		tasks = append(tasks, &task.ScheduledTask{
			Name: name,
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

func (sqlite *Sqlite3Storage) insert(task *task.ScheduledTask) error {
	stmt, err := sqlite.db.Prepare(`
        INSERT INTO task_store(name, duration, last_run, next_run, is_recurring)
        VALUES(?, ?, ?, ?, ?)`)

	if err != nil {
		return fmt.Errorf("Error while pareparing insert task statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Name,
		task.Duration,
		task.LastRun.Format(time.RFC3339),
		task.NextRun.Format(time.RFC3339),
		task.IsRecurring,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting task: %s", err)
	}

	return nil
}

func (sqlite *Sqlite3Storage) update(task *task.ScheduledTask) error {
	stmt, err := sqlite.db.Prepare(`
        UPDATE task_store SET duration=?, last_run=?, next_run=?, is_recurring=?
        WHERE name=?`)

	if err != nil {
		return fmt.Errorf("Error while preparing update task statement: %s", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		task.Duration,
		task.LastRun.Format(time.RFC3339),
		task.NextRun.Format(time.RFC3339),
		task.IsRecurring,
		task.Name,
	)
	if err != nil {
		return fmt.Errorf("Error while updating task: %s", err)
	}

	return nil
}
