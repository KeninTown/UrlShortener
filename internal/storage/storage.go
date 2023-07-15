package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	op := "Storage.New()"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	//посмотреть запрос
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url(
		id  INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
	`)

	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	op := "Storage.SaveUrl()"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")

	if err != nil {
		return -1, fmt.Errorf("%s : %w", op, err)
	}

	result, err := stmt.Exec(urlToSave, alias)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s : %w", op, ErrAliasExist)
		}

		return -1, fmt.Errorf("%s : %w", op, err)
	}

	urlId, err := result.LastInsertId()

	if err != nil {
		return -1, fmt.Errorf("%s: failed to get  last id: %w", op, err)
	}

	return urlId, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	op := "Storage.GetUrl()"

	stmt, err := s.db.Prepare("select url from url where alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var url string

	err = stmt.QueryRow(alias).Scan(&url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return url, nil
}

//TODO: implement delete url
