package filer

import "context"

func NewFiler() Filer {
	return &mdsher{}
}

type mdsher struct {
}

func (m *mdsher) Upsert(ctx context.Context, files []Files, userID string) ([]FilesResponse, error) {
	return nil, nil
}
func (m *mdsher) Read(ctx context.Context, filename string, userID string) (string, error) {
	return "", nil
}
