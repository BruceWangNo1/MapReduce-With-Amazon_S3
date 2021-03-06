package sophie

import (
	"fmt"
	"net/rpc"
)

// What follows are RPC types and methods.
// Field names must start with capital letters, otherwise RPC will break.

// DoTaskArgs holds the arguments that are passed to a worker when a job is 
// scheduled on it.
type DoTaskArgs struct {
	JobName string
	File string // the file to process
	Phase jobPhase // are we in mapPhase or reducePhase?
	TaskNumber int // this task's index in the current phase

	// NumOtherPhase is the total number of tasks in other phase; mappers
	// need this to compute the number of output bins, and reducers needs 
	// this to know how many input files to collect.
	NumOtherPhase int
}

// ShutdownReply is the response to a WorkerShutdown.
// It holds the number of tasks this worker has processed since it was started.
type ShutdownReply struct {
	Ntasks int
}

// RegisterArgs is the argument passed when a worker registers with the master.
type RegisterArgs struct {
	Worker string
}

// call() sends an RPC to the rpcname handler on server server 
// with arguments args, waits for the reply, and leaves the 
// reply in reply. the reply argument should be the address 
// of a reply structure.
//
// call() returns true if the server responded, and false 
// if call() was not able to contact the server. in particular,
// reply's contents are valid if and only if call() returned true.
//
// you should assume that call() will time out and return an 
// error after a while if it does not get a reply from the server.
//
// please use call() to send all RPCs, in master.go, mapreduce.go,
// and worker.go. please do not change this function.
//
func call(srv string, rpcname string,
	args interface{}, reply interface{}) bool {
	// c, errx := rpc.Dial("unix", srv)
	//fmt.Println("send an RPC to the rpcname handler")
	//fmt.Println(srv, rpcname)
	c, errx := rpc.Dial("tcp", srv)
	//fmt.Println("dial done", "err:", errx)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	//fmt.Println("rpc call done")
	if err == nil {
		//fmt.Println("send an RPC to the rpcname handler done")
		return true
	}

	fmt.Println(err)

	return false
}