package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/format"
	"github.com/vchimishuk/opt"
)

type InfoCommand struct {
}

func NewInfoCommand() *InfoCommand {
	return &InfoCommand{}
}

func (c *InfoCommand) Name() string {
	return "info"
}

func (c *InfoCommand) Usage() string {
	return c.Name() + " [-f] [-h host] [-p port] id"
}

func (c *InfoCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"f", "", opt.ArgNone, "", "show files list instead of general info"},
		// TODO: {"H", "", opt.ArgNone, "", "do not display column header"},
	}
}

func (c *InfoCommand) Args() (int, int) {
	return 1, 1
}

func (c *InfoCommand) Exec(client *transmission.Client, opts opt.Options, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.New("invalid torrent ID")
	}

	if opts.Has("f") {
		return c.execFiles(client, id)
	} else {
		return c.execInfo(client, id)
	}
}

func (c *InfoCommand) execInfo(client *transmission.Client, id int) error {
	cols := []Column{
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

	trs, err := client.GetTorrents(context.Background(), transmission.ID(id),
		ColumnsToFields(cols)...)
	if err != nil {
		return err
	}

	maxTitle := 0
	for _, c := range cols {
		maxTitle = max(maxTitle, len(c.Title()))
	}
	for _, col := range cols {
		idnt := strings.Repeat(" ", maxTitle-len(col.Title()))
		fmt.Printf("%s%s: %s\n", idnt, col.Title(), col.Value(trs[0]))
	}

	return nil
}

func (c *InfoCommand) execFiles(client *transmission.Client, ID int) error {
	trs, err := client.GetTorrents(context.Background(), transmission.ID(ID),
		transmission.TorrentFieldFiles)
	if err != nil {
		return err
	}
	files := trs[0].Files

	cols := make([][]string, 2)
	cols[0] = make([]string, len(files)+1)
	cols[0][0] = "FILE"
	cols[1] = make([]string, len(files)+1)
	cols[1][0] = "SIZE"
	for i, f := range files {
		cols[0][i+1] = f.Name
		cols[1][i+1] = format.Size(f.Size)
	}

	fileFmtr := format.NewColumnFormatter(false, cols[0])
	sizeFmtr := format.NewColumnFormatter(true, cols[1])

	for i := 0; i < len(files)+1; i++ {
		fmt.Print(fileFmtr.Format(cols[0][i]))
		fmt.Print("  ")
		fmt.Print(sizeFmtr.Format(cols[1][i]))
		fmt.Println()
	}

	return nil
}
