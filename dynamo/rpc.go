package dynamo

import "github.com/pauloaguiar/ces27-lab2/api"

// RPC is the struct on which the public API operations will be called.
type RPC struct {
	server *Server
}

// NewRPC returns a new RPC struct instance that contains the procedures specified
// in the public API.
func NewRPC(server *Server) *RPC {
	return &RPC{server}
}

// Get handles the read operations on the database. It'll redirect all operations
// to the router that will find a suitable coordinator for the operation and wait
// for it to coordinate the operation and return a reply.
func (rpc *RPC) Get(args *api.GetArgs, reply *api.GetReply) error {
	var (
		err   error
		value string
	)

	value, err = rpc.server.RouteGet(args.Key, args.Quorum)

	if err != nil {
		return err
	}

	*reply = api.GetReply{value}
	return nil
}

// Put handles the write operations on the database. It'll redirect all operations
// to the router that will find a suitable coordinator for the operation and wait
// for it to coordinate the operation and return.
func (rpc *RPC) Put(args *api.PutArgs, _ *struct{}) error {
	rpc.server.RoutePut(args.Key, args.Value, args.Quorum)
	return nil
}
