package internal

import (
	"fmt"
	"io"
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
		var buff = make([]byte, 2048)
		n, err := conn.Read(buff)
		if err == io.EOF {
			log.Println("EOF received, closing connection:", conn.RemoteAddr())
			return
		}
		if err != nil {
			log.Println(err)
			return
		}
		command := string(buff[:n])
		log.Println(command)
		conn.Write([]byte("+PONG\r\n"))

	}

}
