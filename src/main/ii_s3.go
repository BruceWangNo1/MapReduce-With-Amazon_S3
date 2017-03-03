package main

import (
	"fmt"
	"sophie"
	"os"
	"strings"
	"unicode"
	"strconv"
	"bytes"
	"sort"
)

func SliceUnique(slice []string) (uniqueslice []string) {
	maps := map[string]bool{}
	for _,item := range slice {
		if _,exists := maps[item]; !exists {
			maps[item] = true
			uniqueslice = append(uniqueslice,item)
		}
	}
	return uniqueslice
}

// The mapping function is called once for each piece of the input.
// In this framework, the key is the name of the file that is being processed,
// and the value is the file's contents. The return value should be a slice of
// key/value pairs, each represented by a sophie.KeyValue.
func mapF(document string, value string) (res []sophie.KeyValue) {
	// TODO: you should complete this to do the inverted index challenge
	f := func (c rune) bool  {
        return !unicode.IsLetter(c) && !unicode.IsNumber(c)
    }

	words := strings.FieldsFunc(value, f)
	unique := SliceUnique(words)

    for _, w := range unique {
        res = append(res, sophie.KeyValue{Key:w,Value:document})
    }

    return
}

// The reduce function is called once for each key generated by Map, with a
// list of that key's string value (merged across all inputs). The return value
// should be a single output value for that key.
func reduceF(key string, values []string) string {
	// TODO: you should complete this to do the inverted index challenge
	var buffer bytes.Buffer  // 用作字符串拼接

	sort.Strings(values)
	length := len(values)

	for index,name := range values{
		buffer.WriteString(name)
		if index != (length-1){
			buffer.WriteString(",")
		}
	}
	return strconv.Itoa(length) + " " + buffer.String()

}

// Can be run in 3 ways:
// 1) Sequential (e.g., go run wc.go master sequential x1.txt .. xN.txt)
// 2) Master (e.g., go run wc.go master localhost:7777 x1.txt .. xN.txt)
// 3) Worker (e.g., go run wc.go worker localhost:7777 localhost:7778 &)
func main() {
	if len(os.Args) < 4 {
		fmt.Printf("%s: see usage comments in file\n", os.Args[0])
	} else if os.Args[1] == "master" {
		var mr *sophie.Master
		if os.Args[2] == "sequential" {
			mr = sophie.Sequential("iiseq", sophie.getKeys("pg"), 3, mapF, reduceF) // os.Args[3:]
		} else {
			mr = sophie.Distributed("iiseq", os.Args[3:], 3, os.Args[2])
		}
		mr.Wait()
	} else {
		sophie.RunWorker(os.Args[2], os.Args[3], mapF, reduceF, 100)
	}
}