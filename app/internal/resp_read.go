package internal

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)
const (
	CRLF = "\r\n"
)
const (
	MaxBulkSize = 512 * 1024 * 1024 // 512MB max bulk size
	MaxArrayLen = 1024 * 1024       // 1M max array elements
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
	byte, err := r.c.ReadByte()
	if err != nil {
		return nil, err
	}

	switch byte {
	case STRING:
		return r.readSimpleString()
	case ERROR:
		return r.readError()
	case INTEGER:
		return r.readInteger()
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		return nil, fmt.Errorf("unknown RESP type: %c", byte)
	}
}

// Add these new methods
func (r *RespReader) readSimpleString() (string, error) {
	line, err := r.readLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}

func (r *RespReader) readError() (error, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	return fmt.Errorf(string(line)), nil
}

// Helper method for reading lines
func (r *RespReader) readLine() ([]byte, error) {
	line, _, err := r.c.ReadLine()
	return line, err
}

func (r *RespReader) readBulk() (string, error) {
	length, err := r.readInteger()
	if err != nil {
		return "", err
	}

	if length == -1 {
		return "", nil // null bulk string
	}

	// Read exactly length bytes
	bulk := make([]byte, length)
	_, err = io.ReadFull(r.c, bulk)
	if err != nil {
		return "", err
	}

	// Read and discard CRLF
	_, _, err = r.c.ReadLine()
	if err != nil {
		return "", err
	}

	return string(bulk), nil
}
func (r *RespReader) readInteger() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}

	num, err := strconv.Atoi(string(line))
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %s", line)
	}

	return num, nil
}
func (r *RespReader) readArray() (interface{}, error) {
	length, err := r.readInteger()
	if err != nil {
		return nil, err
	}

	if length == -1 {
		return nil, nil // null array
	}

	values := make([]interface{}, length)
	for i := 0; i < length; i++ {
		val, err := r.CommandRead()
		if err != nil {
			return nil, err
		}
		values[i] = val
	}

	return values, nil
}
