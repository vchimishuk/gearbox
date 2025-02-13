package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/opt"
)

const ProgName = "gearbox"

type Command interface {
	Name() string
	Usage() string
	Options() []*opt.Desc
	Args() (int, int)
	Exec(client *transmission.Client, opts opt.Options, args []string) error
}

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", ProgName, err)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s command [opt]... [arg]...\n",
		ProgName)
}

func main() {
	commands := []Command{
		NewInfoCommand(),
		NewListCommand(),
	}

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	i := slices.IndexFunc(commands, func(c Command) bool {
		return c.Name() == os.Args[1]
	})
	if i == -1 {
		usage()
		os.Exit(1)
	}

	cmd := commands[i]
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
