package print

import "fmt"

func Ln(args ...interface{}) {

	panic("<p/>" + fmt.Sprintln(args))
}

func Err(args ...interface{}) {

	panic(fmt.Sprintln(args))
}

func Errf(format string, a ...interface{}) {

	panic(fmt.Sprintf(format, a...))
}
