package bucket_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/igomonov88/nimbler_writer/config"
	"github.com/igomonov88/nimbler_writer/internal/bucket"
	"github.com/igomonov88/nimbler_writer/internal/tests"
)

func TestObjectStorage(t *testing.T) {
	// read config file to connect to test bucket
	cfg, err := config.Parse("config.yaml")
	if err != nil {
		t.Fatalf("failed to read config.")
	}

	// configure creds to connect to bucket
	creds := credentials.NewStaticCredentials(cfg.S3.AccessKeyID, cfg.S3.SecretKey, "")
	awsCfg := aws.Config{Region: aws.String("eu-north-1"), Credentials: creds}

	// create a session to use in bucket calls.
	sess, err := session.NewSession(&awsCfg)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	t.Log("Given the need to test bucket functionality:")

	{
		if err := bucket.Store(context.Background(), sess, cfg.S3.BucketName, "my_key", "my_body"); err != nil {
			t.Fatalf("\t%s\tShould be able to add new file to bucket: %s", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to add new file to bucket.", tests.Success)

		body, err := bucket.Retrieve(context.Background(), sess, cfg.S3.BucketName, "my_key")
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve file from bucket: %s", tests.Failed, err)
		}

		if string(body) != "my_body" {
			t.Fatalf("\t%s\tShould be able to retrieve the same file from bucket: %s", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to retrieve file from bucket.", tests.Success)

		if err := bucket.Delete(context.Background(), sess, cfg.S3.BucketName, "my_key"); err != nil {
			t.Fatalf("\t%s\tShould be able to delete file from bucket: %s", tests.Failed, err)
		}

		t.Logf("\t%s\tShould be able to delete file from bucket.", tests.Success)

	}
}
