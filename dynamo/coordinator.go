package dynamo

import (
	"errors"
	"log"
	"time"
)

var notEnoughQuorumErr = errors.New("not enough quorum")

// Replicate will coordinate the replication("sharding") of the data into the
// preferred nodes.
func (server *Server) Replicate(key string, value string, quorum int) error {
	var (
		err        error
		nodes      []string
		reportChan chan error
		successes  int
		failures   int
		baseNode   string
		timestamp  int64
	)

	log.Printf("[COORDINATOR] Coordinating replication of KV('%v', '%v') with Quorum '%v'\n", key, value, quorum)

	timestamp = time.Now().Unix()
	log.Printf("[COORDINATOR] Operation timestamp: '%v'\n", timestamp)

	baseNode, _ = server.ring.GetNode(key)

	nodes, err = server.ring.GetNodes(baseNode, server.replicas)

	if err != nil {
		log.Printf("[COORDINATOR] Failed to replicate nodes. Error: %v\n", err)
		return err
	}

	reportChan = make(chan error, server.replicas)
	for _, node := range nodes {
		if server.connHostname == node {
			go func() {
				server.cache.Put(key, value, timestamp)

				reportChan <- nil
			}()
		} else {
			go func(hostname string) {
				var err error

				err = server.CallInternalHost(hostname, "Replicate", &ReplicateArgs{key, value, timestamp}, nil)

				reportChan <- err
			}(node)
		}
	}

	successes = 0
	failures = 0
	for report := range reportChan {
		if report != nil {
			log.Printf("[COORDINATOR] Error on replication: %v\n", report)
			failures++
		} else {
			successes++
		}

		if successes == quorum {
			log.Printf("[COORDINATOR] Replication with quorum '%v' succeded.\n", quorum)
			break
		}

		if failures+successes == len(nodes) {
			log.Printf("[COORDINATOR] Replication failed. Not enough quorum\n")
			return notEnoughQuorumErr
		}
	}

	go func() {
		for report := range reportChan {
			if report != nil {
				log.Printf("[COORDINATOR] Error on replication: %v\n", report)
				continue
			}
		}
	}()

	return nil
}

// vote is a value object to be passed around in channel inside the Voting
// operation.
type vote struct {
	value     string
	timestamp int64
	err       error
}

// Replicate will coordinate the voting(to decide on a value to be returned) of
// the data replicated in the preferred nodes.
func (server *Server) Voting(key string, quorum int) (string, error) {
	var (
		err        error
		nodes      []string
		reportChan chan *vote
		successes  int
		failures   int
		votes      []*vote
		baseNode   string
	)

	log.Printf("[COORDINATOR] Coordinating voting of K('%v') with Quorum '%v'\n", key, quorum)

	baseNode, _ = server.ring.GetNode(key)

	nodes, err = server.ring.GetNodes(baseNode, server.replicas)

	if err != nil {
		log.Printf("[COORDINATOR] Failed to gather votes from nodes.\n")
		return "", err
	}

	reportChan = make(chan *vote, server.replicas)

	for _, node := range nodes {
		if server.connHostname == node {
			go func() {
				var (
					v         *vote
					value     string
					timestamp int64
				)
				value, timestamp = server.cache.Get(key)
				v = &vote{value, timestamp, nil}
				reportChan <- v
			}()
		} else {
			go func(hostname string) {
				var (
					err   error
					reply VoteReply
					v     *vote
				)

				err = server.CallInternalHost(hostname, "Vote", &VoteArgs{key}, &reply)

				v = &vote{reply.Value, reply.Timestamp, err}
				reportChan <- v
			}(node)
		}
	}

	successes = 0
	failures = 0
	votes = make([]*vote, 0)
	for report := range reportChan {
		if report.err != nil {
			log.Printf("[COORDINATOR] Error on vote: %v\n", report.err)
			failures++
		} else {
			successes++
			votes = append(votes, report)
		}

		if successes == quorum {
			log.Printf("[COORDINATOR] Voting with quorum '%v' succeded.\n", quorum)
			break
		}

		if failures+successes == len(nodes) {
			log.Printf("[COORDINATOR] Voting failed. Not enough quorum\n")
			return "", notEnoughQuorumErr
		}
	}

	go func() {
		for report := range reportChan {
			if report.err != nil {
				log.Printf("[COORDINATOR] Error on vote: %v\n", report.err)
				continue
			}
		}
	}()

	return aggregateVotes(votes), nil
}

// aggregateVotes will select the right value from the votes received.
func aggregateVotes(votes []*vote) (result string) {
	for _, vote := range votes {
		log.Printf("[COORDINATOR] Vote: %v\n", vote.value)
	}

	/////////////////////////
	// YOUR CODE GOES HERE //
	/////////////////////////
	result = votes[0].value
	return
}
