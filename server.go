package main

import (
	"flag"

	"github.com/pauloaguiar/ces27-lab2/dynamo"
)

var (
	localAddress = flag.String("addr", "localhost", "IP Address to listen on for client connections")
	localPort    = flag.String("port", "3000", "TCP port to listen on for client connections")
	ringAddress  = flag.String("ring", "", "Address of ring coordinator")
	instanceId   = flag.String("id", "", "ID of the instance")
)

// Entry point for our Dynamo server application
// This will start a instance of the server, as well as a cache(the in-memory store)
// and a console that will create a CLI to perform operation on the server.
func main() {
	var (
		hostname string
		id       string
		cache    *dynamo.Cache
		server   *dynamo.Server
		console  *dynamo.Console
	)

	flag.Parse()

	cache = dynamo.NewCache()

	hostname = *localAddress + ":" + *localPort

	// If an id isn't provided, we use the hostname instead
	if *instanceId != "" {
		id = *instanceId
	} else {
		id = hostname
	}

	server = dynamo.NewServer(id, hostname, cache)

	console = dynamo.NewConsole(cache, server)

	// Spawn goroutines to handle both interfaces
	go server.Run(*ringAddress)
	go console.Run()

	// Wait fo the server to finish
	<-server.Done()
}
