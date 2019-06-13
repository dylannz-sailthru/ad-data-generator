package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type DistributionType int

const (
	DistrubutionRegular DistributionType = iota
	DistrubutionDisrupted
)

type Distribution struct {
	MinMean int `toml:"min_mean"`
	MaxMean int `toml:"max_mean"`
	MinStd  int `toml:"min_std"`
	MaxStd  int `toml:"max_std"`
}

type Config struct {
	Nodes                 int          `toml:"nodes"`
	Queries               int          `toml:"queries"`
	Metrics               int          `toml:"metrics"`
	Hours                 int          `toml:"hours"`
	Disruptions           int          `toml:"disruptions"`
	Threads               int          `toml:"threads"`
	RegularDistribution   Distribution `toml:"regular_distribution"`
	DisruptedDistribution Distribution `toml:"disrupted_distribution"`
	ES                    ES           `toml:"es"`
}

type ES struct {
	Mapping         string `toml:"mapping"`
	HotcloudMapping string `toml:"hotcloud_mapping"`
	Query           string `toml:"query"`
	BulkSize        int    `toml:"bulk_size"`
}

func NewFromTOML(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "open/read toml file")
	}

	cfg := &Config{}
	_, err = toml.Decode(string(b), cfg)
	return cfg, errors.Wrap(err, "decode toml")
}
