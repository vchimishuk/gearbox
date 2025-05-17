package config

import "github.com/vchimishuk/config"

var spec = &config.Spec{
	Properties: []*config.PropertySpec{
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
		&config.PropertySpec{
			Type:    config.TypeString,
			Name:    "list-columns",
			Repeat:  false,
			Require: false,
		},
		&config.PropertySpec{
			Type:    config.TypeString,
			Name:    "list-sort",
			Repeat:  false,
			Require: false,
		},
		&config.PropertySpec{
			Type:    config.TypeBool,
			Name:    "list-reverse",
			Repeat:  false,
			Require: false,
		},
		&config.PropertySpec{
			Type:    config.TypeInt,
			Name:    "list-count",
			Repeat:  false,
			Require: false,
		},
	},
}

type Config struct {
	Host        string
	Port        int
	ListColumns string
	ListSort    string
	ListReverse bool
	ListCount   int
}

func Parse(f string) (*Config, error) {
	cfg, err := config.ParseFile(spec, f)
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:        cfg.StringOr("host", ""),
		Port:        cfg.IntOr("port", 0),
		ListColumns: cfg.StringOr("list-columns", ""),
		ListSort:    cfg.StringOr("list-sort", ""),
		ListReverse: cfg.BoolOr("list-reverse", false),
		ListCount:   cfg.IntOr("list-count", 0),
	}, nil
}
