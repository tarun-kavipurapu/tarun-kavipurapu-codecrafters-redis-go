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
	config
}

func NewServer(addr string) *Server {
	if addr == "" {
		addr = DefaultAddr
	}
	return &Server{
		config: config{addr: addr},
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
	defer func() {
		conn.Close()
	}()
	for {
		output, err := s.readCommands(conn)
		if err == io.EOF {
			log.Println("EOF received, closing connection:", conn.RemoteAddr())
			return
		}
		if err != nil {
			log.Println("Error Reading Prefix", err)
		}
		if output.Cmd == "ECHO" {
			ans := output.Args[0]
			conn.Write([]byte("+" + ans + "/r/n"))
		}

	}

}
func (s *Server) readCommands(c io.ReadWriter) (*Command, error) {
	reader := NewRespReader(c)
	command := Command{}
	values, err := reader.CommandRead()
	if err != nil {
		return nil, err
	}
	arrayValue, ok := values.([]interface{})
	if !ok {
		return nil, fmt.Errorf("expected array, got %T", values)

	}
	tokens, err := toArrayString(arrayValue)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty command")

	}
	command.Cmd = strings.ToUpper(tokens[0])

	for _, v := range tokens[1:] {
		command.Args = append(command.Args, strings.ToUpper(v))
	}

	return &command, nil

}
func toArrayString(val []interface{}) ([]string, error) {
	ans := make([]string, len(val))
	for i, v := range val {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a string", i)
		}
		ans[i] = s
	}

	return ans, nil

}
