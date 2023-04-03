package logger

import (
	"fmt"
	"os"
)

func Log(message string, args ...interface{}) {
	fmt.Printf("tailproxy: "+message+"\n", args...)
}
func Err(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "tailproxy: "+message+"\n", args...)
}
func Fatal(message string, args ...interface{}) {
	Err(message, args...)
	os.Exit(1)
}
