package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/opt"
)

type TorrentCommand struct {
}

func NewTorrentCommand() *TorrentCommand {
	return &TorrentCommand{}
}

func (c *TorrentCommand) Name() string {
	return "torrent"
}

func (c *TorrentCommand) Usage() string {
	return c.Name() + " [-DdSs] [-h host] [-l labels] [-p port] id..."
}

func (c *TorrentCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"S", "", opt.ArgNone, "", "stop torrent"},
		{"D", "", opt.ArgNone, "", "delete also torrent data when deleting torrent"},
		{"d", "", opt.ArgNone, "", "delete torrent"},
		{"l", "", opt.ArgString, "", "set labels for torrent"},
		{"s", "", opt.ArgNone, "", "start torrent"},
	}
}

func (c *TorrentCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *TorrentCommand) Exec(client *transmission.Client, opts opt.Options, args []string) error {
	if opts.Has("S") && opts.Has("s") {
		return errors.New("either -S or -s can be used")
	}
	if (opts.Has("S") || opts.Has("s") || opts.Has("l")) && opts.Has("r") {
		return errors.New("-r cannot be used simultaneously with -S -s or -l")
	}
	if opts.Has("d") && !opts.Has("r") {
		return errors.New("-d requires -r")
	}

	var ids transmission.IDList
	for _, a := range args {
		i, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("invalid torrent: %s", a)
		}
		ids = append(ids, transmission.ID(i))
	}

	if opts.Has("s") {
		err := client.StartTorrents(context.Background(), ids)
		if err != nil {
			return err
		}
	}

	if opts.Has("S") {
		err := client.StopTorrents(context.Background(), ids)
		if err != nil {
			return err
		}
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

	if opts.Has("d") {
		d := opts.Has("D")
		err := client.RemoveTorrents(context.Background(), ids, d)
		if err != nil {
			return err
		}
	}

	return nil
}
