package beater

import (
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/config"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
)

// Logbeat is a beater object.
type Logbeat struct {
	inputFile string
	config    config.Config
	events    beat.Client
	done      chan struct{}
}

func newBeater(b *beat.Beat, rawConfig *common.Config) (beat.Beater, error) {

}

func (lb *Logbeat) Run(b *beat.Beat) error {
	panic("implement me")
}

func (lb *Logbeat) Stop() {
	panic("implement me")
}
