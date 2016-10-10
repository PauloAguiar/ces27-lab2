package dynamo

import (
    "log"

    "github.com/pauloaguiar/ces27-lab2/common/consistenthash"
)

// Ring is the struct that will hold the current status of the consistent hash
// data structure that will be used to replicate data among nodes.
type Ring struct {
    hashring     *consistenthash.Ring
    idToHostname map[string]string
    server       *Server
}

// NewRing creates a new Ring object and return a pointer to it.
func NewRing(server *Server) *Ring {
    var newRing *Ring

    newRing = new(Ring)
    newRing.hashring = consistenthash.NewRing()
    newRing.idToHostname = make(map[string]string)
    newRing.server = server
    return newRing
}

// Add a node to the ring structure. This also handle collisions and update it
// in case of changes. The idToHostname is always update to represent the last
// state of a known node.
func (ring *Ring) AddNode(id, hostname string) {
    var (
        node   *consistenthash.Node
        exists bool
    )

    if exists, node = ring.hashring.Exists(id); !exists {
        node = ring.hashring.AddNode(id)
        ring.idToHostname[id] = hostname
        log.Printf("[RING] Added node '%v'(hostname: '%v') to the local ring(hash: '%v')\n", id, hostname, node.HashId)
    } else {
        log.Printf("[RING] Node '%v'(hostname: '%v') already exists in the local ring(hash: '%v')\n", id, hostname, node.HashId)

        if ring.idToHostname[id] != hostname {
            log.Printf("[RING] Updating '%v' hostname from '%v' to '%v'\n", id, ring.idToHostname[id], hostname)
            ring.idToHostname[id] = hostname
        }
    }
}

// GetCoordinator will return the main coordinator for a key. The main
// coordinator is the first node that has a hash higher than the hash of the key.
func (ring *Ring) GetCoordinator(key string) (id, hostname string) {
    id = ring.hashring.Get(key)
    hostname = ring.idToHostname[id]
    return
}

// GetNextCoordinator return the next node in the ring starting from the node
// represented by id.
func (ring *Ring) GetNextCoordinator(id string) (nextId, hostname string, err error) {
    nextId, err = ring.hashring.GetNext(id)

    if err != nil {
        return
    }

    hostname = ring.idToHostname[nextId]
    return
}

// GetNode return id and hostname of the node responsible for a key.
func (ring *Ring) GetNode(key string) (id, hostname string) {
    id = ring.hashring.Get(key)
    hostname = ring.idToHostname[id]
    return
}

// GetNodes return count nodes starting from the one represented by id.
func (ring *Ring) GetNodes(id string, count int) ([]string, error) {
    var (
        err       error
        nodes     []string
        currentId string
    )

    nodes = make([]string, 0)
    nodes = append(nodes, ring.idToHostname[id])

    currentId = id

    for i := 0; i < count-1; i++ {
        currentId, err = ring.hashring.GetNext(currentId)

        if err != nil {
            return nil, err
        }

        nodes = append(nodes, ring.idToHostname[currentId])
    }

    return nodes, nil
}

// Synchronizes this instance ring struct with the one in joinAddress node.
func (ring *Ring) Sync(joinAddress string) {
    var (
        err     error
        reply   SyncRingsReply
        ringMap *map[string]string
    )

    err = ring.server.CallInternalHost(joinAddress, "SyncRings", new(struct{}), &reply)

    if err != nil {
        log.Printf("[RING] Failed to Sync ring with host '%v'. Error: %v\n", joinAddress, err)
        return
    }

    for id, hostname := range reply.RingMap {
        ring.AddNode(id, hostname)
    }

    // The node always add itself after getting the seed. This will ensure that
    // the ring is updated with new data before reporting the data.
    ring.AddNode(ring.server.id, ring.server.connHostname)

    ringMap = ring.GetMap()
    for id, hostname := range *ringMap {
        // A server should not report to itself
        if id != ring.server.id {
            go ring.Report(hostname)
        }
    }
}

// Report calls a hostname node and report that this node has been added to the
// ring.
func (ring *Ring) Report(hostname string) {
    var (
        err error
    )

    log.Printf("[RING] Reporting to host '%v'\n", hostname)
    err = ring.server.CallInternalHost(hostname, "AddNode", &AddNodeArgs{ring.server.id, ring.server.connHostname}, new(struct{}))

    if err != nil {
        log.Printf("[RING] Failed to Report to host '%v'. Error: %v\n", hostname, err)
        return
    }
}

// Return a map of nodes in which the key is the id of the node and the value is
// its hostname.
func (ring *Ring) GetMap() *map[string]string {
    var (
        ringMap map[string]string
    )

    ringMap = make(map[string]string)

    ring.hashring.Lock()

    for _, hashnode := range ring.hashring.Nodes {
        ringMap[hashnode.Id] = ring.idToHostname[hashnode.Id]
    }

    ring.hashring.Unlock()

    return &ringMap
}
