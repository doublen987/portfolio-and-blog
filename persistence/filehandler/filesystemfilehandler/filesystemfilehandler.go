package filesystemfilehandler

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type FileSystemFileHandler struct {
}

func NewFileSystemFileHandler() (*FileSystemFileHandler, error) {
	return &FileSystemFileHandler{}, nil
}

func (fh FileSystemFileHandler) AddFile(file []byte, fileName string) (string, error) {
	fmt.Println("Uploading File")

	//3. write temporary file on our server
	tempFile, err := ioutil.TempFile("webportal/content/images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer tempFile.Close()

	byteslen, err := tempFile.Write(file)
	if err != nil {
		return "", err
	}

	fmt.Printf("Sucessfully uploaded file %s, %d\n", tempFile.Name(), byteslen)
	fmt.Println(strings.ReplaceAll(tempFile.Name(), "webportal/content/images/", ""))
	return strings.ReplaceAll(tempFile.Name(), "webportal/content/images/", ""), nil
}

func (fh FileSystemFileHandler) GetFile(fileName string) ([]byte, error) {
	//fmt.Println("Retrieving File", fileName)
	file, err := os.Open("webportal/content/images/" + fileName)
	if err != nil {
		return []byte{}, err
	}
	var retrievedBytes []byte
	buf := make([]byte, 1)
	for {
		_, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				return []byte{}, nil
			}
			break
		}
		retrievedBytes = append(retrievedBytes, buf...)
	}

	if err != nil {
		return []byte{}, err
	}
	//fmt.Printf("Sucessfully retrieving file %s %d\n", fileName, len(retrievedBytes))
	return retrievedBytes, nil
}

func (fh FileSystemFileHandler) RemoveFile(fileName string) error {
	fmt.Println("Uploading File")

	err := os.Remove("webportal/content/images/" + fileName)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Sucessfully deleted file %s\n", fileName)
	return nil
}
