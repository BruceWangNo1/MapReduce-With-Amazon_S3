package sophie

import (
	"testing"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"fmt"
	"log"
	"strings"
	//"compress/gzip"
	//"io"
	"os"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func TestUploadFileStream(t *testing.T) {
	file, err := os.Open("upload.txt")
	if err != nil {
		log.Fatal("failed to open file", err)
	}

	// reader, writer := io.Pipe()
	// go func() {
	// 	gw := gzip.NewWriter(writer)
	// 	io.Copy(gw, file)

	// 	file.Close()
	// 	gw.Close()
	// 	writer.Close()
	// }()

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body: file,
		Bucket: aws.String("mybucketbruce"),
		Key: aws.String("upload.txt"),
		})
	if err != nil {
		log.Fatalln("Failed to upload", err)
	}

	log.Println("Successfully uploaded to", result.Location)
}

func TestListBuckets(t *testing.T) {
	fmt.Println("list buckets")
	svc := s3.New(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
	    log.Println("Failed to list buckets", err)
	    return
	}

	log.Println("Buckets:")
	for _, bucket := range result.Buckets {
	    log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
	}
}

func TestCreateBucket(t *testing.T) {
	bucket := "mybucketbrucelee"
	key := "TestFile.txt"

	svc := s3.New(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
	    Bucket: &bucket,
	})
	if err != nil {
	    log.Println("Failed to create bucket", err)
	    return
	}

	if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucket}); err != nil {
	    log.Printf("Failed to wait for bucket to exist %s, %s\n", bucket, err)
	    return
	}

	_, err = svc.PutObject(&s3.PutObjectInput{
	    Body:   strings.NewReader("Hello World!"),
	    Bucket: &bucket,
	    Key:    &key,
	})
	if err != nil {
	    log.Printf("Failed to upload data to %s/%s, %s\n", bucket, key, err)
	    return
	}

	log.Printf("Successfully created bucket %s and uploaded data with key %s\n", bucket, key)
}

func TestBasic(t *testing.T) {
}