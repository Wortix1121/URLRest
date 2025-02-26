package rest

import (
	"api/internal/config"
	"api/internal/storage"
	"api/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	db *pgx.Conn
}

func Connect() (*Storage, error) {
	const op = "storage.postgres.pgconn.Connect"
	var storagePathPg = config.MustLoad().StoragePathPg

	db, err := pgx.Connect(context.Background(), storagePathPg)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s Storage) SaveUrl(ctx context.Context, alias string, urlToSave string) (int64, error) {
	const op = "storage.postgres.rest.SaveUrl"

	query := "INSERT INTO url(url, alias) VALUES($1, $2)"

	_, err := s.db.Exec(ctx, query, urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return 1, err

}

func (s Storage) GetUrlall(ctx context.Context) ([]model.UrlModel, error) {
	const op = "storage.postgres.rest.GetUrlall"

	query := "SELECT id, alias, url FROM url LIMIT 10"

	stmt, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	// var urls = []model.UrlModel{}
	// for stmt.Next() {
	// 	u := model.UrlModel{}
	// 	err := stmt.Scan(&u.Alias, &u.Url)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%s: %w", op, err)
	// 	}
	// 	urls = append(urls, u)
	// }
	// return urls, nil

	return pgx.CollectRows(stmt, pgx.RowToStructByNameLax[model.UrlModel])
}

func (s Storage) GetUrl(alias string, ctx context.Context) (string, error) {
	const op = "storage.postgres.rest.GetUrl"

	query := "SELECT url FROM url WHERE alias=$1"

	var url string
	err := s.db.QueryRow(ctx, query, alias).Scan(&url)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}

func (s Storage) DeleteUrl(alias string, ctx context.Context) (int64, error) {
	const op = "storage.postgres.rest.DeleteUrl"

	query := "DELETE FROM url WHERE alias=$1"

	_, err := s.db.Exec(ctx, query, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return 1, err

}

func (s Storage) NewTableUrls(ctx context.Context) (*pgconn.CommandTag, error) {
	const op = "storage.postgres.pgconn.NewTableUrls"

	query := `CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`
	// IF NOT EXISTS
	new, err := s.db.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &new, nil

}
