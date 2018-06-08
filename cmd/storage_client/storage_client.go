package storage_client

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"bufio"
	"github.com/ashriths/go-graph/storage"
	"github.com/ashriths/go-graph/graph"
	"github.com/google/uuid"
)

func noError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func logError(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
	}
}

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

func V(data string) *graph.GoGraphVertex{
	return &graph.GoGraphVertex{GoGraphElement:graph.GoGraphElement{Label:data}}
}

const help = `Usage:
   storage-client <server address> [command <args...>]

With no command specified to enter interactive mode. 
` + cmdHelp

const cmdHelp = `Command list:
   get <key>
   set <key> <value>
   keys [<prefix> [<suffix>]]
   list-get <key>
   list-append <key> <value>
   list-remove <key> <value>
   list-keys [<prefix> [<suffix]]
   clock [<atleast=0>]
   help
   exit
`

func runCmd(s *storage.StorageClient, args []string) bool {
	var u uuid.UUID

	cmd := args[0]

	switch cmd {
	case "add-v":
		logError( s.StoreVertex(graph.V(args[1]), &u) )
		fmt.Println(u)
	case "help":
		fmt.Println(cmdHelp)
	case "exit":
		return true
	default:
		logError(fmt.Errorf("bad command, try \"help\"."))
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
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}

	addr := args[0]
	s := storage.NewStorageClient(addr)

	cmdArgs := args[1:]
	if len(cmdArgs) == 0 {
		runPrompt(s)
		fmt.Println()
	} else {
		runCmd(s, cmdArgs)
	}
}