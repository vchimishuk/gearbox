package main

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/format"
	"github.com/vchimishuk/opt"
)

type ListCommand struct {
}

func NewListCommand() *ListCommand {
	return &ListCommand{}
}

func (c *ListCommand) Name() string {
	return "list"
}

func (c *ListCommand) Usage() string {
	return c.Name() + " [-ar] [-c column] [-h host] [-n count] [-p port] [-s column]"
}

func (c *ListCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"a", "", opt.ArgNone, "", "display only recently active torrents"},
		{"c", "", opt.ArgString, "columns", "list of columns to display"},
		// TODO: {"H", "", opt.ArgNone, "", "do not display column header"},
		{"n", "", opt.ArgInt, "count", "display only N first rows"},
		{"r", "", opt.ArgNone, "", "sort in reverse order"},
		{"s", "", opt.ArgString, "column", "column to sort results by"},
	}
}

func (c *ListCommand) Args() (int, int) {
	return 0, 0
}

func (c *ListCommand) Exec(client *transmission.Client, opts opt.Options, args []string) error {
	cols, err := parseColumns(opts.StringOr("c", "id,name"))
	if err != nil {
		return err
	}

	fields := ColumnsToFields(cols)

	var sortCol Column
	if opts.Has("s") {
		var err error
		sortCol, err = GetColumn(opts.StringOr("s", ""))
		if err != nil {
			return err
		}
		fields = append(fields, sortCol.Field())
	}

	id := transmission.All()
	if opts.Has("a") {
		id = transmission.RecentlyActive()
	}
	trs, err := client.GetTorrents(context.Background(), id, fields...)
	if err != nil {
		return err
	}

	if sortCol != nil {
		c := sortCol.Comparator()
		if opts.Has("r") {
			c = reversed(c)
		}
		slices.SortFunc(trs, c)
	}

	count := min(len(trs), opts.IntOr("n", math.MaxInt))
	colVals := make([][]string, len(cols))
	for i, c := range cols {
		colVals[i] = make([]string, count+1)
		colVals[i][0] = c.ListTitle()
	}
	for i := 0; i < count; i++ {
		for j, c := range cols {
			colVals[j][i+1] = c.Value(trs[i])
		}
	}

	fmtrs := make([]format.ColumnFormatter, len(cols))
	for i, c := range cols {
		fmtrs[i] = c.Formatter(colVals[i])
	}

	for i := 0; i < count+1; i++ {
		for j, _ := range cols {
			if j > 0 {
				fmt.Print("  ")
			}
			if j == len(cols)-1 && !fmtrs[j].Right() {
				fmt.Print(colVals[j][i])
			} else {
				fmt.Print(fmtrs[j].Format(colVals[j][i]))
			}
		}
		fmt.Println()
	}

	return nil
}

func parseColumns(s string) ([]Column, error) {
	var cols []Column

	for _, c := range strings.Split(s, ",") {
		col, err := GetColumn(c)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}

	return cols, nil
}

func reversed[T any](f func(a, b T) int) func(a, b T) int {
	return func(a, b T) int {
		return f(a, b) * -1
	}
}
