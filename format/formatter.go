package format

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ColumnFormatter interface {
	Right() bool
	Format(s string) string
}

type columnFormatter struct {
	right  bool
	format func(s string) string
}

func (f *columnFormatter) Right() bool {
	return f.right
}

func (f *columnFormatter) Format(s string) string {
	return f.format(s)
}

func NewColumnFormatter(right bool, rows []string) ColumnFormatter {
	sz := 0
	for _, r := range rows {
		n := utf8.RuneCountInString(r)
		if n > sz {
			sz = n
		}
	}

	if right {
		return &columnFormatter{
			true,
			func(s string) string {
				idnt := strings.Repeat(" ", sz-utf8.RuneCountInString(s))
				return fmt.Sprintf("%s%s", idnt, s)
			},
		}
	} else {
		return &columnFormatter{
			false,
			func(s string) string {
				idnt := strings.Repeat(" ", sz-utf8.RuneCountInString(s))
				return fmt.Sprintf("%s%s", s, idnt)
			},
		}
	}
}
