package mapreduce

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Master holds all the states that the master needs to keep track off
type Master struct {
	sync.Mutex

	address string
	registerChannel chan string
	doneChannel chan bool
	workers []string //protected by mutex

	// Per-task information
	jobName string //name of the currently executing job
	files []string //input files
	nReduce int //number of reduce partitions

	shutdown chan struct{}
	l net.Listener
	stats []int

	StartTime time.Time
	MapTime float64
	ElapsedTime float64
}

// Register is an RPC method that is called by workers after they have started 
// up to report that they are ready to receive tasks.
func (mr *Master) Register(args *RegisterArgs, _ *struct{}) error {
	mr.Lock()
	defer mr.Unlock()
	debug("Register: worker %s\n", args.Worker)
	mr.workers = append(mr.workers, args.Worker)
	go func() {
		mr.registerChannel <- args.Worker
	}()
	return nil
}

// newMaster initialize a new Map/Reduce Master
func newMaster(master string) (mr *Master) {
	mr = new(Master)
	mr.address = master
	mr.shutdown = make(chan struct{})
	mr.registerChannel = make(chan string)
	mr.doneChannel = make(chan bool)
	return
}
// Sequential runs map and reduce tasks sequentially, waiting for 
// each task to finish before scheduling the next
func Sequential(jobName string, files []string, nreduce int,
	mapF func(string, string) []KeyValue,
	reduceF func(string, []string) string,
) (mr *Master) {
	mr = newMaster("master")
	go mr.run(jobName, files, nreduce, func(phase jobPhase) {
		switch phase {
		case mapPhase:
			for i, f := range files {
				doMap(mr.jobName, i, f, mr.nReduce, mapF)
			}
		case reducePhase:
			for i := 0; i < mr.nReduce; i++ {
				doReduce(mr.jobName, i, len(mr.files), reduceF)
			}
		}
		}, func() {
			mr.stats = []int{len(files) + nreduce}
			})
	return
}

//Distributed schedules map and reduce tasks on workers that register with the 
// master over RPC
func Distributed(jobName string, files []string, nreduce int, master string) (mr *Master) {
	mr = newMaster(master)
	mr.startRPCServer()
	go mr.run(jobName, files, nreduce, mr.schedule, func() {
		mr.stats = mr.killWorkers()
		mr.stopRPCServer()
		})
	return
}
// run executes a mapreduce job on the given number of mappers and reducers
//
func (mr *Master) run(jobName string, files []string, nreduce int,
	schedule func(phase jobPhase),
	finish func(),
) {
	mr.jobName = jobName
	mr.files = files
	mr.nReduce = nreduce

	fmt.Printf("%s: starting Map/Reduce task %s\n", mr.address, mr.jobName)

	mr.StartTime = time.Now()

	schedule(mapPhase)
	mr.MapTime = time.Since(mr.StartTime).Seconds()

	schedule(reducePhase)
	finish()
	mr.merge()

	fmt.Printf("%s: Map/Reduce task completed\n", mr.address)
	mr.ElapsedTime = time.Since(mr.StartTime).Seconds()
	fmt.Printf("Map time is %v s\n", mr.MapTime)
	fmt.Printf("The program finished in %v s\n", mr.ElapsedTime)

	mr.doneChannel <- true
}

// Wait blocks until the currently scheduled work has completed.
// This happens when all tasks have scheduled and completed, the final output
// have been computed, and all workers have been shut down.
func (mr *Master) Wait() {
	<-mr.doneChannel
}

// killWorkers cleans up all workers by sending each one a Shutdown RPC.
// It also collects and returns the number of tasks each worker has performed.
func (mr *Master) killWorkers() []int {
	mr.Lock() // lock to ensure exclusive access
	defer mr.Unlock()
	ntasks := make([]int, 0, len(mr.workers))
	for _, w := range mr.workers {
		debug("Master: shutdown worker %s\n", w)
		var reply ShutdownReply
		ok := call(w, "Worker.Shutdown", new(struct{}), &reply)
		if ok == false {
			fmt.Printf("Master: RPC %s shutdown error\n", w)
		} else {
			ntasks = append(ntasks, reply.Ntasks)
		}
	}
	return ntasks
}