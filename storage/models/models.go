package models

import "time"

type File struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Size        int       `json:"size"`
	Data        []byte    `json:"data"`
	ContentType string    `json:"content_type"`
	AddedAt     time.Time `json:"added_at"`
	UserID      int       `json:"user_id"`
	Hash        string    `json:"hash"`
}

type SimpleFileView struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Size        int       `json:"size"`
	ContentType string    `json:"content_type"`
	AddedAt     time.Time `json:"added_at"`
	UserID      int       `json:"user_id"`
}
