package filer

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
)

func NewFiler() Filer {
	return &mdsher{}
}

type mdsher struct {
}

/*
TODO: Добавить префикс для расположений файлов (через конфиг)
TODO: Научить создавать директории при необходимости
TODO: Добавить в файлы уникальные идентификаторы (например, по номеру топика).
		Для пердотвращения коллизий в случае двух оджинаковых имен файлов в разных топиках
*/

func (m *mdsher) Upsert(ctx context.Context, files []Files, userID string) ([]FilesResponse, error) {
	if files != nil {
		resp := make([]FilesResponse, len(files))

		for i := range files {
			file, err := os.Create(fmt.Sprintf("./%s/notis/%s.md", userID, files[i]))
			if err != nil {
				log.Warn().Msg(fmt.Sprintf("file cannot create: %v", err))
			}
			_, err = file.WriteString(files[i].Message)
			if err != nil {
				log.Warn().Msg(fmt.Sprintf("file cannot write: %v", err))
				resp[i].Message = fmt.Sprintf("%s %s", ErrFileCannotUpdate.Error(), files[i].Filename)
			} else {
				info, _ := file.Stat()

				resp[i].Message = fmt.Sprintf("file %s writed", files[i].Filename)
				resp[i].UpdatedAt = info.ModTime()
			}

			file.Close()
		}

		return resp, nil
	}

	return nil, ErrFilesEmpty
}
func (m *mdsher) Read(ctx context.Context, filename string, userID string) (string, error) {
	file, err := os.Open(fmt.Sprintf("./%s/notis/%s.md", userID, filename))
	if err != nil {
		return "", ErrFilesCannotRead
	}
	defer file.Close()

	message := make([]byte, 64)

	for {
		_, err := file.Read(message)
		if err == io.EOF {
			break
		}
	}

	return string(message), nil
}
