package main

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//create a global store
	store := internal.NewStore()
	//create the central Server needed
	server := internal.NewServer(internal.DefaultAddr, store)
	err := server.ListenAndAccept()
	if err != nil {
		log.Println(err)
	}
}
