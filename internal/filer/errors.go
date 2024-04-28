package filer

import "errors"

var (
	ErrFileCannotUpdate = errors.New("couldn't update the file")
	ErrFilesEmpty       = errors.New("there is no information about the files")
	ErrFilesCannotRead  = errors.New("couldn't read the file")
)
