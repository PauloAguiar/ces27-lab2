package dynamo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type Console struct {
	cache  *Cache
	server *Server
}

func NewConsole(cache *Cache, server *Server) *Console {
	return &Console{cache, server}
}

func (console *Console) Run() {
	var (
		scanner *bufio.Scanner
		input   string
		tokens  []string
	)

	scanner = bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input = scanner.Text()
		tokens = strings.Split(input, " ")

		switch tokens[0] {
		case "get":
			if len(tokens) < 2 {
				fmt.Println("[CONSOLE] usage: get <key>")
				break
			}

			key := tokens[1]
			value, timestamp := console.cache.Get(key)

			fmt.Printf("[CONSOLE] Get result: '%v' = '%v (ts: '%v')\n", key, value, timestamp)

		case "rget":
			var (
				value string
				err   error
			)
			if len(tokens) < 3 {
				fmt.Println("[CONSOLE] usage: rget <key> <quorum>")
				break
			}

			key := tokens[1]
			quorum, err := strconv.Atoi(tokens[2])

			if err != nil {
				fmt.Println("[CONSOLE] invalid quorum value.")
				break
			}

			value, err = console.server.RouteGet(key, quorum)

			if err != nil {
				fmt.Println("[CONSOLE] Rget result: failed.")
				break
			}

			fmt.Printf("[CONSOLE] Rget result: '%v' = '%v'\n", key, value)
		case "put":
			if len(tokens) < 3 {
				fmt.Println("[CONSOLE] usage: put <key> <value>")
				break
			}

			key := tokens[1]
			value := tokens[2]
			timestamp := time.Now().Unix()
			console.cache.Put(key, value, timestamp)

		case "rput":
			if len(tokens) < 4 {
				fmt.Println("[CONSOLE] usage: rput <key> <value> <quorum>")
				break
			}

			key := tokens[1]
			value := tokens[2]
			quorum, err := strconv.Atoi(tokens[3])

			if err != nil {
				fmt.Println("[CONSOLE] invalid quorum value.")
				break
			}

			console.server.RoutePut(key, value, quorum)

		case "print":
			w := tabwriter.NewWriter(os.Stdout, 5, 0, 1, '_', tabwriter.Debug)

			cacheMap, cacheTimestamps := console.cache.getAll()
			fmt.Fprintf(w, "[CONSOLE] KEY\tVALUE\tTIMESTAMP\t\n")
			for key, value := range cacheMap {
				fmt.Fprintf(w, "[CONSOLE] '%v'\t'%v'\t%v\t\n", key, value, cacheTimestamps[key])
			}
			w.Flush()

		case "ring":
			w := tabwriter.NewWriter(os.Stdout, 5, 0, 1, '_', tabwriter.Debug)
			fmt.Fprintf(w, "[CONSOLE] HASH\tID\t\n")
			for _, node := range console.server.ring.hashring.Nodes {
				fmt.Fprintf(w, "[CONSOLE] '%v'\t'%v'\t\n", node.HashId, node.Id)
			}
			w.Flush()

		case "down":
			fmt.Println("[CONSOLE] Putting server DOWN.")
			go console.server.Stop()

		case "up":
			fmt.Println("[CONSOLE] Putting server UP.")
			go console.server.Start()
		}
	}
}
