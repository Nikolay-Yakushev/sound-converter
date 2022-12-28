package loggers

import (
	"log"
	"os"
)

var (
	InfLogger *log.Logger
	ErrLogger *log.Logger
)

func init() {
    InfLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags | log.Lshortfile | log.Ldate|log.Ltime|log.Lshortfile)
	ErrLogger = log.New(os.Stdout, "Error: ",log.LstdFlags | log.Lshortfile | log.Ldate|log.Ltime|log.Lshortfile)
}