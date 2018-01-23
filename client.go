package storage

import (
	"io"
	"os"

	minio "github.com/minio/minio-go"
)

// Clienter client interface
type Clienter interface {
	PutFile(reader io.Reader, objectName string, objectSize int64, contentType string) (n int64, err error)
	DelFile(objectName string) error
	GetFileURL(fileName string) string
	SetupClient(c *minio.Client)
}

// Client storage client
type Client struct {
	c          *minio.Client
	bucketName string
	basePath   string
}

// PutFile put file object into storage
func (cl *Client) PutFile(reader io.Reader, objectName string, objectSize int64, contentType string) (n int64, err error) {
	opts := minio.PutObjectOptions{
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
		ContentType:  contentType,
	}
	return cl.c.PutObject(cl.bucketName, objectName, reader, objectSize, opts)
}

// DelFile delete file object from storage
func (cl *Client) DelFile(objectName string) error {
	return cl.c.RemoveObject(cl.bucketName, objectName)
}

// GetFileURL returns file URL
func (cl *Client) GetFileURL(fileName string) string {
	// cl.c.PresignedGetObject(cl.bucketName, cl.fileName, expires, reqParams)
	return "https://" + cl.bucketName + "." + cl.basePath + "/" + fileName
}

// SetupClient returns file URL
func (cl *Client) SetupClient(c *minio.Client) {
	cl.c = c
}

// GetClient get filestorage client
func GetClient() (cl Clienter, err error) {
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secKey := os.Getenv("STORAGE_SECRET_ACCESS_KEY")
	endPoint := os.Getenv("STORAGE_ENDPOINT")
	cl = &Client{
		bucketName: os.Getenv("STORAGE_BUCKET_NAME"),
		basePath:   os.Getenv("STORAGE_BASE_PATH"),
	}
	minioClient, err := minio.New(endPoint, accessKey, secKey, false)
	if err != nil {
		return
	}
	cl.SetupClient(minioClient)
	return
}
