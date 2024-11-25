package internal

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var DefaultAddr = "0.0.0.0:6379"

type config struct {
	addr string
}

type Server struct {
	listener net.Listener
	store    *Store
	config
}

func NewServer(addr string, store *Store) *Server {
	if addr == "" {
		addr = DefaultAddr
	}
	return &Server{
		config: config{addr: addr},
		store:  store,
	}
}
func (s *Server) ListenAndAccept() error {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	fmt.Printf("Your dummy redis Server Started Listening on %s\n", DefaultAddr)

	return s.acceptLoop()
}
func (s *Server) acceptLoop() error {
	for {

		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnLoop(conn)
	}
}
func (s *Server) handleConnLoop(conn net.Conn) {
	defer conn.Close()

	for {
		output, err := s.readCommands(conn)
		if err == io.EOF {
			log.Printf("Client disconnected: %s\n", conn.RemoteAddr())
			return
		}

		if err != nil {
			log.Printf("Error reading command: %v\n", err)
			conn.Write([]byte(fmt.Sprintf("-ERR %v\r\n", err)))
			continue
		}

		if output == nil {
			log.Println("Invalid RESP command format")
			conn.Write([]byte("-ERR invalid command format\r\n"))
			continue
		}

		outputString, err := EvaluateFunc(output, s.store)
		if err != nil {
			log.Printf("Error evaluating command: %v\n", err)
			conn.Write([]byte(fmt.Sprintf("-ERR %v\r\n", err)))
			continue
		}

		_, writeErr := conn.Write(outputString)
		if writeErr != nil {
			log.Printf("Error writing response: %v\n", writeErr)
			return
		}
	}
}
func (s *Server) readCommands(c io.ReadWriter) (*Command, error) {
	reader := NewRespReader(c)
	values := make([]interface{}, 0)

	// Read multiple RESP messages until the buffer is empty
	for {
		value, err := reader.CommandRead()
		if err == io.EOF {
			return nil, io.EOF
		}
		if err != nil {
			if len(values) == 0 {
				return nil, fmt.Errorf("error reading command: %v", err)
			}
			break
		}
		values = append(values, value)

		// Break the loop if there is no more buffered data
		if reader.c.Buffered() == 0 {
			break
		}
	}

	log.Println("Parsed RESP values:", values)

	if len(values) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	tokens, err := toFlatStringArray(values)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty command array")
	}

	return &Command{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

// Recursive function to flatten `[]interface{}` into a `[]string`
func toFlatStringArray(values []interface{}) ([]string, error) {
	var result []string

	for _, value := range values {
		switch v := value.(type) {
		case string:
			result = append(result, v)
		case []interface{}:
			flattened, err := toFlatStringArray(v)
			if err != nil {
				return nil, err
			}
			result = append(result, flattened...)
		default:
			return nil, fmt.Errorf("unsupported type %T in values", v)
		}
	}

	return result, nil
}
