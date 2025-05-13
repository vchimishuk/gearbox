package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/opt"
)

const ProgName = "gearbox"
const Usage = "usage: " + ProgName + " [-H] command [opt]... [arg]..."

type Command interface {
	Name() string
	Usage() string
	Options() []*opt.Desc
	Args() (int, int)
	Exec(client *transmission.Client, opts opt.Options, args []string) error
}

var Commands []Command = []Command{
	NewAddCommand(),
	NewInfoCommand(),
	NewListCommand(),
	NewStatsCommand(),
	NewTorrentCommand(),
}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", ProgName, err)
}

func usage() {
	fmt.Fprintf(os.Stderr, "%s\n", Usage)
}

func help() {
	fmt.Printf("gearbox is a non-interactive client for transmission-daemon\n")
	fmt.Printf("\n")
	fmt.Printf("%s\n", Usage)
	fmt.Printf("\n")
	fmt.Printf("commands:\n")
	for _, c := range Commands {
		fmt.Printf("  %s\n", c.Usage())
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	if os.Args[1] == "-H" {
		help()
		os.Exit(1)
	}

	i := slices.IndexFunc(Commands, func(c Command) bool {
		return c.Name() == os.Args[1]
	})
	if i == -1 {
		usage()
		os.Exit(1)
	}

	cmd := Commands[i]
	cmdOpts := cmd.Options()
	cmdOpts = append(cmdOpts, &opt.Desc{
		"h", "", opt.ArgString, "host", "server host name",
	})
	cmdOpts = append(cmdOpts, &opt.Desc{
		"p", "", opt.ArgInt, "port", "server port",
	})
	minArgs, maxArgs := cmd.Args()

	opts, args, err := opt.Parse(os.Args[2:], cmdOpts)
	if err != nil || len(args) < minArgs || len(args) > maxArgs {
		fmt.Fprintf(os.Stderr, "usage: %s %s\n", ProgName, cmd.Usage())
		os.Exit(1)
	}

	host := "j4105.localdomain"
	port := 9091

	uri := fmt.Sprintf("http://%s:%d", host, port)
	client, err := transmission.New(uri)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}

	err = cmd.Exec(client, opts, args)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
}
