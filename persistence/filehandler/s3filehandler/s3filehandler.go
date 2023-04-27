package s3filehandler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	uuid "github.com/nu7hatch/gouuid"
)

type S3FileHandler struct {
	client       *s3.Client
	imagesBucket string
}

var InvalidS3ImageBucketName error = errors.New("invalid S3 image bucket name")
var InvalidS3Region error = errors.New("invalid S3 region")
var InvalidS3AccessKey error = errors.New("invalid S3 access key")

func NewS3FileHandler() (*S3FileHandler, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
	}

	imagesBucketName := os.Getenv("AWS_S3_IMAGE_BUCKET_NAME")
	if imagesBucketName == "" {
		return &S3FileHandler{}, InvalidS3ImageBucketName
	}

	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return &S3FileHandler{}, InvalidS3Region
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return &S3FileHandler{}, InvalidS3AccessKey
	}

	client := s3.NewFromConfig(cfg)

	return &S3FileHandler{
		client:       client,
		imagesBucket: imagesBucketName,
	}, nil
}

var InvalidFileName = errors.New("invalid file name")

func (fh S3FileHandler) AddFile(file []byte, filename string) (string, error) {

	uuid, err := uuid.NewV4()
	if err != nil {
		//fmt.Println(err.Error())
		return "", err
	}

	s := strings.Split(filename, ".")
	if len(s) < 1 {
		return "", InvalidFileName
	}

	fileextension := s[len(s)-1]
	newfilename := uuid.String() + "." + fileextension

	input := s3.PutObjectInput{
		Bucket: aws.String(fh.imagesBucket),
		Key:    aws.String(newfilename),
		//Body:          bytes.NewReader([]byte("PAYLOAD")),
		Body:          bytes.NewReader(file),
		ContentLength: int64(len(file)),
	}

	_, err = fh.client.PutObject(context.TODO(), &input)
	if err != nil {
		//fmt.Println(err.Error())
		return "", err
	}

	return newfilename, nil
}

func (fh S3FileHandler) GetFile(filename string) ([]byte, error) {

	output, err := fh.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(fh.imagesBucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		//fmt.Println(err.Error())
		return []byte{}, err
	}

	imageBytes, err := ioutil.ReadAll(output.Body)
	if err != nil {
		//fmt.Println(err.Error())
		return []byte{}, err
	}

	return imageBytes, nil
}

func (fh S3FileHandler) RemoveFile(filename string) error {
	input := s3.DeleteObjectInput{
		Bucket: aws.String(fh.imagesBucket),
		Key:    aws.String(filename),
	}

	_, err := fh.client.DeleteObject(context.TODO(), &input)
	return err
}
