package dynamo

import "log"

type SyncRingsReply struct {
	RingMap map[string]string
}

type AddNodeArgs struct {
	Id       string
	Hostname string
}

type InternalRPC struct {
	server *Server
}

type CoordinatePutArgs struct {
	Key    string
	Value  string
	Quorum int
}

type CoordinateGetArgs struct {
	Key    string
	Quorum int
}

type CoordinateGetReply struct {
	Value string
}

type ReplicateArgs struct {
	Key       string
	Value     string
	Timestamp int64
}

type VoteArgs struct {
	Key string
}

type VoteReply struct {
	Value     string
	Timestamp int64
}

func NewInternalRPC(server *Server) *InternalRPC {
	return &InternalRPC{server}
}

// SyncRings will reply with the current state of the ring in this node.
func (rpc *InternalRPC) SyncRings(_ *struct{}, reply *SyncRingsReply) error {
	reply.RingMap = *rpc.server.ring.GetMap()
	return nil
}

// AddNode will add a new node to the ring in this node.
func (rpc *InternalRPC) AddNode(args *AddNodeArgs, _ *struct{}) error {
	go rpc.server.ring.AddNode(args.Id, args.Hostname)
	return nil
}

// CoordinatePut will start a replication operation.
func (rpc *InternalRPC) CoordinatePut(args *CoordinatePutArgs, _ *struct{}) error {
	return rpc.server.Replicate(args.Key, args.Value, args.Quorum)
}

// CoordinateGet will start a voting operation and return the value voted.
func (rpc *InternalRPC) CoordinateGet(args *CoordinateGetArgs, reply *CoordinateGetReply) error {
	var (
		err   error
		value string
	)
	value, err = rpc.server.Voting(args.Key, args.Quorum)

	if err != nil {
		return err
	}

	*reply = CoordinateGetReply{value}
	return nil
}

// Replicate will be called by a replication operation to store data.
func (rpc *InternalRPC) Replicate(args *ReplicateArgs, _ *struct{}) error {
	log.Printf("[INTERNAL] Replicating '%v' = '%v' (timestamp: '%v')\n", args.Key, args.Value, args.Timestamp)
	rpc.server.cache.Put(args.Key, args.Value, args.Timestamp)
	return nil
}

// Vote will be called by a voting operation to read the stored data.
func (rpc *InternalRPC) Vote(args *VoteArgs, reply *VoteReply) error {
	var (
		value     string
		timestamp int64
	)

	value, timestamp = rpc.server.cache.Get(args.Key)

	log.Printf("[INTERNAL] Voting for '%v' = '%v' (timestamp: %v)\n", args.Key, value, timestamp)

	*reply = VoteReply{value, timestamp}
	return nil
}
