package database

import (
	"context"
	"database/sql"
	"strings"

	"github.com/zeebo/errs"

	"phatic_dialogue/types"
)

// Templates provides access to templates db.
//
// architecture: Database
type Templates struct {
	conn *sql.DB
}

// Create creates template in the Database.
func (collectionsDB *Templates) Create(ctx context.Context, template types.Template) error {
	template.Template = strings.ToLower(template.Template)
	query := `INSERT INTO templates(template, topic) VALUES ($1, $2)`

	_, err := collectionsDB.conn.ExecContext(ctx, query, template.Template, template.Topic)

	return Error.Wrap(err)
}

// List returns all templates from the Database.
func (collectionsDB *Templates) List(ctx context.Context) (_ []types.Template, err error) {
	var list []types.Template

	query := `SELECT template, topic
 	          FROM templates
 	          ORDER BY topic ASC`

	rows, err := collectionsDB.conn.QueryContext(ctx, query)
	if err != nil {
		return list, Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	for rows.Next() {
		var template types.Template
		err := rows.Scan(&template.Template, &template.Topic)
		if err != nil {
			return list, Error.Wrap(err)
		}

		list = append(list, template)
	}
	if err = rows.Err(); err != nil {
		return list, Error.Wrap(err)
	}

	return list, nil
}
