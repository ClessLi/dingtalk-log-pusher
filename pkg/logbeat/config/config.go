package config

import (
	"github.com/elastic/beats/libbeat/common"
	"time"
)

type Config struct {
	Period       time.Duration `config:"period"`
	Inputs       []Input       `config:"inputs"`
	OutputFormat string        `config:"output_format"` // ${Env}: Config.Inputs[*].Fields.Env, ${Title}: CInputs[*].Fields.Title, ${Date}: Filter.Date, ${Msg}: Filter.Msg
}

var DefaultConfig = Config{
	Period:       1 * time.Second,
	OutputFormat: "${Env}环境-${Title}日志：\n\t日志时间：${Date}\n\t日志信息：${Msg}",
}

type Input struct {
	*common.Config
	Filters []Filter `config:"filters"`
}

type Filter struct {
	// TODO: Encoding
	Keyword      string `config:"keyword"`
	DateReg      string `config:"date_reg"`
	TimeTemplate string `config:"time_template"`
}
