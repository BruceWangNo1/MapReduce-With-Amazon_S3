package sophie

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

// merge combines the results of the many reduce jobs into a single output file
// XXX use merge sort
func (mr *Master) merge() {
	debug("Merge phase")
	kvs := make(map[string]string)
	for i := 0; i < mr.NReduce; i++ {
		p := mergeName(mr.JobName, i)
		fmt.Printf("Merge: read %s\n", p)
		file, err := readFromS3(p)
		if err != nil {
			log.Fatalln("do Reduce 2: ", err)
		}

		dec := json.NewDecoder(file)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err != nil {
				break
			}
			kvs[kv.Key] = kv.Value
		}
		file.Close()
		removeFile(p)
	}
	var keys []string
	for k := range kvs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	file, err := os.Create("mrtmp." + mr.JobName)
	if err != nil {
		log.Fatal("Merge: create ", err)
	}
	w := bufio.NewWriter(file)
	for _, k := range keys {
		fmt.Fprintf(w, "%s: %s\n", k, kvs[k])
	}
	w.Flush()
	
	writeToS3(file.Name())
}

// removeFile is a simple wrapper around os.Remove that logs errors.
func removeFile(n string) {
	err := os.Remove(n)
	if err != nil {
		log.Fatal("CleanupFiles ", err)
	}
}

// CleanupFiles removes all intermediate files produced by running mapreduce.
func (mr *Master) CleanupFiles() {
	for i := range mr.Files {
		for j := 0; j < mr.NReduce; j++ {
			removeFile(reduceName(mr.JobName, i, j))
		}
	}
	for i := 0; i < mr.NReduce; i++ {
		removeFile(mergeName(mr.JobName, i))
	}
	removeFile("mrtmp." + mr.JobName)
}