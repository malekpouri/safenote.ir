package repositories

import (
	"database/sql"
	"safenote/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type NoteRepository struct {
	DB *sql.DB
}

func NewNoteRepository(dbPath string) (*NoteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// For development simplicity, drop table to ensure schema update
	// In production, use migrations!
	// _, _ = db.Exec("DROP TABLE IF EXISTS notes")

	query := `
	CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		encrypted_data TEXT,
		is_password_protected BOOLEAN,
		views_remaining INTEGER,
		expires_at DATETIME,
		created_at DATETIME
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &NoteRepository{DB: db}, nil
}

func (r *NoteRepository) Create(note *models.Note) error {
	query := `INSERT INTO notes (id, encrypted_data, is_password_protected, views_remaining, expires_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, note.ID, note.EncryptedData, note.IsPasswordProtected, note.ViewsRemaining, note.ExpiresAt, note.CreatedAt)
	return err
}

func (r *NoteRepository) GetByID(id string) (*models.Note, error) {
	query := `SELECT id, encrypted_data, is_password_protected, views_remaining, expires_at, created_at FROM notes WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var note models.Note
	err := row.Scan(&note.ID, &note.EncryptedData, &note.IsPasswordProtected, &note.ViewsRemaining, &note.ExpiresAt, &note.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepository) DecrementViews(id string) error {
	query := `UPDATE notes SET views_remaining = views_remaining - 1 WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *NoteRepository) Delete(id string) error {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
