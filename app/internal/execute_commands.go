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
		return nil, fmt.Errorf("either the key or value is missing")
	}
	key := output.Args[0]

	value := output.Args[1]
	record := &Record{
		value:     value,
		createdAt: time.Now(),
	}

	if len(output.Args) > 3 && output.Args[2] == "PX" {
		expiry, err := strconv.Atoi(output.Args[3])
		if err != nil {
			return nil, fmt.Errorf("error getting the expiry from the args")
		}
		log.Println("Expiry Triggered", expiry)

		//Another way  of expiring this would be to maintain an expiry map in the Store struct and then while accesing it wwith the golang you can simply expire it

		record.expiryAt = time.Now().Add(time.Duration(expiry) * time.Millisecond)
	}
	s.kv[key] = record

	// s.kv[output.Cmd] =

	return respOK, nil
}

func executeGet(output *Command, s *Store) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(output.Args) < 1 {
		return nil, fmt.Errorf("KEy is missing")
	}
	key := output.Args[0]
	val, prsnt := s.kv[key]
	if !prsnt {
		return respNull, nil
	}
	if time.Now().After(val.expiryAt) && !val.expiryAt.IsZero() {

		delete(s.kv, key)

		return respNull, nil

	}
	log.Println(val.value)
	log.Println(val.expiryAt)
	log.Println(val.createdAt)

	return respString(val.value), nil
}
