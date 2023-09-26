package database

import (
	"context"
	"database/sql"
	"strings"

	"github.com/zeebo/errs"

	"phatic_dialogue/types"
)

// Answers provides access to answers db.
//
// architecture: Database
type Answers struct {
	conn *sql.DB
}

// Create creates general answers in the Database.
func (collectionsDB *Answers) Create(ctx context.Context, answer types.Answer) error {
	answer.Answer = strings.ToLower(answer.Answer)
	query := `INSERT INTO answers(answer, topic) VALUES ($1, $2)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, answer.Answer, answer.Topic)

	return Error.Wrap(err)
}

// List returns all general answers or by topic from the Database.
func (collectionsDB *Answers) List(ctx context.Context, topic types.Topic) (_ []types.Answer, err error) {
	var list []types.Answer
	var args = make([]any, 0, 1)

	query := `SELECT answer, topic
 	          FROM answers
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
		var answer types.Answer
		err := rows.Scan(&answer.Answer, &answer.Topic)
		if err != nil {
			return list, Error.Wrap(err)
		}

		list = append(list, answer)
	}
	if err = rows.Err(); err != nil {
		return list, Error.Wrap(err)
	}

	return list, nil
}
