package internal

import "fmt"

type Resp struct {
}

var respNull = []byte("$-1\r\n")
var respOK = []byte("+OK\r\n")

// $<length>\r\n<value>\r\n

func respString(value string) []byte {
	outputString := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)

	return ([]byte(outputString))
}

func encodeSimpleString(s string) []byte {

	return []byte("+" + s + "\r\n")

}
