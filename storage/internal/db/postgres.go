package db

import (
	"cmd/main.go/configs"
	"cmd/main.go/models"
	"context"
	"github.com/jackc/pgx"
	"log"
	"time"
)

type Database interface {
	Migrate() error
	GetFile(name string, userID string) (models.File, error)
	UploadFile(file models.File) error
	DeleteFile(name string, userID string) error
	GetFilesByUserID(userID string) ([]models.SimpleFileView, error)
}
type database struct {
	conn   *pgx.Conn
	config *configs.Config
}

func NewDatabase(cfg *configs.Config) (Database, error) {
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Database,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
	})
	if err != nil {
		return nil, err
	}
	for {
		err := conn.Ping(context.Background())
		if err == nil {
			break
		}
		log.Println("Database is still starting up, retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	return &database{conn: conn, config: cfg}, nil
}

func (d database) Migrate() error {
	sql := `
    CREATE TABLE IF NOT EXISTS files (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        data bytea NOT NULL,
        file_size INTEGER NOT NULL,
        content_type VARCHAR(255) NOT NULL,
        added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        user_id INTEGER NOT NULL ,
        hash VARCHAR(255) UNIQUE NOT NULL
    );
    `
	_, err := d.conn.Exec(sql)
	return err
}

func (d database) GetFile(name string, userID string) (models.File, error) {
	sql := `SELECT id, name,file_size,content_type,added_at,user_id,hash WHERE name=$1 AND user_id=$2`
	row := d.conn.QueryRow(sql, name, userID)
	var file models.File
	row.Scan(file)
	return file, nil
}

func (d database) UploadFile(file models.File) error {
	sql := `INSERT INTO files VALUES($1,$2,$3,$4,$5,$6)`
	_, err := d.conn.Exec(sql, file.Name, file.Size, file.Data, file.ContentType, file.UserID, file.Hash)
	return err

}

func (d database) DeleteFile(name string, userID string) error {
	sql := `DELETE FROM files WHERE name=$1 AND user_id=$2`
	_, err := d.conn.Exec(sql, name, userID)
	return err
}

func (d database) GetFilesByUserID(userID string) ([]models.SimpleFileView, error) {
	sql := `SELECT id, name, file_size, created_at FROM files WHERE user_id = $1`
	var files []models.SimpleFileView
	rows, err := d.conn.Query(sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Важно закрывать ресурсы после использования

	for rows.Next() {
		var file models.SimpleFileView
		err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.AddedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	// Проверка на ошибки, возникшие во время итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
