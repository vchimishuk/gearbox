package main

import (
	"context"
	"math"
	"strings"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/config"
	"github.com/vchimishuk/opt"
)

type EditCommand struct {
}

func NewEditCommand() *EditCommand {
	return &EditCommand{}
}

func (c *EditCommand) Name() string {
	return "edit"
}

func (c *EditCommand) Usage() string {
	return c.Name() + " [-l label[,label...]] id..."
}

func (c *EditCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"l", "", opt.ArgString, "", "edit labels"},
	}
}

func (c *EditCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *EditCommand) Exec(client *transmission.Client, cfg *config.Config,
	opts opt.Options, args []string) error {

	ids, err := parseIDArgs(args)
	if err != nil {
		return err
	}

	if opts.Has("l") {
		req := &transmission.SetTorrentReq{
			Labels: strings.Split(opts.StringOr("l", ""), ","),
		}
		err := client.SetTorrents(context.Background(), ids, req)
		if err != nil {
			return err
		}
	}

	return nil
}
