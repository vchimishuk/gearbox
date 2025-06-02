package main

import (
	"context"
	"math"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/config"
	"github.com/vchimishuk/opt"
)

type StopCommand struct {
}

func NewStopCommand() *StopCommand {
	return &StopCommand{}
}

func (c *StopCommand) Name() string {
	return "stop"
}

func (c *StopCommand) Usage() string {
	return c.Name() + " id..."
}

func (c *StopCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c *StopCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *StopCommand) Exec(client *transmission.Client, cfg *config.Config,
	opts opt.Options, args []string) error {

	ids, err := parseIDArgs(args)
	if err != nil {
		return err
	}

	err = client.StopTorrents(context.Background(), ids)
	if err != nil {
		return err
	}

	return nil
}
