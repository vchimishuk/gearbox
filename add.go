package main

import (
	"context"
	"math"
	"os"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/opt"
)

type AddCommand struct {
}

func NewAddCommand() *AddCommand {
	return &AddCommand{}
}

func (c *AddCommand) Name() string {
	return "add"
}

func (c *AddCommand) Usage() string {
	return c.Name() + " [-h host] [-p port] file..."
}

func (c *AddCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"P", "", opt.ArgNone, "", "do not automatically start torrent"},
	}
}

func (c *AddCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *AddCommand) Exec(client *transmission.Client, opts opt.Options, args []string) error {
	for _, n := range args {
		f, err := os.Open(n)
		if err != nil {
			return err
		}
		defer f.Close()

		paused := opts.Has("P")
		req := &transmission.AddTorrentReq{
			Meta:   f,
			Paused: &paused,
		}

		_, err = client.AddTorrent(context.Background(), req)
		if err != nil {
			return err
		}
	}

	return nil
}
