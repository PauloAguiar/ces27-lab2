package dynamo

import "net/rpc"

// CallInternalHost will communicate to another host through it's InternalRPC
// API.
func (server *Server) CallInternalHost(hostname string, method string, args interface{}, reply interface{}) error {
    var (
        err    error
        client *rpc.Client
    )

    client, err = rpc.Dial(server.connType, hostname)
    if err != nil {
        return err

    }
    defer client.Close()

    err = client.Call("InternalRPC."+method, args, reply)

    if err != nil {
        return err
    }

    return nil
}

// CallHost will communicate to another host through it's RPC public API.
func (server *Server) CallHost(hostname string, method string, args interface{}, reply interface{}) error {
    var (
        err    error
        client *rpc.Client
    )

    client, err = rpc.Dial(server.connType, hostname)
    if err != nil {
        return err
    }

    defer client.Close()

    err = client.Call("RPC."+method, args, reply)

    if err != nil {
        return err
    }

    return nil
}
