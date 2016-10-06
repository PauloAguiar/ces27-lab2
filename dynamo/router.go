package dynamo

import "log"

// RoutePut will handle the search for a coordinator of a Put operation. It'll
// start a replication if it finds itself as coordinator or delegate it to the
// remote coordinator and wait for it.
func (server *Server) RoutePut(key string, value string, quorum int) error {
	var (
		err                 error
		coordinatorId       string
		coordinatorHostname string
	)

	log.Printf("[ROUTER] Routing Put of KV('%v', '%v') with quorum Q('%v').\n", key, value, quorum)

	coordinatorId, coordinatorHostname = server.ring.GetCoordinator(key)

	for server.id != coordinatorId {
		log.Printf("[ROUTER] Trying '%v' as coordinator.\n", coordinatorId)
		err = server.CallInternalHost(coordinatorHostname, "CoordinatePut", &CoordinatePutArgs{key, value, quorum}, nil)

		if err == nil {
			log.Printf("[ROUTER] Coordinate succeded.\n")
			break
		}

		if err == notEnoughQuorumErr {
			log.Printf("[ROUTER] Not enough quorum error.\n")
			return err
		}

		log.Printf("[ROUTER] Coordinator tryout failed. Error: %v.\n", err)

		coordinatorId, coordinatorHostname, err = server.ring.GetNextCoordinator(coordinatorId)

		if err != nil {
			log.Printf("[ROUTER] Failed to find next coordinator to '%v'.\n", coordinatorId)
			return err
		}
	}

	if server.id == coordinatorId {
		err = server.Replicate(key, value, quorum)
	}

	return err
}

// RouteGet will handle the search for a coordinator of a Get operation. It'll
// start a votation if it finds itself as coordinator or delegate the it to the
// remote coordinator and wait for it.
func (server *Server) RouteGet(key string, quorum int) (value string, err error) {
	var (
		coordinatorId       string
		coordinatorHostname string
		reply               CoordinateGetReply
	)

	log.Printf("[ROUTER] Routing Get of K('%v') with quorum Q('%v').\n", key, quorum)

	coordinatorId, coordinatorHostname = server.ring.GetCoordinator(key)

	for server.id != coordinatorId {
		log.Printf("[ROUTER] Trying '%v' as coordinator.\n", coordinatorId)
		err = server.CallInternalHost(coordinatorHostname, "CoordinateGet", &CoordinateGetArgs{key, quorum}, &reply)

		if err == nil {
			log.Printf("[ROUTER] Coordinate succeded.\n")
			value = reply.Value
			break
		}

		if err == notEnoughQuorumErr {
			log.Printf("[ROUTER] Not enough quorum error.\n")
			return "", err
		}

		log.Printf("[ROUTER] Coordinator tryout failed. Error: %v.\n", err)

		coordinatorId, coordinatorHostname, err = server.ring.GetNextCoordinator(coordinatorId)

		if err != nil {
			log.Printf("[ROUTER] Failed to find next coordinator to '%v'.\n", coordinatorId)
			return "", err
		}
	}

	if server.id == coordinatorId {
		value, err = server.Voting(key, quorum)

		if err != nil {
			return "", err
		}
	}

	return value, nil
}
