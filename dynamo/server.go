package dynamo

import (
    "log"
    "net"
    "net/rpc"
)

// Server is the struct that hold all information that is used by this instance
// of dynamo server.
type Server struct {
    // Network
    listener     net.Listener
    connType     string
    connHostname string

    // Partitioning
    id       string
    ring     *Ring
    replicas int

    // Context
    err  error
    done chan struct{}

    // Storage
    cache *Cache
}

// Create a new server object and return a pointer to it.
func NewServer(id string, hostname string, cache *Cache) *Server {
    var server Server

    server.done = make(chan struct{})

    server.cache = cache
    server.replicas = 3

    server.connType = "tcp"
    server.connHostname = hostname

    server.id = id

    server.ring = NewRing(&server)

    return &server
}

// Run a dynamo server. The parameter joinAddress is the address to the server
// which this instance should use get the current status of the consistenthash
// ring.
func (server *Server) Run(joinAddress string) {
    log.Printf("[SERVER] Running Dynamo Server with id '%v'\n", server.id)

    if joinAddress != "" {
        server.ring.Sync(joinAddress)
    } else {
        server.ring.AddNode(server.id, server.connHostname)
    }

    // Start RPC servers.
    // RPC is the public API and InternalRPC is the internal server-to-server
    // API.
    rpcs := NewRPC(server)
    internalrpcs := NewInternalRPC(server)
    rpc.Register(rpcs)
    rpc.Register(internalrpcs)

    // Start the network interface
    go server.Start()
}

// Start will create a new listener on the given interface and handle
// all the connections afterwards until the server is stopped.
func (server *Server) Start() {
    var (
        err error
    )

    server.listener, err = net.Listen(server.connType, server.connHostname)

    if err != nil {
        server.err = err
        close(server.done)
        return
    }

    log.Printf("[SERVER] Listening on '%v'\n", server.connHostname)

    for {
        conn, err := server.listener.Accept()

        if err != nil {
            return
        }

        go rpc.ServeConn(conn)
    }
}

// Stop will close the listener create in the Start method, causing it to return.
func (server *Server) Stop() {
    log.Printf("[SERVER] Server Stopped \n")
    server.listener.Close()
}

// Done returns a channel that is closed when server is done
func (server *Server) Done() <-chan struct{} {
    return server.done
}

// Err indicates why this context was canceled
func (server *Server) Err() error {
    return server.err
}
