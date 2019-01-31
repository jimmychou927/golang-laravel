package print

import (
	"encoding/json"
	"fmt"
)

func Msg(args ...interface{}) {

	message := ""

	for _, value := range args {
		s, err := json.MarshalIndent(value, "", "   ")

		if err != nil {
			message += fmt.Sprintln(err)
		} else {
			message += fmt.Sprintln(string(s))
		}
	}

	panic("<pre>\n" + message + "\n</pre>")
}

func Err(args ...interface{}) {

	panic(fmt.Sprintln(args))
}

func Errf(format string, a ...interface{}) {

	panic(fmt.Sprintf(format, a...))
}
