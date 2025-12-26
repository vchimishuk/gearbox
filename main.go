package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/config"
	"github.com/vchimishuk/opt"
)

const ProgName = "gearbox"
const Usage = "usage: " + ProgName + " [-H] [-h host] [-p port] command [opt]... [arg]..."

const DefaultHost = "localhost"
const DefaultPort = 9091

type Command interface {
	Name() string
	Usage() string
	Options() []*opt.Desc
	Args() (int, int)
	Exec(client *transmission.Client, cfg *config.Config,
		opts opt.Options, args []string) error
}

var Commands []Command = []Command{
	NewAddCommand(),
	NewDeleteCommand(),
	NewEditCommand(),
	NewInfoCommand(),
	NewListCommand(),
	NewStartCommand(),
	NewStatsCommand(),
	NewStopCommand(),
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

func configDir() (string, error) {
	ch := os.Getenv("XDG_CONFIG_HOME")
	if ch != "" {
		return ch, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return filepath.Join(u.HomeDir, ".config/"+ProgName), nil
}

func loadConfig() (*config.Config, error) {
	d, err := configDir()
	if err != nil {
		return nil, err
	}

	cfg, err := config.Parse(filepath.Join(d, ProgName+".conf"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg = &config.Config{}
		} else {
			return nil, err
		}
	}

	return cfg, nil
}

func splitArgs(args []string) ([]string, []string) {
	i := 0
	val := false
	for i < len(args) {
		a := args[i]

		if val {
			// Consume option's value.
			val = false
		} else if a == "-H" {
			// Do nothing.
		} else if strings.HasPrefix(a, "-h") {
			if len(a) == len("-h") {
				val = true
			}
		} else if strings.HasPrefix(a, "-p") {
			if len(a) == len("-p") {
				val = true
			}
		} else {
			break
		}
		i++
	}

	return args[:i], args[i:]
}

func parseMainArgs(args []string) (opt.Options, error) {
	desc := []*opt.Desc{
		{"H", "", opt.ArgNone, "", "display help information and exit"},
		{"h", "", opt.ArgString, "host", "server host name"},
		{"p", "", opt.ArgInt, "port", "server port"},
	}
	opts, args, err := opt.Parse(args, desc)
	if len(args) > 0 {
		panic("invalid state")
	}

	return opts, err
}

func parseCommand(cfg *config.Config, cmd string, args []string) (Command, []string, error) {
	aliases := make(map[string]*config.Alias)
	for _, a := range cfg.Aliases {
		aliases[a.Name] = a
	}

	var pcmd string
	var pargs []string
	alias, ok := aliases[cmd]
	if ok {
		pcmd = alias.Command
		if alias.Args != "" {
			pargs = append(pargs, strings.Split(alias.Args, " ")...)
		}
		pargs = append(pargs, args...)
	} else {
		pcmd = cmd
		pargs = args
	}

	i := slices.IndexFunc(Commands, func(c Command) bool {
		return c.Name() == pcmd
	})
	if i == -1 {
		return nil, nil, fmt.Errorf("invalid command: %s", pcmd)
	}

	return Commands[i], pargs, nil
}

func main() {
	mainArgs, cmdArgs := splitArgs(os.Args[1:])

	mainOpts, err := parseMainArgs(mainArgs)
	if err != nil {
		usage()
		os.Exit(1)
	}
	if mainOpts.Has("H") {
		help()
		os.Exit(1)
	}

	if len(cmdArgs) < 1 {
		usage()
		os.Exit(1)
	}

	cfg, err := loadConfig()
	if err != nil {
		printErr(err)
		os.Exit(1)
	}

	cmd, cargs, err := parseCommand(cfg, cmdArgs[0], cmdArgs[1:])
	if err != nil {
		usage()
		os.Exit(1)
	}

	minArgs, maxArgs := cmd.Args()

	opts, args, err := opt.Parse(cargs, cmd.Options())
	if err != nil || len(args) < minArgs || len(args) > maxArgs {
		fmt.Fprintf(os.Stderr, "usage: %s %s\n", ProgName, cmd.Usage())
		os.Exit(1)
	}

	host := DefaultHost
	port := DefaultPort
	if h, ok := mainOpts.String("h"); ok {
		host = h
	} else if cfg.Host != "" {
		host = cfg.Host
	}
	if p, ok := mainOpts.Int("p"); ok {
		port = p
	} else if cfg.Port != 0 {
		port = cfg.Port
	}

	uri := fmt.Sprintf("http://%s:%d", host, port)
	client, err := transmission.New(uri)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}

	err = cmd.Exec(client, cfg, opts, args)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
}
