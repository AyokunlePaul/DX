package file

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"time"
)

const (
	bucketUrl          = "https://storage.googleapis.com"
	bucketName         = "errand_bucket"
	categoryBucketName = "category_icons"
)

type repository struct {
	*storage.Client
}

func NewRepository(client *storage.Client) Repository {
	return &repository{
		Client: client,
	}
}

func (r *repository) Create(userId string, errandId string, file *File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	path := fmt.Sprintf("%s/%s/%s/%s", file.Folder, userId, errandId, file.Header.Filename)
	storageWriter := r.Bucket(bucketName).Object(path).NewWriter(ctx)

	headerFile, err := file.Header.Open()
	if err != nil {
		return err
	}

	if _, err = io.Copy(storageWriter, headerFile); err != nil {
		return err
	}

	if err = storageWriter.Close(); err != nil {
		return err
	}

	file.UploadedUrl = bucketUrl + "/" + bucketName + "/" + path
	return nil
}

func (r *repository) CreateList(userId string, errandId string, files []*File) error {
	for _, file := range files {
		err := r.Create(userId, errandId, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) UploadCategoryIcon(adminId string, file *File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	path := fmt.Sprintf("%s/%s/%s", file.Folder, adminId, file.Header.Filename)
	storageWriter := r.Bucket(categoryBucketName).Object(path).NewWriter(ctx)

	headerFile, err := file.Header.Open()
	if err != nil {
		return err
	}

	if _, err = io.Copy(storageWriter, headerFile); err != nil {
		return err
	}

	if err = storageWriter.Close(); err != nil {
		return err
	}

	file.UploadedUrl = bucketUrl + "/" + categoryBucketName + "/" + path
	return nil
}
