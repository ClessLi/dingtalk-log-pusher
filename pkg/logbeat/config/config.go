package config

import (
	"github.com/elastic/beats/libbeat/common"
	"time"
)

type Config struct {
	Inputs          []Input       `config:"inputs"`
	ShutdownTimeout time.Duration `config:"shutdown_timeout"`
}

type Input struct {
	*common.Config
	Filters []Filter
}
