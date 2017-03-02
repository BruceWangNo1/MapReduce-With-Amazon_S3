package mapreduce

import (
	"os"
	"encoding/json"
	"io"
	"log"
)

// doReduce does the job of a reduce worker: it reads the intermediate 
// key/value pairs (produced by the map phase) for this task, sorts the 
// intermediate key/value pairs by key, calls the user-defined reduce function 
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	kvMap := make(map[string][]string) // pair (Key, []Values)
	for m := 0; m < nMap; m++ {
		fi, err := os.Open(reduceName(jobName, m, reduceTaskNumber))
		if err != nil {
			log.Fatal("doReduce 2: ", err)
		}
		defer fi.Close() // this will run just before the function finishes

		// Decoder
		dec := json.NewDecoder(fi)
		// Decode
		for {
			var kv KeyValue
			if err := dec.Decode(&kv); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			kvMap[kv.Key] = append(kvMap[kv.Key], kv.Value)
		}
	}

	// Create merge file.
	mergeFile, err := os.Create(mergeName(jobName, reduceTaskNumber))
	if err != nil {
		log.Fatal("doReduce 1: ", err)
	}
	defer mergeFile.Close()

	// Write merge file
	enc := json.NewEncoder(mergeFile)
	for key, value := range kvMap {
		enc.Encode(KeyValue{key, reduceF(key, value)})
	}
}