package mysql

import (
	"database/sql"
	"errors"

	"montgomery.wg/http-server/pkg/models"
)

// Define a TemplateModel type which wraps a sql.DB connection pool.
type TemplateModel struct {
    DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *TemplateModel) Insert(repo, template, version, expires string) (int, error) {
    stmt := `INSERT INTO templates (repo, template, version, created, expires)
    VALUES(?, ?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

    result, err := m.DB.Exec(stmt, repo, template, version)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *TemplateModel) Get(id int) (*models.Template, error) {
    stmt := `SELECT id, repo, template, version, created, expires FROM templates
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	
	row := m.DB.QueryRow(stmt, id)

	s := &models.Template{}

	err := row.Scan(&s.ID, &s.Repo, &s.Template, &s.Version, &s.Created, &s.Expires)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, models.ErrNoRecord
        } else {
             return nil, err
        }
	}
	
	return s, nil
}

// This will return the 10 most recently created templates.
func (m *TemplateModel) Latest() ([]*models.Template, error) {
    return nil, nil
}