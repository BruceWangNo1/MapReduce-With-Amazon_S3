package mapreduce

import (
	"hash/fnv"
	"os"
	"fmt"
	"encoding/json"
)

// doMap does the job of a map worker: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce tasks that will be run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	file, err := os.Open(inFile)
	if err == nil {
		fmt.Printf("file: %s opened\n", inFile)
	} else {
		fmt.Println(err)
	}
	inf, err := file.Stat()

	contents := make([]byte, inf.Size())
	file.Read(contents)
	file.Close()

	kv := mapF(inFile, string(contents))
	filesenc := make([]*json.Encoder, nReduce)
	files := make([]*os.File, nReduce)

	for i := range filesenc {
		file, err := os.Create(reduceName(jobName, mapTaskNumber, i))
		if err != nil {
			fmt.Printf("%s Create Failed\n", reduceName(jobName, mapTaskNumber, i))
		} else {
			filesenc[i] = json.NewEncoder(file)
			files[i] = file
		}
	}

	for _, v := range kv {
		err := filesenc[ihash(v.Key) % uint32(nReduce)].Encode(&v)
		if err != nil {
			fmt.Printf("%s Encode Failed %v\n", v, err)
		}
	}

	for _, f := range files {
		f.Close()
	}
}

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}