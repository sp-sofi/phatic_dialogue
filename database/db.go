package database

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq" // using postgres driver.
	"github.com/zeebo/errs"
)

var (
	// Error is the default db error class.
	Error      = errs.Class("db error")
	ErrNoTopic = errors.New("topic does not exist")
)

// Database provides access to Database tables.
//
// architecture: Master Database
type Database struct {
	conn *sql.DB

	templates     *Templates
	answers       *Answers
	singleInserts *SingleInserts
	groupInserts  *GroupInserts
	topics        *Topics
}

// New is a constructor for Database.
func New(databaseURL string) (*Database, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	db := Database{conn: conn}
	return &db, nil
}

// CreateSchema creates schema for all tables and databases.
func (db *Database) CreateSchema(ctx context.Context) (err error) {
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS topics (
            topic   VARCHAR   PRIMARY KEY   NOT NULL
        );
        CREATE TABLE IF NOT EXISTS single_inserts (
            id      SERIAL    PRIMARY KEY                NOT NULL,
            word    VARCHAR                              NOT NULL,
            topic   VARCHAR   REFERENCES topics(topic)   NOT NULL
        );
        CREATE TABLE IF NOT EXISTS group_inserts (
            id      SERIAL    PRIMARY KEY                NOT NULL,
            words   VARCHAR                              NOT NULL,
            topic   VARCHAR   REFERENCES topics(topic)   NOT NULL
        );
 		CREATE TABLE IF NOT EXISTS templates (
 		    id          SERIAL    PRIMARY KEY                NOT NULL,
            template    VARCHAR                              NOT NULL,
            topic       VARCHAR   REFERENCES topics(topic)   NOT NULL
        );
		CREATE TABLE IF NOT EXISTS answers (
		    id         SERIAL    PRIMARY KEY                NOT NULL,
            answer     VARCHAR                              NOT NULL,
		    topic      VARCHAR   REFERENCES topics(topic)   NOT NULL
        );
       `

	_, err = db.conn.ExecContext(ctx, createTableQuery)
	if err != nil {
		return Error.Wrap(err)
	}

	return nil
}

// Answers returns connection to answers db.
func (db *Database) Answers() *Answers {
	if db.answers == nil {
		db.answers = &Answers{conn: db.conn}
	}

	return db.answers
}

// Templates returns connection to templates db.
func (db *Database) Templates() *Templates {
	if db.templates == nil {
		db.templates = &Templates{conn: db.conn}
	}

	return db.templates
}

// SingleInserts returns connection to singleInserts db.
func (db *Database) SingleInserts() *SingleInserts {
	if db.singleInserts == nil {
		db.singleInserts = &SingleInserts{conn: db.conn}
	}

	return db.singleInserts
}

// GroupInserts returns connection to groupInserts db.
func (db *Database) GroupInserts() *GroupInserts {
	if db.groupInserts == nil {
		db.groupInserts = &GroupInserts{conn: db.conn}
	}

	return db.groupInserts
}

// Topics returns connection to topics db.
func (db *Database) Topics() *Topics {
	if db.topics == nil {
		db.topics = &Topics{conn: db.conn}
	}

	return db.topics
}

// Close closes underlying db connection.
func (db *Database) Close() error {
	return Error.Wrap(db.conn.Close())
}
