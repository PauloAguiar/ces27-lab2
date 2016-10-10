package main

import (
  "fmt"
  "log"
  "net/rpc"

  "github.com/pauloaguiar/ces27-lab2/api"
)

func main() {
  client, err := rpc.Dial("tcp", "localhost:3000")

  if err != nil {
    log.Fatal("dialing:", err)
  }

  e := client.Call("RPCServer.Put", &api.PutArgs{"SomeKey", "Test"}, nil)

  if e != nil {
    log.Fatalf("Something went wrong: %v", e.Error())
  }

  var reply api.GetReply

  e = client.Call("RPCServer.Get", &api.GetArgs{"SomeKey"}, &reply)

  fmt.Printf("The 'reply' pointer value has been changed to: %s", reply)
  client.Close()
}
