package main

import (
	"fmt"
	"os"

	stdio "github.com/anacpe9/std-io"
	logger "github.com/anacpe9/std-io/example/basic/logger"
)

const module = "Main"

var log = logger.NewLogger(module)

func main() {
	fmt.Println("")
	defer func() {
		<-stdio.WaitLoggerUntilEnd()
	}()

	stdio.InitWriter()
	logger.InitLogger(
		&stdio.StdOut,
		&stdio.StdErr,
	)

	logLevel := logger.LOG_LEVEL_DEBUG
	logger.SetLevel(logger.LOG_LEVEL_DEBUG)
	log.Info("===============================")
	log.Info("Started PID   : ", os.Getpid())
	log.Info("Version number: ", "0.0.0")
	log.Info("Version build : ", "2024-05-15")
	log.Info("Environment   : ", os.Getenv("GO_ENV"))
	log.Info("Logger Level  : ", logLevel.ToString())

	for cnt := 0; cnt < 1000000; cnt++ {
		log.Info("Count : ", cnt)
	}

	log.Info("Good bye. Have a nice day.\n")
}
