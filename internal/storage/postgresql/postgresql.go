package postgresql

import (
	"Url-shortner/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

// создаёт соединение с базой данных
func New() (*Storage, error) {
	operation := "storage.postgresql.New" // передать в ошибку(чтобы удобно было находить)

	connStr := "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable"
	fmt.Printf("connection: %s, ", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{
		db: db,
	}, nil
}

// сохранение урла в базу данных
func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	op := "storage.postgresql.SaveURL"
	var id int64

	err := s.db.QueryRow("INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id", urlToSave, alias).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrUrlExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// получение урла из базы данных по id
func (s *Storage) GetURLById(id int64) (string, error) {
	op := "storage.postgresql.GetURL"
	var url string

	err := s.db.QueryRow("SELECT url FROM url WHERE id = $1", id).Scan(&url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

// получение урла из базы данных по алиасу
func (s *Storage) GetURLByAlias(alias string) (string, error) {
	op := "storage.postgresql.GetURLByAlias"
	var url string

	err := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias).Scan(&url)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

// получение всех урлов из базы данных
func (s *Storage) GetAllURLs() ([]string, error) {
	op := "storage.postgresql.GetAllURLs"
	var urls []string

	rows, err := s.db.Query("SELECT url FROM url")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var url string
		err := rows.Scan(&url)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		urls = append(append(urls, "\n"), url)
	}

	return urls, nil
}

// удаление урла из базы данных по алиасу
func (s *Storage) DeleteURLByAlias(alias string) (string, error) {
	op := "storage.postgresql.DeleteURLByAlias"

	_, err := s.db.Exec("DELETE FROM url WHERE alias = $1", alias)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "Запись успешно удалена", nil
}
