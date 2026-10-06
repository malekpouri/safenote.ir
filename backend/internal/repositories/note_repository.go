package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"safenote/internal/models"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound     = errors.New("note not found")
	ErrDuplicateID  = errors.New("duplicate note id")
	ErrUnauthorized = errors.New("unauthorized")
)

type NoteRepository struct {
	DB *sql.DB
}

const noteColumns = `id, encrypted_data, password_hash, access_token_hash, kdf_salt, is_password_protected, views_remaining, expires_at, created_at`

func NewNoteRepository(dbPath string) (*NoteRepository, error) {
	params := url.Values{}
	params.Set("_journal_mode", "WAL")
	params.Set("_busy_timeout", "5000")
	params.Set("_txlock", "immediate")
	db, err := sql.Open("sqlite3", "file:"+dbPath+"?"+params.Encode())
	if err != nil {
		return nil, err
	}
	// A single connection serializes all access, which makes the
	// read-check-decrement sequence in Consume atomic and avoids
	// "database is locked" errors under concurrent writes.
	db.SetMaxOpenConns(1)

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return &NoteRepository{DB: db}, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			encrypted_data TEXT,
			password_hash TEXT,
			is_password_protected BOOLEAN,
			views_remaining INTEGER,
			expires_at DATETIME,
			created_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS statistics (
			key TEXT PRIMARY KEY,
			value INTEGER
		)`,
		`INSERT OR IGNORE INTO statistics (key, value) VALUES ('total_notes', 0)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}

	for _, col := range []string{"access_token_hash", "kdf_salt"} {
		exists, err := hasColumn(db, "notes", col)
		if err != nil {
			return err
		}
		if !exists {
			if _, err := db.Exec(`ALTER TABLE notes ADD COLUMN ` + col + ` TEXT NOT NULL DEFAULT ''`); err != nil {
				return fmt.Errorf("migrate: add %s: %w", col, err)
			}
		}
	}

	_, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_notes_expires_at ON notes (expires_at)`)
	return err
}

func hasColumn(db *sql.DB, table, column string) (bool, error) {
	rows, err := db.Query(`SELECT name FROM pragma_table_info(?)`, table)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return false, err
		}
		if strings.EqualFold(name, column) {
			return true, nil
		}
	}
	return false, rows.Err()
}

func (r *NoteRepository) Close() error {
	return r.DB.Close()
}

func (r *NoteRepository) Create(note *models.Note) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`INSERT INTO notes (`+noteColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		note.ID, note.EncryptedData, note.PasswordHash, note.AccessTokenHash, note.KdfSalt,
		note.IsPasswordProtected, note.ViewsRemaining, note.ExpiresAt.UTC(), note.CreatedAt.UTC(),
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			return ErrDuplicateID
		}
		return err
	}

	if _, err = tx.Exec(`UPDATE statistics SET value = value + 1 WHERE key = 'total_notes'`); err != nil {
		return err
	}
	return tx.Commit()
}

type Stats struct {
	ActiveNotes       int `json:"active_notes"`
	TotalNotesCreated int `json:"total_notes_created"`
}

func (r *NoteRepository) GetStats(now time.Time) (*Stats, error) {
	var stats Stats
	err := r.DB.QueryRow(
		`SELECT COUNT(*) FROM notes WHERE expires_at > ? AND views_remaining > 0`, now.UTC(),
	).Scan(&stats.ActiveNotes)
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(`SELECT value FROM statistics WHERE key = 'total_notes'`).Scan(&stats.TotalNotesCreated)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return &stats, nil
}

func scanNote(row interface{ Scan(...any) error }) (*models.Note, error) {
	var n models.Note
	var passwordHash sql.NullString
	err := row.Scan(&n.ID, &n.EncryptedData, &passwordHash, &n.AccessTokenHash, &n.KdfSalt,
		&n.IsPasswordProtected, &n.ViewsRemaining, &n.ExpiresAt, &n.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	n.PasswordHash = passwordHash.String
	return &n, nil
}

func isLive(n *models.Note, now time.Time) bool {
	return n.ViewsRemaining > 0 && now.Before(n.ExpiresAt)
}

// GetByID returns a live note without consuming a view. Dead notes found
// along the way are deleted and reported as ErrNotFound.
func (r *NoteRepository) GetByID(id string, now time.Time) (*models.Note, error) {
	n, err := scanNote(r.DB.QueryRow(`SELECT `+noteColumns+` FROM notes WHERE id = ?`, id))
	if err != nil {
		return nil, err
	}
	if !isLive(n, now) {
		_ = r.Delete(id)
		return nil, ErrNotFound
	}
	return n, nil
}

// Consume atomically authorizes, decrements and (on the last view) deletes a
// note, returning it as it was before the decrement with ViewsRemaining set to
// the new count. authorize may return ErrUnauthorized to abort without
// consuming a view.
func (r *NoteRepository) Consume(id string, now time.Time, authorize func(*models.Note) error) (*models.Note, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	n, err := scanNote(tx.QueryRow(`SELECT `+noteColumns+` FROM notes WHERE id = ?`, id))
	if err != nil {
		return nil, err
	}
	if !isLive(n, now) {
		if _, err := tx.Exec(`DELETE FROM notes WHERE id = ?`, id); err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return nil, ErrNotFound
	}
	if err := authorize(n); err != nil {
		return nil, err
	}

	n.ViewsRemaining--
	if n.ViewsRemaining <= 0 {
		_, err = tx.Exec(`DELETE FROM notes WHERE id = ?`, id)
	} else {
		_, err = tx.Exec(`UPDATE notes SET views_remaining = ? WHERE id = ?`, n.ViewsRemaining, id)
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return n, nil
}

// DeleteAuthorized deletes a note if authorize allows it.
func (r *NoteRepository) DeleteAuthorized(id string, now time.Time, authorize func(*models.Note) error) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	n, err := scanNote(tx.QueryRow(`SELECT `+noteColumns+` FROM notes WHERE id = ?`, id))
	if err != nil {
		return err
	}
	if isLive(n, now) {
		if err := authorize(n); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM notes WHERE id = ?`, id); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	if !isLive(n, now) {
		return ErrNotFound
	}
	return nil
}

func (r *NoteRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM notes WHERE id = ?`, id)
	return err
}

func (r *NoteRepository) DeleteExpired(now time.Time) (int64, error) {
	result, err := r.DB.Exec(`DELETE FROM notes WHERE expires_at <= ? OR views_remaining <= 0`, now.UTC())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
