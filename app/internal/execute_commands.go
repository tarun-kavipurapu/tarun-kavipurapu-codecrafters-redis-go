package internal

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

/*
- Create a channel  where expiry of keys can be stored
- Now i want to concurrently watch that datastructure if an new entries is made into
-
*/
// type KVExpiry struct {
// 	key    string
// 	expiry time.Time
// }
type Record struct {
	value     string
	createdAt time.Time
	expiryAt  time.Time
}
type Store struct {
	kv map[string]*Record
	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		kv: make(map[string]*Record),
		mu: sync.RWMutex{},
	}
}

func EvaluateFunc(output *Command, s *Store) ([]byte, error) {
	var err error
	var outputString []byte
	switch output.Cmd {
	case "PING":
		outputString, err = executePing(output)
	case "ECHO":
		outputString, err = executeEcho(output)
	case "SET":
		outputString, err = executeSet(output, s)
	case "GET":
		outputString, err = executeGet(output, s)
	default:
		err = fmt.Errorf("-ERR unknown command '" + output.Cmd + "'\r\n")
	}

	return outputString, err

}
func executePing(output *Command) ([]byte, error) {
	return encodeSimpleString("PONG"), nil

}
func executeEcho(output *Command) ([]byte, error) {
	ans := output.Args[0]
	return encodeSimpleString(ans), nil

}
func executeSet(output *Command, s *Store) ([]byte, error) {
	s.mu.Lock()

	defer s.mu.Unlock()
	log.Println(output)
	if len(output.Args) < 2 {
		return nil, fmt.Errorf("either the key or value is missing")
	}
	key := output.Args[0]

	value := output.Args[1]
	record := &Record{
		value:     value,
		createdAt: time.Now(),
	}
	log.Println(output.Args)

	if len(output.Args) == 4 && output.Args[2] == "PX" {
		expiration, err := strconv.Atoi(output.Args[3])

		if err != nil {

			fmt.Println(err)

			return nil, fmt.Errorf("-ERR wrong expiration time provided for the record")

		}

		record.expiryAt = time.Now().Add(time.Duration(expiration) * time.Millisecond) // Converting milliseconds to seconds
	}
	s.kv[key] = record

	log.Println("Expiry Time", s.kv[key].expiryAt)

	return respOK, nil
}

func executeGet(output *Command, s *Store) ([]byte, error) {
	// s.mu.RLock()
	// defer s.mu.RUnlock()
	if len(output.Args) < 1 {
		return nil, fmt.Errorf("Key is missing")
	}

	key := output.Args[0]
	val, ok := s.kv[key]
	if ok {
		if val.expiryAt.IsZero() || val.expiryAt.After(time.Now()) {
			return respString(val.value), nil
		}
		// If expired, delete the key and return null
		delete(s.kv, key)
		return respNull, nil
	} else {
		return respNull, nil
	}

}
