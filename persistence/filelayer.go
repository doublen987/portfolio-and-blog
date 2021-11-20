package persistence

import (
	"errors"

	"github.com/doublen987/Projects/MySite/server/persistence/filehandler/filesystemfilehandler"
	"github.com/doublen987/Projects/MySite/server/persistence/filehandler/s3filehandler"
)

const (
	S3         string = "s3"
	FILESYSTEM        = "filesystem"
)

type FileHandler interface {
	AddFile(file []byte, filename string) (string, error)
	GetFile(fileName string) ([]byte, error)
}

var FileHandlerTypeNotSupported = errors.New("The Database type provided is not supported...")

func GetFileHandler(fileHandlerType string, connection string) (FileHandler, error) {
	switch fileHandlerType {
	case "s3":
		return s3filehandler.NewS3FileHandler()
	case "filesystem":
		return filesystemfilehandler.NewFileSystemFileHandler()
	}
	return nil, FileHandlerTypeNotSupported
}
