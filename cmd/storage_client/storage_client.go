package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ashriths/go-graph/common"
	"github.com/ashriths/go-graph/graph"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/system"
	"github.com/google/uuid"
	"os"
	"strings"
)

func kv(k, v string) *storage.KeyValue {
	return &storage.KeyValue{k, v}
}

func pat(pre, suf string) *storage.Pattern {
	return &storage.Pattern{pre, suf}
}

func kva(args []string) *storage.KeyValue {
	if len(args) == 1 {
		return kv("", "")
	} else if len(args) == 2 {
		return kv(args[1], "")
	}
	return kv(args[1], args[2])
}

func pata(args []string) *storage.Pattern {
	if len(args) == 1 {
		return pat("", "")
	} else if len(args) == 2 {
		return pat(args[1], "")
	}
	return pat(args[1], args[2])
}

func single(args []string) string {
	if len(args) == 1 {
		return ""
	}
	return args[1]
}

func printList(lst []string) {
	for _, e := range lst {
		fmt.Println(e)
	}
}

var guuid uuid.UUID

const help = `Usage:
   storage-client <server address> [command <args...>]

With no command specified to enter interactive mode. 
` + cmdHelp

const cmdHelp = `Command list:
   add-v <data>
   get-v <uuid>
   del-v <uuid>
   add-e <src-uuid> <dest-uuid> data
   help
   exit
`

func runCmd(s *storage.StorageClient, args []string) bool {
	var u uuid.UUID
	var u1 uuid.UUID
	var u2 uuid.UUID
	var succ bool
	var data graph.ElementProperty
	var e error
	var v graph.Vertex
	var ed graph.Edge
	cmd := args[0]

	switch cmd {
	case "add-v":
		if e := json.Unmarshal([]byte(args[1]), &data); e != nil {
			common.LogError(e)
			return false
		}
		u, e = uuid.NewUUID()
		if e := s.StoreVertex(graph.V(guuid, u, data), &succ); e != nil {
			common.LogError(e)
			return false
		}
		fmt.Println(succ)
	case "get-v":
		u, e = uuid.Parse(args[1])
		if e != nil {
			common.LogError(e)
			return false
		}
		if e := s.GetVertexById(u, &v); e != nil {
			common.LogError(e)
			return false
		}
		fmt.Println(v.Json())
	case "del-v":
		u, e = uuid.Parse(args[1])
		if e != nil {
			common.LogError(e)
			return false
		}
		if e := s.RemoveVertex(u, &succ); e != nil {
			common.LogError(e)
			return false
		}
		fmt.Println(succ)
	case "add-e":
		u, e = uuid.NewUUID()
		u1, e = uuid.Parse(args[1])
		if e != nil {
			common.LogError(e)
			return false
		}
		u2, e = uuid.Parse(args[2])
		if e != nil {
			common.LogError(e)
			return false
		}
		if e := s.StoreEdge(graph.E(guuid, u, u1, u2, data), &succ); e != nil {
			common.LogError(e)
			return false
		}
		fmt.Println(succ)
	case "get-e":
		u, e = uuid.Parse(args[1])
		if e != nil {
			common.LogError(e)
			return false
		}
		if e := s.GetEdgeById(u, &ed); e != nil {
			common.LogError(e)
			return false
		}
		fmt.Println(ed.Json())
	case "del-e":
		u, e = uuid.Parse(args[1])
		if e != nil {
			common.LogError(e)
			return false
		}
		if e := s.RemoveEdge(u, &succ); e != nil {
			common.LogError(e)
			return false
		}
		fmt.Println(succ)
	case "help":
		fmt.Println(cmdHelp)
	case "exit":
		return true
	default:
		common.LogError(fmt.Errorf("bad command, try \"help\"."))
	}
	return false
}

func fields(s string) []string {
	return strings.Fields(s)
}

func runPrompt(s *storage.StorageClient) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")

	for scanner.Scan() {
		line := scanner.Text()
		args := fields(line)
		if len(args) > 0 {
			if runCmd(s, args) {
				break
			}
		}
		fmt.Print("> ")
	}

	e := scanner.Err()
	if e != nil {
		panic(e)
	}
}

func main() {
	system.Logging = true
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}

	addr := args[0]
	s := storage.NewStorageClient(addr)

	var e error
	guuid,e   = uuid.NewUUID()
	common.NoError(e)

	cmdArgs := args[1:]
	if len(cmdArgs) == 0 {
		runPrompt(s)
		fmt.Println()
	} else {
		runCmd(s, cmdArgs)
	}
}
