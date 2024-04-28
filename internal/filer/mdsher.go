package filer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	fileMask = "../%s/%s/notis/%d_%s.md"
	pathMask = "../%s/%s/notis"
)

func NewFiler(baseDir string) Filer {
	return &mdsher{
		baseDir: baseDir,
	}
}

type mdsher struct {
	baseDir string
}

func (m *mdsher) Upsert(ctx context.Context, files []Files, userID string) ([]FilesResponse, error) {
	if files != nil {
		resp := make([]FilesResponse, len(files))
		os.MkdirAll(fmt.Sprintf(pathMask, m.baseDir, userID), os.ModePerm)

		for i := range files {
			file, err := os.Create(fmt.Sprintf(fileMask, m.baseDir, userID, files[i].Topic, files[i].Filename))
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
func (m *mdsher) Read(ctx context.Context, rd ReadData) (string, error) {
	/*mess, err := os.ReadFile(fmt.Sprintf(fileMask, m.baseDir, rd.UserID, rd.Topic, rd.FileName))
	return string(mess), err*/

	file, err := os.Open(fmt.Sprintf(fileMask, m.baseDir, rd.UserID, rd.Topic, rd.FileName))
	if err != nil {
		return "", ErrFilesCannotRead
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var sb strings.Builder

	for scanner.Scan() {
		sb.WriteString(scanner.Text())
	}

	return sb.String(), scanner.Err()
}
