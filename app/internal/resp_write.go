package internal

import "fmt"

type Resp struct {
}

var respNull = "$-1\r\n"

func respOK() string {
	return ("+OK\r\n")
}

// $<length>\r\n<value>\r\n

func respString(value string) string {
	outputString := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)

	return (outputString)
}
