package server

import (
	"notis_filer/internal/filer"
	"notis_filer/pkg/api"
)

func FilesToWorker(files []*api.UpsertRequest_Files) []filer.Files {
	res := make([]filer.Files, len(files))
	for i := range files {
		res[i] = filer.Files{
			Filename: files[i].Filename,
			Message:  files[i].Message,
			Topic:    files[i].Topic,
		}
	}
	return res
}

func WorkerResponseToFiles(files []filer.FilesResponse) []*api.UpsertResponse_Files {
	res := make([]*api.UpsertResponse_Files, len(files))
	for i := range files {
		res[i] = &api.UpsertResponse_Files{
			Message: files[i].Message,
		}
	}
	return res
}
