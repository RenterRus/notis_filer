package server

import (
	"context"
	"fmt"
	"notis_filer/internal/filer"
	"notis_filer/pkg/api"

	"google.golang.org/grpc/metadata"
)

type noteServer struct {
	MDsher filer.Filer
	api.UnimplementedNotesServer
}

func NewService(mdsher filer.Filer) *noteServer {
	return &noteServer{
		MDsher: mdsher,
	}
}

func (n *noteServer) Upsert(ctx context.Context, in *api.UpsertRequest) (*api.UpsertResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("X-User")) == 0 {
		return nil, fmt.Errorf("Upsert: %w", ErrUserNotFound)
	}

	resp, err := n.MDsher.Upsert(ctx, FilesToWorker(in.Files), md.Get("X-User")[0])
	return &api.UpsertResponse{
		Files: WorkerResponseToFiles(resp),
	}, err
}

func (n *noteServer) Read(ctx context.Context, in *api.ReadRequest) (*api.ReadResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get("X-User")) == 0 {
		return nil, fmt.Errorf("Read: %w", ErrUserNotFound)
	}

	message, err := n.MDsher.Read(ctx, filer.ReadData{
		FileName: in.Filename,
		UserID:   md.Get("X-User")[0],
		Topic:    in.Topic,
	})

	return &api.ReadResponse{
		Message: message,
	}, err
}
