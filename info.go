package main

import (
	"context"
	"fmt"
	"math"
	"sort"
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

type InfoCommand struct {
}

func NewInfoCommand() *InfoCommand {
	return &InfoCommand{}
}

func (c *InfoCommand) Name() string {
	return "info"
}

func (c *InfoCommand) Usage() string {
	return c.Name() + " [-fi] id..."
}

func (c *InfoCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"f", "", opt.ArgNone, "", "show torrent files list"},
		{"i", "", opt.ArgNone, "", "show torrent information"},
	}
}

func (c *InfoCommand) Args() (int, int) {
	return 1, math.MaxInt
}

func (c *InfoCommand) Exec(client *transmission.Client, cfg *config.Config,
	opts opt.Options, args []string) error {

	ids, err := parseIDArgs(args)
	if err != nil {
		return err
	}

	info := opts.Has("i") || !opts.Has("f")
	files := opts.Has("f")

	var cols []transmission.TorrentField
	if info {
		cols = append(cols, ColumnsToFields(infoColumns)...)
	}
	if files {
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

		if info {
			for _, col := range infoColumns {
				fmt.Print(strings.Repeat(" ",
					maxTitle-len(col.Title())))
				fmt.Print(col.Title())
				fmt.Print(": ")
				fmt.Println(col.Value(t))
			}
		}

		if files {
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

	return nil
}
