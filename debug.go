package debug

import (
	"fmt"
	"os"
)

var enabled bool = len(os.Getenv("DEBUG")) > 0

type sprintfType func(format string, a ...interface{}) string

var sprintfFunc = func(format string, a ...interface{}) {
	return fmt.Sprintf(format, a)
}

func Debug(namespace string) sprintfType {
	
	if !enabled {
		return func(format string, a ...interface{}) {}
	}
	
	sprintfFunc
}

type Debugger struct {
	
}

func (d Debugger) Print2() {
	fmt.Println("Hello")
}