package main

import (
	"context"
	"math"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/config"
	"github.com/vchimishuk/opt"
)

type DeleteCommand struct {
}

func NewDeleteCommand() *DeleteCommand {
	return &DeleteCommand{}
}

func (c *DeleteCommand) Name() string {
	return "delete"
}

func (c *DeleteCommand) Usage() string {
	return c.Name() + " [-d] id..."
}

func (c *DeleteCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"d", "", opt.ArgNone, "", "delete also torrent data"},
	}
}

func (c *DeleteCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *DeleteCommand) Exec(client *transmission.Client, cfg *config.Config,
	opts opt.Options, args []string) error {

	ids, err := parseIDArgs(args)
	if err != nil {
		return err
	}

	d := opts.Has("d")
	err = client.RemoveTorrents(context.Background(), ids, d)
	if err != nil {
		return err
	}

	return nil
}
