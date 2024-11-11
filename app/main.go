package main

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

func main() {
	server := internal.NewServer(internal.DefaultAddr)
	err := server.ListenAndAccept()
	if err != nil {
		log.Println(err)
	}
}
