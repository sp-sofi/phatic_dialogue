package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/zeebo/errs"

	"phatic_dialogue/types"
)

// Topics provides access to topics db.
//
// architecture: Database
type Topics struct {
	conn *sql.DB
}

// Create creates topic in the Database.
func (collectionsDB *Topics) Create(ctx context.Context, topic types.Topic) error {
	topic = types.Topic(strings.ToLower(string(topic)))
	query := `INSERT INTO topics(topic) VALUES ($1)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, topic)

	return Error.Wrap(err)
}

// Get returns topic by id from the Database.
func (collectionsDB *Topics) Get(ctx context.Context, topic types.Topic) (types.Topic, error) {
	query := `SELECT topic
 	          FROM topics
 	          WHERE topic = $1`

	err := collectionsDB.conn.QueryRowContext(ctx, query, topic).Scan(&topic)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoTopic
	}

	return topic, Error.Wrap(err)
}

// List returns all topics from the Database.
func (collectionsDB *Topics) List(ctx context.Context) (_ []types.Topic, err error) {
	var list []types.Topic

	query := `SELECT topic
 	          FROM topics`

	rows, err := collectionsDB.conn.QueryContext(ctx, query)
	if err != nil {
		return list, Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var topic types.Topic
		err := rows.Scan(&topic)
		if err != nil {
			return list, Error.Wrap(err)
		}

		list = append(list, topic)
	}
	if err = rows.Err(); err != nil {
		return list, Error.Wrap(err)
	}

	return list, nil
}
