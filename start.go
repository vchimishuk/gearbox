package main

import (
	"context"
	"math"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/config"
	"github.com/vchimishuk/opt"
)

type StartCommand struct {
}

func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

func (c *StartCommand) Name() string {
	return "start"
}

func (c *StartCommand) Usage() string {
	return c.Name() + " id..."
}

func (c *StartCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c *StartCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *StartCommand) Exec(client *transmission.Client, cfg *config.Config,
	opts opt.Options, args []string) error {

	ids, err := parseIDArgs(args)
	if err != nil {
		return err
	}

	err = client.StartTorrents(context.Background(), ids)
	if err != nil {
		return err
	}

	return nil
}
