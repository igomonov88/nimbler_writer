package bucket

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Store is used for storing file in provided object storage with key as a filename and
// body as a file content.
func Store(ctx context.Context, session *session.Session, bucketName, key, body string) error {
	ctx, span := trace.StartSpan(ctx, "internal.object_storage.Store")
	defer span.End()

	fileName := fmt.Sprintf("%v.txt", key)
	file, err := ioutil.TempFile("", fileName)
	if err != nil {
		return errors.Wrapf(err, "creating a temp file with name: %v", fileName)
	}

	if err := ioutil.WriteFile(file.Name(), []byte(body), 0644); err != nil {
		return errors.Wrapf(err, "failed to write data to file.")
	}

	defer os.Remove(fileName)

	uploader := s3manager.NewUploader(session)
	input := s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	}

	if _, err = uploader.UploadWithContext(ctx, &input); err != nil {
		return errors.Wrapf(err, "upload file: %v to bucket: %v", fileName, bucketName)
	}

	if err := file.Close(); err != nil {
		return errors.Wrapf(err, "closing file: %v", fileName)
	}

	return nil
}

// Retrieve is retrieve content of file with key as file name from provided storage.
func Retrieve(ctx context.Context, session *session.Session, bucket, key string) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "internal.object_storage.Retrieve")
	defer span.End()

	file, err := ioutil.TempFile("", "temp.txt")
	if err != nil {
		return nil, errors.Wrapf(err, "creating file.")
	}

	defer os.Remove("temp.txt")

	downloader := s3manager.NewDownloader(session)

	input := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	if _, err := downloader.DownloadWithContext(ctx, file, &input); err != nil {
		return nil, errors.Wrapf(err, "download file with key: %v from bucket.", key)
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(file); err != nil {
		return nil, errors.Wrap(err, "read from file.")
	}

	return buf.Bytes(), nil
}

// Delete is deleting file with provided key from bucket.
func Delete(ctx context.Context, session *session.Session, bucket, key string) error {
	ctx, span := trace.StartSpan(ctx, "internal.object_storage.Delete")
	defer span.End()

	svc := s3.New(session)
	input := s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}
	if _, err := svc.DeleteObject(&input); err != nil {
		return errors.Wrapf(err, "delete file: %v from bucket", key)
	}

	return nil
}
