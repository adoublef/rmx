package main

import (
	"log"
	"os"
	"sort"
	"sync"

	"github.com/choria-io/fisk"
)

const (
	appName = "rmx"
	appHelp = ""
)

func main() {
	cli := fisk.New(appName, appHelp)
	// use lipgloss for ✨vibezzz✨
	cli.UsageWriter(os.Stdout)
	// may not be required
	sort.Slice(commands, func(i int, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	for _, c := range commands {
		c.Command(cli)
	}

	_, err := cli.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

type cmd struct {
	Name    string
	Order   int
	Command func(a app)
}

type app interface {
	Command(name string, help string) *fisk.CmdClause
}

var (
	commands = []*cmd{}
	mu       sync.Mutex
)

func setCmd(name string, order int, c func(a app)) {
	mu.Lock()
	commands = append(commands, &cmd{name, order, c})
	mu.Unlock()
}
