package database

import (
	"context"
	"database/sql"
	"strings"

	"github.com/zeebo/errs"

	"phatic_dialogue/types"
)

// SingleInserts provides access to single_inserts db.
//
// architecture: Database
type SingleInserts struct {
	conn *sql.DB
}

// Create creates singleInsert in the Database.
func (collectionsDB *SingleInserts) Create(ctx context.Context, singleInsert types.SingleInsert) error {
	singleInsert.Word = strings.ToLower(singleInsert.Word)
	query := `INSERT INTO single_inserts(word, topic) VALUES ($1, $2)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, singleInsert.Word, singleInsert.Topic)

	return Error.Wrap(err)
}

// List returns all single inserts or by topic from the Database.
func (collectionsDB *SingleInserts) List(ctx context.Context, topic types.Topic) (_ []types.SingleInsert, err error) {
	var list []types.SingleInsert
	var args = make([]any, 0, 1)

	query := `SELECT word, topic
 	          FROM single_inserts
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
		var singleInsert types.SingleInsert
		err := rows.Scan(&singleInsert.Word, &singleInsert.Topic)
		if err != nil {
			return list, Error.Wrap(err)
		}

		list = append(list, singleInsert)
	}
	if err = rows.Err(); err != nil {
		return list, Error.Wrap(err)
	}

	return list, nil
}
