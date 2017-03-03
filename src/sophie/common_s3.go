package sophie

import (
	"fmt"
	"strconv"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "log"
)

// Debugging enabled?
const debugEnabled = true

// DPrintf will only print if the debugEnabled const has been set to true
func debug(format string, a ...interface{}) (n int, err error) {
	if debugEnabled {
		n, err = fmt.Printf(format, a...)
	}
	return
}

//jobPhase indicates whether a task is scheduled as a map or reduce task.
type jobPhase string

const (
	mapPhase jobPhase = "Map"
	reducePhase jobPhase = "Reduce"
)

// KeyValue is a type used to hold the key/value pairs passed to the map and 
// reduce functions.
type KeyValue struct {
	Key string
	Value string
}

// reduceName constructs the name of the intermediate file which map task
// <mapTask> produces for reduce task <reduceTask>.
func reduceName(jobName string, mapTask int, reduceTask int) string {
	return "mrtmp." + jobName + "-" + strconv.Itoa(mapTask) + "-" + strconv.Itoa(reduceTask)
}

// mergeName constructs the name of the output file of reduce task <reduceTask>
func mergeName(jobName string, reduceTask int) string {
	return "mrtmp." + jobName + "-res-" + strconv.Itoa(reduceTask)
}

func getKeys(prefix string) []string {
	fmt.Println("list objects in a bucket")
	svc := s3.New(session.New(&aws.Config{Region: aws.String("ap-northeast-1")}))
	result, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String("mybucketbruce"),
		Prefix: aws.String(prefix),
		})
	if err != nil {
		log.Println("Failed to list objects", err)
		return make([]string, 0)
	}

	keys := make([]string, len(result.Contents))
	log.Println("Buckets:")
	fmt.Println(aws.StringValue(result.Contents[0].Key))
	for i, object := range result.Contents {
		keys[i] = aws.StringValue(object.Key) 
		//log.Printf("%s\n", aws.StringValue(object.Key))
		log.Printf("%s\n", keys[i])
	}

	return keys
}