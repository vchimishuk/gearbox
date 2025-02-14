package main

import (
	"context"
	"fmt"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/format"
	"github.com/vchimishuk/opt"
)

type StatsCommand struct {
}

func NewStatsCommand() *StatsCommand {
	return &StatsCommand{}
}

func (c *StatsCommand) Name() string {
	return "stats"
}

func (c *StatsCommand) Usage() string {
	return c.Name() + " [-h host] [-p port]"
}

func (c *StatsCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		// {"P", "", opt.ArgNone, "", "do not automatically start torrent"},
	}
}

func (c *StatsCommand) Args() (int, int) {
	return 0, 0
}

func (c *StatsCommand) Exec(client *transmission.Client, opts opt.Options, args []string) error {
	st, err := client.GetSessionStats(context.Background())
	if err != nil {
		return err
	}

	names := []string{
		"Total torrents",
		"Active torrents",
		"Paused torrents",
		"Downloaded",
		"Uploaded",
		"Ratio",
		"Download rate",
		"Upload rate",
	}
	vals := []string{
		fmt.Sprintf("%d", st.Torrents),
		fmt.Sprintf("%d", st.ActiveTorrents),
		fmt.Sprintf("%d", st.PausedTorrents),
		format.Size(st.AllSessions.Downloaded),
		format.Size(st.AllSessions.Uploaded),
		fmt.Sprintf("%.1f", rate(st.AllSessions.Downloaded, st.AllSessions.Uploaded)),
		format.Rate(st.DownloadRate),
		format.Rate(st.UploadRate),
	}

	namesFmtr := format.NewColumnFormatter(true, names)
	for i := 0; i < len(names); i++ {
		fmt.Printf("%s: %s\n", namesFmtr.Format(names[i]), vals[i])
	}

	return nil
}

func rate(down int64, up int64) float64 {
	if down == 0 {
		return 0
	}

	return float64(up) / float64(down)
}
