package functionality

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func SaveFile(fileBytes []byte) (string, error) {
	fmt.Println("Uploading File\n")

	//3. write temporary file on our server
	tempFile, err := ioutil.TempFile("webportal/content/images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer tempFile.Close()

	tempFile.Write(fileBytes)

	fmt.Printf("Sucessfully uploaded file %s\n", tempFile.Name())
	fmt.Println(strings.ReplaceAll(tempFile.Name(), "webportal/content/images/", ""))
	return strings.ReplaceAll(tempFile.Name(), "webportal/content/images/", ""), nil
}
