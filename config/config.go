package config

import (
	"errors"
	"strings"

	"github.com/vchimishuk/config"
)

var spec = &config.Spec{
	Properties: []*config.PropertySpec{
		&config.PropertySpec{
			Type:    config.TypeString,
			Name:    "alias",
			Repeat:  true,
			Require: false,
			Parser: func(v any) (any, error) {
				return parseAlias(v.(string))
			},
		},
		&config.PropertySpec{
			Type:    config.TypeString,
			Name:    "host",
			Repeat:  false,
			Require: false,
		},
		&config.PropertySpec{
			Type:    config.TypeInt,
			Name:    "port",
			Repeat:  false,
			Require: false,
		},
	},
}

type Alias struct {
	Name    string
	Command string
	Args    string
}

type Config struct {
	Host    string
	Port    int
	Aliases []*Alias
}

func Parse(f string) (*Config, error) {
	cfg, err := config.ParseFile(spec, f)
	if err != nil {
		return nil, err
	}

	var aliases []*Alias
	for _, a := range cfg.Anys("alias") {
		aliases = append(aliases, a.(*Alias))
	}

	return &Config{
		Host:    cfg.StringOr("host", ""),
		Port:    cfg.IntOr("port", 0),
		Aliases: aliases,
	}, nil
}

func parseAlias(s string) (*Alias, error) {
	pts := strings.SplitN(s, "=", 2)
	if len(pts) != 2 {
		return nil, errors.New("invaid alias format")
	}

	name := strings.TrimSpace(pts[0])
	pts = strings.SplitN(strings.TrimSpace(pts[1]), " ", 2)
	cmd := pts[0]
	args := ""
	if len(pts) == 2 {
		args = pts[1]
	}

	return &Alias{
		Name:    name,
		Command: cmd,
		Args:    args,
	}, nil
}
