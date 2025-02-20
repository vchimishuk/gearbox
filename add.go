package main

import (
	"context"
	"math"
	"os"
	"strings"

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
	return c.Name() + " [-S] [-h host] [-l label] [-p port] file..."
}

func (c *AddCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"l", "", opt.ArgString, "", "labels attache to the torrent to"},
		{"S", "", opt.ArgNone, "", "do not start torrent automatically"},
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

		var labels []string
		if opts.Has("l") {
			labels = strings.Split(opts.StringOr("l", ""), ",")
		}

		paused := opts.Has("S")
		req := &transmission.AddTorrentReq{
			Meta:   f,
			Paused: &paused,
			Labels: labels,
		}

		_, err = client.AddTorrent(context.Background(), req)
		if err != nil {
			return err
		}
	}

	return nil
}
