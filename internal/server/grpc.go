package server

import (
	"context"
	"errors"
	"notis_filer/pkg/api"
)

type noteServer struct {
	api.UnimplementedNotesServer
}

func NewService() *noteServer {
	return &noteServer{}
}

func (n *noteServer) Upsert(ctx context.Context, in *api.UpsertRequest) (*api.UpsertResponse, error) {
	return nil, errors.New("unimplemented")
}

func (n *noteServer) Read(ctx context.Context, in *api.ReadRequest) (*api.ReadResponse, error) {
	return nil, errors.New("unimplemented")
}
