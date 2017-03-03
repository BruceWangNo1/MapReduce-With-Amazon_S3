package sophie

import (
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"fmt"
	"log"
	//"strings"
	"os"
)

// read file from amazon s3 to local disk. content extracted and file deleted locally.
func readFromS3(infile string) (*os.File, error) {
	file, err := os.Create(infile)
	if err != nil {
		log.Fatalln("Failed to create file", err)
	}

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	numBytes, err := downloader.Download(file, 
		&s3.GetObjectInput{
			Bucket: aws.String("mybucketbruce"),
			Key: aws.String(infile),
			})
	if err != nil {
		log.Fatalln("Failed to download file", err)
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")

	return file, err
}

// upload file to amazon s3 and delete local file.
func writeToS3(infile string) {
	file, err := os.Open(infile)
	if err != nil {
		log.Fatal("failed to open file", err)
	}

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body: file,
		Bucket: aws.String("mybucketbruce"),
		Key: aws.String(infile),
		})
	if err != nil {
		log.Fatalln("Failed to upload", err)
	}

	log.Println("Successfully uploaded to", result.Location)
}


// // upload file to amazon s3 and delete local file.
// func writeToS3(infile string) {
// 	file, err := os.Open(infile)
// 	if err != nil {
// 		log.Fatal("failed to open file", err)
// 	}

// 	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
// 	result, err := uploader.Upload(&s3manager.UploadInput{
// 		Body: file,
// 		Bucket: aws.String("mybucketbruce"),
// 		Key: aws.String(infile),
// 		})
// 	if err != nil {
// 		log.Fatalln("Failed to upload", err)
// 	}

// 	log.Println("Successfully uploaded to", result.Location)

// 	removeFile(infile)
// }

// // read file from amazon s3 to local disk. content extracted and file deleted locally.
// func readFromS3(infile string) (contents []byte) {
// 	file, err := os.Create(infile)
// 	if err != nil {
// 		log.Fatalln("Failed to create file", err)
// 	}

// 	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
// 	numBytes, err := downloader.Download(file, 
// 		&s3.GetObjectInput{
// 			Bucket: aws.String("mybucketbruce"),
// 			Key: aws.String(infile),
// 			})
// 	if err != nil {
// 		fmt.Println("Failed to download file", err)
// 		return
// 	}

// 	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")

// 	inf, err := file.Stat()
// 	contents := make([]byte, inf.Size())
// 	file.Read(contents)
// 	file.Close()

// 	removeFile(infile)

// 	return contents
// }