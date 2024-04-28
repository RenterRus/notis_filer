package filer

import (
	"context"
	"time"
)

type ReadData struct {
	FileName string
	UserID   string
	Topic    int64
}

type Filer interface {
	Upsert(ctx context.Context, files []Files, userID string) ([]FilesResponse, error)
	Read(ctx context.Context, rd ReadData) (string, error)
}

type Files struct {
	Filename string
	Message  string
	Topic    int64
}

type FilesResponse struct {
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
