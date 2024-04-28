package filer

import (
	"context"
	"time"
)

type Filer interface {
	Upsert(ctx context.Context, files []Files, userID string) ([]FilesResponse, error)
	Read(ctx context.Context, filename string, userID string) (string, error)
}

type Files struct {
	Filename string
	Message  string
}

type FilesResponse struct {
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
