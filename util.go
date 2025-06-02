package main

import (
	"fmt"
	"strconv"

	"github.com/pborzenkov/go-transmission/transmission"
)

func parseIDArgs(args []string) (transmission.IDList, error) {
	var ids transmission.IDList
	for _, a := range args {
		i, err := strconv.Atoi(a)
		if err != nil {
			return nil, fmt.Errorf("invalid torrent ID: %s", a)
		}
		ids = append(ids, transmission.ID(i))
	}

	return ids, nil
}
