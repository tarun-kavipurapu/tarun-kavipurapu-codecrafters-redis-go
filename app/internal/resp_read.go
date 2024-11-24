package internal

import (
	"bufio"
	"io"
	"log"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type RespReader struct {
	c *bufio.Reader
}

func NewRespReader(r io.Reader) *RespReader {
	return &RespReader{
		c: bufio.NewReader(r),
	}
}

func (r *RespReader) CommandRead() (interface{}, error) {
	// input := "$5\r\nAhmed\r\n"

	var err error
	byte, err := r.c.ReadByte()
	if err != nil {
		return nil, err
	}
	switch byte {
	case STRING:
	case ERROR:
	case INTEGER:
		return r.readInteger()
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		log.Println("Unknown Type\n")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *RespReader) readBulk() (string, error) {
	n, err := r.readInteger()
	if err != nil {
		return "", err
	}
	if n == -1 {
		return "", nil
	}
	command, _, err := r.c.ReadLine()
	if err != nil {
		return "", err
	}
	// log.Println(string(command))

	return string(command), nil
}

func (r *RespReader) readInteger() (int, error) {
	byte, _, err := r.c.ReadLine()
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(string(byte))
	if err != nil {
		return 0, nil
	}

	return num, nil

}

func (r *RespReader) readArray() (interface{}, error) {
	size, err := r.readInteger()
	if err != nil {
		return nil, err
	}
	values := make([]interface{}, size)
	for i := 0; i < size; i++ {
		val, err := r.CommandRead()
		log.Println(val)
		if err != nil {
			return nil, err
		}
		values[i] = val

	}
	return values, nil
}
