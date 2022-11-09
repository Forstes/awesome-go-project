package storage

import (
	"io"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/minio/minio-go"
)

type ObjectStorage interface {
	CreateBucket(bucket string, region string) error
	UploadObject(bucket string, objectName string, reader io.Reader, objSize int64, contentType string) error
	GetObjectPresigned(bucket string, path string) (string, error)
	GetObject(bucket string, path string) ([]byte, error)
}

type MinioStore struct {
	Client *minio.Client
}

func (m *MinioStore) GetObjectPresigned(bucket string, path string) (string, error) {
	urlPresigned, err := m.Client.PresignedGetObject(bucket, path, time.Second*60*60, make(url.Values))
	if err != nil {
		return "", err
	}
	return urlPresigned.Redacted(), nil
}

func (m *MinioStore) GetObject(bucket string, path string) ([]byte, error) {

	reader, err := m.Client.GetObject(bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	/*
		b := make([]byte, 4)
		for {
			_, err := reader.Read(b)
			if err == io.EOF {
				break
			}
		}*/

	println(b)
	return b, nil
}

func (m *MinioStore) UploadObject(bucket string, objectName string, reader io.Reader, objSize int64, contentType string) error {
	_, err := m.Client.PutObject(bucket, objectName, reader, objSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}

func (m *MinioStore) CreateBucket(bucket string, region string) error {
	err := m.Client.MakeBucket(bucket, region)
	if err != nil {
		exists, err := m.Client.BucketExists(bucket)
		if err == nil && exists {
			return ErrBucketExists
		} else {
			return err
		}
	}
	return nil
}
