package sophie

import (
	"fmt"
	"net"
	"sync"
	"os/user"
	"time"
)

// Master holds all the states that the master needs to keep track off
type Master struct {
	sync.Mutex

	Address string
	RegisterChannel chan string
	DoneChannel chan bool
	Workers []string //protected by mutex

	// Per-task information
	JobName string //name of the currently executing job
	Files []string //input files
	NReduce int //number of reduce partitions

	shutdown chan struct{}
	l net.Listener
	Stats []int

	// added for panel package
	User string
	StartTime time.Time
	SchedulingMode string
	ElapsedTime float64

}

// Register is an RPC method that is called by workers after they have started 
// up to report that they are ready to receive tasks.
func (mr *Master) Register(args *RegisterArgs, _ *struct{}) error {
	mr.Lock()
	defer mr.Unlock()
	debug("Register: worker %s\n", args.Worker)
	mr.Workers = append(mr.Workers, args.Worker)
	go func() {
		mr.RegisterChannel <- args.Worker
	}()
	return nil
}

// newMaster initialize a new Map/Reduce Master
func newMaster(master string) (mr *Master) {
	mr = new(Master)
	mr.Address = master
	mr.shutdown = make(chan struct{})
	mr.RegisterChannel = make(chan string)
	mr.DoneChannel = make(chan bool)

	// added for panel package
	userStruct, err := user.Current()
	if err == nil {
		mr.User = userStruct.Username
	}
	mr.StartTime = time.Now()
	mr.SchedulingMode = "No Mode"

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
				doMap(mr.JobName, i, f, mr.NReduce, mapF)
			}
		case reducePhase:
			for i := 0; i < mr.NReduce; i++ {
				doReduce(mr.JobName, i, len(mr.Files), reduceF)
			}
		}
		}, func() {
			mr.Stats = []int{len(files) + nreduce}
			})
	return
}

//Distributed schedules map and reduce tasks on workers that register with the 
// master over RPC
func Distributed(jobName string, files []string, nreduce int, master string) (mr *Master) {
	mr = newMaster(master)
	mr.startRPCServer()
	go mr.run(jobName, files, nreduce, mr.schedule, func() {
		mr.Stats = mr.killWorkers()
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
	mr.JobName = jobName
	mr.Files = files
	mr.NReduce = nreduce

	fmt.Printf("%s: starting Map/Reduce task %s\n", mr.Address, mr.JobName)

	schedule(mapPhase)
	schedule(reducePhase)
	finish()
	mr.merge()

	fmt.Printf("%s: Map/Reduce task completed\n", mr.Address)

	mr.DoneChannel <- true
}

// Wait blocks until the currently scheduled work has completed.
// This happens when all tasks have scheduled and completed, the final output
// have been computed, and all workers have been shut down.
func (mr *Master) Wait() {
	<-mr.DoneChannel
}

// killWorkers cleans up all workers by sending each one a Shutdown RPC.
// It also collects and returns the number of tasks each worker has performed.
func (mr *Master) killWorkers() []int {
	mr.Lock() // lock to ensure exclusive access
	defer mr.Unlock()
	ntasks := make([]int, 0, len(mr.Workers))
	for _, w := range mr.Workers {
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