package stdio

import (
	"fmt"
	"io"
	"os"
	"sync"

	syncpool "github.com/anacpe9/sync-pool"
)

var (
	out    = io.Writer(os.Stdout)
	err    = io.Writer(os.Stderr)
	StdOut = NewWriter(&out)
	StdErr = NewWriter(&err)
)

type Writer struct {
	io.Writer
	// writer *os.File
}

type writerPayload struct {
	writer *io.Writer
	data   *string
}

var isInit = false
var logEnd chan bool
var logChannel chan *writerPayload
var logPool *syncpool.Pool[writerPayload] // *sync.Pool
var mutex sync.Mutex

func loopLog() {
	for {
		log, ok := <-logChannel
		// fmt.Println(">>> log:", ok)
		if !ok {
			break
		}

		fmt.Fprint(*log.writer, *log.data)
		logPool.Put(log)
	}

	logEnd <- true
}

func InitWriter() {
	if isInit {
		return
	}

	isInit = true
	logEnd = make(chan bool)
	logChannel = make(chan *writerPayload)
	logPool = syncpool.GetPool[writerPayload]()
	// logPool = &sync.Pool{
	// 	New: func() interface{} {
	// 		return new(writerPayload)
	// 	},
	// }

	go loopLog()
}

func WaitLoggerUntilEnd() chan bool {
	defer close(logEnd)
	async := make(chan bool, 1)

	go func(c chan bool) {
		<-logEnd
		c <- true
	}(async)

	close(logChannel)

	block := <-async
	async <- block

	return async
}

func NewWriter(file *io.Writer) io.Writer {
	w := new(Writer)
	// w.writer = file
	w.Writer = *file //io.Writer(file)

	return w
}

func (w *Writer) Write(data []byte) (int, error) {
	mutex.Lock()
	log := logPool.Get() // .(*writerPayload)
	mutex.Unlock()

	dst := string(data)
	log.writer = &w.Writer
	log.data = &dst

	logChannel <- log

	return 0, nil
}
