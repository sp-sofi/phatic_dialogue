package database

import (
	"context"
	"database/sql"
	"strings"

	"github.com/zeebo/errs"

	"phatic_dialogue/types"
)

// GroupInserts provides access to group_inserts db.
//
// architecture: Database
type GroupInserts struct {
	conn *sql.DB
}

// Create creates groupInsert in the Database.
func (collectionsDB *GroupInserts) Create(ctx context.Context, groupInsert types.GroupInsert) error {
	groupInsert.Words = strings.ToLower(groupInsert.Words)
	query := `INSERT INTO group_inserts(words, topic) VALUES ($1, $2)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, groupInsert.Words, groupInsert.Topic)

	return Error.Wrap(err)
}

// List returns all group inserts or by topic from the Database.
func (collectionsDB *GroupInserts) List(ctx context.Context, topic types.Topic) (_ []types.GroupInsert, err error) {
	var list []types.GroupInsert
	var args = make([]any, 0, 1)

	query := `SELECT words, topic
 	          FROM group_inserts
 	          `

	if len(topic) != 0 {
		query += `WHERE topic = $1`
		args = append(args, topic)
	}

	rows, err := collectionsDB.conn.QueryContext(ctx, query, args...)
	if err != nil {
		return list, Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var groupInsert types.GroupInsert
		err := rows.Scan(&groupInsert.Words, &groupInsert.Topic)
		if err != nil {
			return list, Error.Wrap(err)
		}

		list = append(list, groupInsert)
	}
	if err = rows.Err(); err != nil {
		return list, Error.Wrap(err)
	}

	return list, nil
}
