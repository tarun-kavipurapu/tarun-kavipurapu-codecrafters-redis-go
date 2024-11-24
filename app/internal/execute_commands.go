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
type Store struct {
	kv map[string]string
	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		kv: make(map[string]string),
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
		outputString = []byte("-ERR unknown command '" + output.Cmd + "'\r\n")
	}

	return outputString, err

}
func executePing(output *Command) ([]byte, error) {
	return encodeSimpleString("+" + "PONG" + "\r\n"), nil

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
		return nil, fmt.Errorf("Either the key or value is missing")
	}
	key := output.Args[0]

	value := output.Args[1]
	s.kv[key] = value
	if len(output.Args) > 3 && output.Args[2] == "PX" {
		expiry, err := strconv.Atoi(output.Args[3])
		if err != nil {
			return nil, fmt.Errorf("Error Getting the Expiry form the Args")
		}
		go deleteAfterExpiry(expiry, s, key)
	}

	// s.kv[output.Cmd] =

	return respOK, nil
}

func deleteAfterExpiry(t int, s *Store, key string) {
	//waiit till the time is tickered
	time.Sleep(time.Duration(t) * time.Millisecond)
	delete(s.kv, key)
	//delete from the key
}

func executeGet(output *Command, s *Store) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(output.Args) < 1 {
		return nil, fmt.Errorf("KEy is missing")
	}

	key := output.Args[0]
	value, ok := s.kv[key]
	if !ok {
		return respNull, fmt.Errorf("Value corresponding to key not found")
	}
	return respString(value), nil
}
