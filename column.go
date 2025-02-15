package main

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pborzenkov/go-transmission/transmission"
	"github.com/vchimishuk/gearbox/format"
)

const TimeLayout = time.RFC822Z

type Column interface {
	Name() string
	Title() string
	ListTitle() string
	Field() transmission.TorrentField
	Comparator() func(a, b *transmission.Torrent) int
	Value(t *transmission.Torrent) string
	Formatter(rows []string) format.ColumnFormatter
}

var Columns []Column = []Column{
	&column{
		"active",
		"Last active",
		"ACTIVE",
		transmission.TorrentFieldLastActiveAt,
		comparator(func(t *transmission.Torrent) int64 {
			return t.LastActiveAt.Unix()
		}),
		func(t *transmission.Torrent) string {
			return t.LastActiveAt.Format(TimeLayout)
		},
		formatter(false),
	},
	&column{
		"added",
		"Added",
		"ADDED",
		transmission.TorrentFieldAddedAt,
		comparator(func(t *transmission.Torrent) int64 {
			return t.AddedAt.Unix()
		}),
		func(t *transmission.Torrent) string {
			return t.AddedAt.Format(TimeLayout)
		},
		formatter(false),
	},
	&column{
		"comment",
		"Comment",
		"COMMENT",
		transmission.TorrentFieldComment,
		comparator(func(t *transmission.Torrent) string {
			return t.Comment
		}),
		func(t *transmission.Torrent) string {
			return t.Comment
		},
		formatter(false),
	},
	&column{
		"created",
		"Created",
		"CREATED",
		transmission.TorrentFieldCreatedAt,
		comparator(func(t *transmission.Torrent) int64 {
			return t.CreatedAt.Unix()
		}),
		func(t *transmission.Torrent) string {
			return t.CreatedAt.Format(TimeLayout)
		},
		formatter(false),
	},
	&column{
		"drate",
		"Download rate",
		"DRATE",
		transmission.TorrentFieldDownloadRate,
		comparator(func(t *transmission.Torrent) int64 {
			return t.DownloadRate
		}),
		func(t *transmission.Torrent) string {
			return format.Rate(t.DownloadRate)
		},
		formatter(true),
	},
	&column{
		"dsize",
		"Downloaded",
		"DSIZE",
		transmission.TorrentFieldDownloadedTotal,
		comparator(func(t *transmission.Torrent) int64 {
			return t.DownloadedTotal
		}),
		func(t *transmission.Torrent) string {
			return format.Size(t.DownloadedTotal)
		},
		formatter(true),
	},
	&column{
		"id",
		"Id",
		"ID",
		transmission.TorrentFieldID,
		comparator(func(t *transmission.Torrent) transmission.ID {
			return t.ID
		}),
		func(t *transmission.Torrent) string {
			return strconv.Itoa(int(t.ID))
		},
		formatter(true),
	},
	&column{
		"labels",
		"Labels",
		"LABELS",
		transmission.TorrentFieldLabels,
		comparator(func(t *transmission.Torrent) string {
			return strings.Join(t.Labels, ",")
		}),
		func(t *transmission.Torrent) string {
			return strings.Join(t.Labels, ",")
		},
		formatter(false),
	},
	&column{
		"name",
		"Name",
		"NAME",
		transmission.TorrentFieldName,
		comparator(func(t *transmission.Torrent) string {
			return t.Name
		}),
		func(t *transmission.Torrent) string {
			return t.Name
		},
		formatter(false),
	},
	&column{
		"ratio",
		"Ratio",
		"RATIO",
		transmission.TorrentFieldUploadRatio,
		comparator(func(t *transmission.Torrent) float64 {
			return t.UploadRatio
		}),
		func(t *transmission.Torrent) string {
			return fmt.Sprintf("%.2f", t.UploadRatio)
		},
		formatter(true),
	},
	&column{
		"size",
		"Size",
		"SIZE",
		transmission.TorrentFieldTotalSize,
		comparator(func(t *transmission.Torrent) int64 {
			return t.TotalSize
		}),
		func(t *transmission.Torrent) string {
			return format.Size(t.TotalSize)
		},
		formatter(true),
	},
	&column{
		"status",
		"Status",
		"STATUS",
		transmission.TorrentFieldStatus,
		comparator(func(t *transmission.Torrent) transmission.Status {
			return t.Status
		}),
		func(t *transmission.Torrent) string {
			return t.Status.String()
		},
		formatter(false),
	},
	&column{
		"urate",
		"Upload rate",
		"URATE",
		transmission.TorrentFieldUploadRate,
		comparator(func(t *transmission.Torrent) int64 {
			return t.UploadRate
		}),
		func(t *transmission.Torrent) string {
			return format.Rate(t.UploadRate)
		},
		formatter(true),
	},
	&column{
		"usize",
		"Uploaded",
		"USIZE",
		transmission.TorrentFieldUploadedTotal,
		comparator(func(t *transmission.Torrent) int64 {
			return t.UploadedTotal
		}),
		func(t *transmission.Torrent) string {
			return format.Size(t.UploadedTotal)
		},
		formatter(true),
	},
}

type column struct {
	name       string
	title      string
	listTitle  string
	field      transmission.TorrentField
	comparator func(a, b *transmission.Torrent) int
	value      func(t *transmission.Torrent) string
	formatter  func(rows []string) format.ColumnFormatter
}

func (c *column) Name() string {
	return c.name
}

func (c *column) Title() string {
	return c.title
}

func (c *column) ListTitle() string {
	return c.listTitle
}

func (c *column) Field() transmission.TorrentField {
	return c.field
}

func (c *column) Comparator() func(a, b *transmission.Torrent) int {
	return c.comparator
}

func (c *column) Value(t *transmission.Torrent) string {
	return c.value(t)
}

func (c *column) Formatter(rows []string) format.ColumnFormatter {
	return c.formatter(rows)
}

func GetColumnMust(name string) Column {
	c, err := GetColumn(name)
	if err != nil {
		panic(err)
	}

	return c
}

func GetColumn(name string) (Column, error) {
	i := slices.IndexFunc(Columns, func(c Column) bool {
		return c.Name() == name
	})
	if i == -1 {
		return nil, errors.New("invalid column: " + name)
	}

	return Columns[i], nil
}

func ColumnsToFields(cols []Column) []transmission.TorrentField {
	fs := make([]transmission.TorrentField, len(cols))
	for i, c := range cols {
		fs[i] = c.Field()
	}

	return fs
}

func comparator[T cmp.Ordered](f func(t *transmission.Torrent) T) func(a, b *transmission.Torrent) int {
	return func(a, b *transmission.Torrent) int {
		return cmp.Compare(f(a), f(b))
	}
}

func formatter(right bool) func(rows []string) format.ColumnFormatter {
	return func(rows []string) format.ColumnFormatter {
		return format.NewColumnFormatter(right, rows)
	}
}
