package log

import (
	// "bytes"
	"log"
	"os"
)

func Logger() *log.Logger {
	// buf := bytes.Buffer{}
	return log.New(os.Stdout, "boo-http:", log.Lshortfile)
}
