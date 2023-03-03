package logger

import (
	"fmt"
	"os"
)

func Log(message string, args ...interface{}) {
	fmt.Printf("tailproxy: "+message, args...)
}
func Err(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "tailproxy: "+message, args...)
}
