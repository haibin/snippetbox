package mysql

import (
	"database/sql"
	"errors"

	"github.com/haibin/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
    DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `
		INSERT INTO snippets (
			title, 
			content, 
			created, 
			expires
		)
		VALUES(
			?, 
			?, 
			UTC_TIMESTAMP(), 
			DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY)
		)
	`
	result, err := m.DB.Exec(stmt, title, content, expires)
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
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires 
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() 
		AND id = ?
	`
	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            return nil, models.ErrNoRecord
        } else {
             return nil, err
        }
	}

    return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `
		SELECT id, title, content, created, expires 
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() 
		ORDER BY created DESC 
		LIMIT 10
	`

	rows, err := m.DB.Query(stmt)
    if err != nil {
        return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
        // Create a pointer to a new zeroed Snippet struct.
        s := &models.Snippet{}
        // Use rows.Scan() to copy the values from each field in the row to the
        // new Snippet object that we created. Again, the arguments to row.Scan()
        // must be pointers to the place you want to copy the data into, and the
        // number of arguments must be exactly the same as the number of
        // columns returned by your statement.
        err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
        if err != nil {
            return nil, err
        }
        // Append it to the slice of snippets.
        snippets = append(snippets, s)
	}
	
	if err = rows.Err(); err != nil {
        return nil, err
    }
	
    return snippets, nil
}
