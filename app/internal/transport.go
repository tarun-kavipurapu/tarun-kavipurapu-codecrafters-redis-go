package internal

import (
	"log"
	"net"
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

	return s.acceptLoop()
}
func (s *Server) acceptLoop() error {
	for {

		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		log.Print(conn)
	}
}
