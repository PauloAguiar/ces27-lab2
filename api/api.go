package api

type GetArgs struct {
	Key    string
	Quorum int
}

type GetReply struct {
	Value string
}

type PutArgs struct {
	Key    string
	Value  string
	Quorum int
}
