package sophie

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	//"os"
	"sync"
	"os"
)

// Worker holds the state for a server waiting for DoTask or Shutdown RPCs
type Worker struct {
	sync.Mutex

	name string
	Map func(string, string) []KeyValue
	Reduce func(string, []string) string
	nRPC int //protected by mutex
	nTasks int // protected by mutex
	l net.Listener
	files []string // files used. record every file used and remove all later
}

// DoTask is called by the master when a new task is being scheduled on this 
// worker.
func (wk *Worker) DoTask(arg *DoTaskArgs, _ *struct{}) error {
	fmt.Printf("%s: given %v task #%d on file %s (nios: %d)\n",
		wk.name, arg.Phase, arg.TaskNumber, arg.File, arg.NumOtherPhase)

	switch arg.Phase {
	case mapPhase:
		// wk.files = append(wk.files, arg.File) // keep track of all the files downloaded and saved locally
		doMap(arg.JobName, arg.TaskNumber, arg.File, arg.NumOtherPhase, wk.Map)
	case reducePhase:
		doReduce(arg.JobName, arg.TaskNumber, arg.NumOtherPhase, wk.Reduce)
	}

	fmt.Printf("%s: %v task #%d done\n", wk.name, arg.Phase, arg.TaskNumber)
	return nil
}

// Shutdown is called by the master when all work has been completed.
// We should respond with the number of tasks we have processed.
func (wk *Worker) Shutdown(_ *struct{}, res *ShutdownReply) error {
	debug("Shutdown %s\n", wk.name)
	wk.Lock()
	defer wk.Unlock()
	res.Ntasks = wk.nTasks
	wk.nRPC = 1
	wk.nTasks-- // Don't count the shutdown RPC

	wk.l.Close()
	// debug("RunWorker %s exit\n", wk.name)
	os.Remove(wk.name)

	return nil
}

// Tell the master we exist and ready to work
func (wk *Worker) register(master string) {
	args := new(RegisterArgs)
	args.Worker = wk.name

	ok := false
	for ok == false {
		ok = call(master, "Master.Register", args, new(struct{}))
	}
}

// RunWorker sets up a connection with the master, registers its address, and 
// waits for tasks to be scheduled.
func RunWorker(MasterAddress string, me string,
	MapFunc func(string, string) []KeyValue,
	ReduceFunc func(string, []string) string, 
	nRPC int, 
) {
	debug("RunWorker %s\n", me)
	wk := new(Worker)
	wk.name = me
	wk.Map = MapFunc
	wk.Reduce = ReduceFunc
	wk.nRPC = nRPC
	rpcs := rpc.NewServer()
	rpcs.Register(wk)
	os.Remove(me) // only need for "unix"
	// l, e := net.Listen("unix", me)
	l, e := net.Listen("tcp", me)
	if e != nil {
		log.Fatal("RunWorker: worker ", me, " error: ", e)
	}
	wk.l = l
	fmt.Println("worker is preparing to register with master")
	wk.register(MasterAddress)
	fmt.Println("registration done")

	// DON'T MODIFY CODE BELOW
	for {
		wk.Lock()
		fmt.Printf("NOTICE: wk.nRPC is %d\n", wk.nRPC)
		if wk.nRPC == 0 {
			wk.Unlock()
			break
		}
		wk.Unlock()
		conn, err := wk.l.Accept()
		if err == nil {
			wk.Lock()
			wk.nRPC--
			wk.Unlock()
			go rpcs.ServeConn(conn)
			wk.Lock()
			wk.nTasks++
			wk.Unlock()
		} else {
			break
		}
	}
	//wk.l.Close()
	debug("RunWorker %s exit\n", me)
}