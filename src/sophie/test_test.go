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
	"strconv"
	"bufio"
	"sort"
	"time"
)

const (
	nNumber = 100000
	nMap = 10
	nReduce = 5
)

func MapFunc(file string, value string) (res []KeyValue) {
	//debug("Map %v\n", value)
	words := strings.Fields(value)
	for _, w := range words {
		kv := KeyValue{w, ""}
		res = append(res, kv)
	}
	return
}

// Just return key
func ReduceFunc(key string, values []string) string {
	// for _, e := range values {
	// 	debug("Reduce %s %v\n", key, e)
	// }
	return ""
}

func TestDownloadFile(t *testing.T) {
	file, err := os.Create("download_file.txt")
	if err != nil {
	    log.Fatal("Failed to create file", err)
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	numBytes, err := downloader.Download(file,
	    &s3.GetObjectInput{
	        Bucket: aws.String("mybucketbruce"),
	        Key:    aws.String("upload.txt"),
	    })
	if err != nil {
	    fmt.Println("Failed to download file", err)
	    return
	}

	fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")
}
func TestUploadFileStream(t *testing.T) {
	file, err := os.Open("upload.txt")
	if err != nil {
		log.Fatal("failed to open file", err)
	}

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

func TestListObjects(t *testing.T) {
	fmt.Println("list objects in a bucket")
	svc := s3.New(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	result, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("mybucketbruce"),
		Prefix: aws.String("pg"),
		})
	if err != nil {
		log.Println("Failed to list objects", err)
		return
	}

	log.Println("Buckets:")
	fmt.Println(aws.StringValue(result.Contents[0].Key))
	for _, object := range result.Contents {
		log.Printf("%s\n", aws.StringValue(object.Key))
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
	mr := setup()
	for i := 0; i < 2; i++ {
		go RunWorker(mr.Address, port("worker"+strconv.Itoa(i)),
			MapFunc, ReduceFunc, -1)
	}
	mr.Wait()
	check(t, mr.Files)
	checkWorker(t, mr.Stats)
	cleanup(mr)
}

func TestOneFailure(t *testing.T) {
	mr := setup()
	// start 2 workers that fail after 10 tasks
	go RunWorker(mr.Address, port("worker"+strconv.Itoa(0)),
		MapFunc, ReduceFunc, 10)
	go RunWorker(mr.Address, port("worker"+strconv.Itoa(1)),
		MapFunc, ReduceFunc, -1)
	mr.Wait()
	check(t, mr.Files)
	checkWorker(t, mr.Stats)
	cleanup(mr)
}
func TestManyFailures(t *testing.T) {
	mr := setup()
	i := 0
	done := false
	for !done {
		select {
		case done = <-mr.DoneChannel:
			check(t, mr.Files)
			cleanup(mr)
			break
		default:
			// Start 2 workers each sec. The workers fail after 10 tasks
			w := port("worker" + strconv.Itoa(i))
			go RunWorker(mr.Address, w, MapFunc, ReduceFunc, 10)
			i++
			w = port("worker" + strconv.Itoa(i))
			go RunWorker(mr.Address, w, MapFunc, ReduceFunc, 10)
			i++
			time.Sleep(1 * time.Second)
		}
	}
}


func setup() *Master {
	files := makeInputs(nMap)
	master := port("master")
	mr := Distributed("test", files, nReduce, master)
	return mr
}

// make input file
func makeInputs(num int) []string {
	var names []string
	var i = 0
	for f := 0; f < num; f++ {
		names = append(names, fmt.Sprintf("824-mrinput-%d.txt", f))
		file, err := os.Create(names[f])
		if err != nil {
			log.Fatal("mkInput: ", err)
		}
		w := bufio.NewWriter(file)
		for i < (f+1)*(nNumber/num) {
			fmt.Fprintf(w, "%d\n", i)
			i++
		}
		w.Flush()
		
		writeToS3(file.Name())
		file.Close()
	}
	return names 
}

func port(suffix string) string {
	s := "/var/tmp/824-"
	s += strconv.Itoa(os.Getuid()) + "/"
	os.Mkdir(s, 0777)
	s += "mr"
	s += strconv.Itoa(os.Getpid()) + "-"
	s += suffix
	return s
}

// Check input file against output file: each input number should show up
// in the output file in string sorted order
func check(t *testing.T, files []string) {
	output, err := readFromS3("mrtmp.test")
	if err != nil {
		log.Fatal("check: ", err)
	}
	defer output.Close()

	var lines []string
	for _, f := range files {
		input, err := os.Open(f)
		if err != nil {
			log.Fatal("check: ", err)
		}
		defer input.Close()
		inputScanner := bufio.NewScanner(input)
		for inputScanner.Scan() {
			lines = append(lines, inputScanner.Text())
		}
	}

	sort.Strings(lines)

	outputScanner := bufio.NewScanner(output)
	i := 0
	for outputScanner.Scan() {
		var v1 int
		var v2 int
		text := outputScanner.Text()
		n, err := fmt.Sscanf(lines[i], "%d", &v1)
		if n == 1 && err == nil {
			n, err = fmt.Sscanf(text, "%d", &v2)
		}
		if err != nil || v1 != v2 {
			t.Fatalf("line %d: %d != %d err %v\n", i, v1, v2, err)
		}
		i++
	}
	if i != nNumber {
		t.Fatalf("Expected %d lines in output\n", nNumber)
	}
}

// Workers report back how many RPCs they have processed in the Shutdown reply.
// Check that they processed at least 1 RPC.
func checkWorker(t *testing.T, l []int) {
	for _, tasks := range l {
		if tasks == 0 {
			t.Fatalf("Some worker didn't do any work\n")
		}
	}
}

func cleanup(mr *Master) {
	mr.CleanupFiles()
	for _, f := range mr.Files {
		removeFile(f)
	}
}