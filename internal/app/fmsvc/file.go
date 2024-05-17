package fmsvc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type File struct {
	Id         int64     `json:"id"`
	FileName   string    `json:"filename"`
	SecondName string    `json:"-"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	UploadedAt time.Time `json:"uploaded"`
}

type FileStore interface {
	Save(ctx context.Context, file File) error
	GetFile(ctx context.Context, filename string) (File, error)
	GetListBy(ctx context.Context, column, orderby string) ([]File, error)
	SerachByName(ctx context.Context, filename string) ([]File, error)
}

type FileStorage struct {
	db *sql.DB
}

func NewFileStorage(db *sql.DB) FileStore {
	return &FileStorage{
		db: db,
	}
}

// Save implements FileService.
func (fs *FileStorage) Save(ctx context.Context, file File) error {

	_, err := fs.db.ExecContext(ctx, "INSERT INTO files ( filename, secondname, path, size, uploaded_at) VALUES ($1, $2, $3, $4, $5)",
		file.FileName,
		file.SecondName,
		file.Path,
		file.Size,
		file.UploadedAt)
	if err != nil {
		return err
	}

	return nil
}

// get file
func (fs *FileStorage) GetFile(ctx context.Context, filename string) (File, error) {

	row := fs.db.QueryRowContext(ctx, "SELECT * FROM files WHERE filename = $1", filename)
	var file File
	err := row.Scan(
		&file.Id,
		&file.FileName,
		&file.SecondName,
		&file.Path,
		&file.Size,
		&file.UploadedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return File{}, errors.New("file not found")
		}
		return File{}, err
	}

	return file, nil
}

// GetAll implements FileService.
func (fs *FileStorage) GetListBy(ctx context.Context, column, orderby string) ([]File, error) {

	if !isValidColumnName(column) {

		return nil, errors.New("invalid column name")
	}

	if !isValidOrder(orderby) {

		return nil, errors.New("invalid order name")
	}

	query := fmt.Sprintf("SELECT * FROM files ORDER BY %s %s", column, orderby)

	rows, err := fs.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File
		if err := rows.Scan(
			&file.Id,
			&file.FileName,
			&file.SecondName,
			&file.Path,
			&file.Size,
			&file.UploadedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func isValidOrder(columnName string) bool {
	// Add your validation logic here, e.g., checking against a predefined list of allowed column names
	allowedColumnNames := map[string]bool{
		"asc":  true,
		"ASC":  true,
		"desc": true,
		"DESC": true,
	}
	return allowedColumnNames[columnName]
}

func isValidColumnName(orderName string) bool {
	// Add your validation logic here, e.g., checking against a predefined list of allowed column names
	allowedOrderNames := map[string]bool{
		"filename":    true,
		"FILENAME":    true,
		"size":        true,
		"SIZE":        true,
		"uploaded_at": true,
		"UPLOADED_AT": true,
	}
	return allowedOrderNames[orderName]
}

func (fs *FileStorage) SerachByName(ctx context.Context, filename string) ([]File, error) {

	rows, err := fs.db.QueryContext(ctx, "SELECT * FROM files WHERE filename ILIKE '%' || $1 || '%' ORDER BY filename ASC", filename)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err == sql.ErrNoRows {
		return nil, err
	}

	var files []File

	for rows.Next() {
		var file File
		if err := rows.Scan(
			&file.Id,
			&file.FileName,
			&file.SecondName,
			&file.Path,
			&file.Size,
			&file.UploadedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
