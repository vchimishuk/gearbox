package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/config"
	"github.com/vchimishuk/gearbox/format"
	"github.com/vchimishuk/opt"
)

var infoColumns []Column = []Column{
	GetColumnMust("id"),
	GetColumnMust("name"),
	GetColumnMust("labels"),
	GetColumnMust("status"),
	GetColumnMust("size"),
	GetColumnMust("dsize"),
	GetColumnMust("usize"),
	GetColumnMust("ratio"),
	GetColumnMust("drate"),
	GetColumnMust("urate"),
	GetColumnMust("active"),
	GetColumnMust("added"),
	GetColumnMust("created"),
	GetColumnMust("comment"),
}

type TorrentCommand struct {
}

func NewTorrentCommand() *TorrentCommand {
	return &TorrentCommand{}
}

func (c *TorrentCommand) Name() string {
	return "torrent"
}

func (c *TorrentCommand) Usage() string {
	return c.Name() + " [-DdfiSs] [-h host] [-l labels] [-p port] id..."
}

func (c *TorrentCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"D", "", opt.ArgNone, "", "delete also torrent data when deleting torrent"},
		{"S", "", opt.ArgNone, "", "stop torrent"},
		{"d", "", opt.ArgNone, "", "delete torrent"},
		{"f", "", opt.ArgNone, "", "show torrent files list"},
		{"i", "", opt.ArgNone, "", "show torrent information"},
		{"l", "", opt.ArgString, "", "set labels for torrent"},
		{"s", "", opt.ArgNone, "", "start torrent"},
	}
}

func (c *TorrentCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *TorrentCommand) Exec(client *transmission.Client, cfg *config.Config,
	opts opt.Options, args []string) error {

	if opts.Has("S") && opts.Has("s") {
		return errors.New("either -S or -s can be used")
	}
	if (opts.Has("S") || opts.Has("s") || opts.Has("l")) && opts.Has("r") {
		return errors.New("-r cannot be used simultaneously with -S -s or -l")
	}
	if opts.Has("D") && !opts.Has("d") {
		return errors.New("-D requires -d")
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

	if opts.Has("i") || opts.Has("f") {
		var cols []transmission.TorrentField

		if opts.Has("i") {
			cols = append(cols, ColumnsToFields(infoColumns)...)
		}
		if opts.Has("f") {
			cols = append(cols, transmission.TorrentFieldFiles)
		}
		trs, err := client.GetTorrents(context.Background(), ids,
			cols...)
		if err != nil {
			return err
		}

		maxTitle := 0
		for _, c := range infoColumns {
			maxTitle = max(maxTitle, len(c.Title()))
		}
		for i, t := range trs {
			if i != 0 {
				fmt.Println()
			}

			if opts.Has("i") {
				for _, col := range infoColumns {
					fmt.Print(strings.Repeat(" ",
						maxTitle-len(col.Title())))
					fmt.Print(col.Title())
					fmt.Print(": ")
					fmt.Println(col.Value(t))
				}
			}

			if opts.Has("f") {
				if opts.Has("i") {
					fmt.Println()
				}
				files := t.Files
				sort.Slice(files, func(i, j int) bool {
					return files[i].Name < files[j].Name
				})

				cols := make([][]string, 2)
				cols[0] = make([]string, len(files)+1)
				cols[0][0] = "FILE"
				cols[1] = make([]string, len(files)+1)
				cols[1][0] = "SIZE"
				for i, f := range files {
					cols[0][i+1] = f.Name
					cols[1][i+1] = format.Size(f.Size)
				}

				fileFmtr := format.NewColumnFormatter(false,
					cols[0])
				sizeFmtr := format.NewColumnFormatter(true,
					cols[1])

				for i := 0; i < len(files)+1; i++ {
					fmt.Print(fileFmtr.Format(cols[0][i]))
					fmt.Print("  ")
					fmt.Print(sizeFmtr.Format(cols[1][i]))
					fmt.Println()
				}
			}
		}
	}

	return nil
}
