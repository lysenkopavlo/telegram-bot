package sqlite

import (
	"context"
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"

	"github.com/lysenkopavlo/telegram-bot/internal/helpers/e"
	"github.com/lysenkopavlo/telegram-bot/internal/storage"
)

const msgDbError = "database connection error"

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		slog.Error(msgDbError, "error", err)
		return nil, e.WrapError(msgDbError, err)
	}
	if err := db.Ping(); err != nil {
		slog.Error(msgDbError, "error", err)
		return nil, e.WrapError(msgDbError, err)
	}

	slog.Info("Successfully connected to db")
	return &Storage{
		db: db,
	}, nil
}

// Init creates database.
func (s *Storage) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		slog.Error("database creation", "error", err)
		return err
	}

	slog.Info("Successfully created database")
	return nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	query := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, query, p.URL, p.UserName)
	if err != nil {
		slog.Error("database insertion", "error", err)
		return err
	}

	slog.Info("Successfully inserted page to db")
	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	query := `
	SELECT url 
	FROM pages 
	WHERE user_name = ?
	ORDERED BY RANDOM() 
	LIMIT 1
	`
	var resURL string
	err := s.db.QueryRowContext(ctx, query, userName).Scan(&resURL)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		slog.Error("database pickRandom", "error", err)
		return nil, err
	}

	slog.Info("Successfully picked random page")
	return &storage.Page{
		URL:      resURL,
		UserName: userName,
	}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	query := "DELETE FROM pages WHERE url = ? AND user_name = ?"
	_, err := s.db.ExecContext(ctx, query, p.URL, p.UserName)
	if err != nil {
		slog.Error("database insertion", "error", err)
		return err
	}

	slog.Info("Successfully removed page")
	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	query := "SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?"
	var count int

	err := s.db.QueryRowContext(ctx, query, p.URL, p.UserName).Scan(&count)
	if err != nil {
		slog.Error("database insertion", "error", err)
		return false, err
	}

	slog.Info("Successfully removed page")

	return count > 0, nil
}
