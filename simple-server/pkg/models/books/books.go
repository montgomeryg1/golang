package books

import (
	"gorm.io/gorm"
	"montgomery.wg/simple-server/pkg/models"
)

// Define a BooksModel type which wraps a sql.DB connection pool.
type BooksModel struct {
    DB *gorm.DB
}

// func (m *BooksModel) Insert(title, content, expires string) (int, error) {
//     stmt := `INSERT INTO snippets (title, content, created, expires)
//     VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

//     result, err := m.DB.Exec(stmt, title, content, expires)
//     if err != nil {
//         return 0, err
//     }

//     id, err := result.LastInsertId()
//     if err != nil {
//         return 0, err
//     }

//     return int(id), nil
// }

// func (m *BooksModel) Get(id int) (*models.Book, error) {
//     stmt := `SELECT id, title, content, created, expires FROM snippets
// 	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	
// 	row := m.DB.QueryRow(stmt, id)

// 	s := &models.Book{}

// 	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
//     if err != nil {
//         if errors.Is(err, sql.ErrNoRows) {
//             return nil, models.ErrNoRecord
//         } else {
//              return nil, err
//         }
// 	}
	
// 	return s, nil
// }

func (m *BooksModel) Latest() ([]*models.Book, error) {
    return nil, nil
}