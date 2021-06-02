package beater

import (
	"errors"
	"fmt"
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/config"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"time"
)

// Logbeat is a beater object.
type Logbeat struct {
	period  time.Duration
	client  beat.Client
	handler Handler
	done    chan struct{}
}

func New() beat.Creator {
	return func(b *beat.Beat, c *common.Config) (beat.Beater, error) {
		return newBeater(b, c)
	}
}

func newBeater(_ *beat.Beat, rawConfig *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := rawConfig.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	//client, err := b.Publisher.Connect()
	//if err != nil {
	//	return nil, fmt.Errorf("Error connecting publisher: %v", err)
	//}

	// initialize Handler
	handler, err := NewHandler(c)
	if err != nil {
		return nil, err
	}

	lb := &Logbeat{
		period:  c.Period,
		handler: handler,
	}
	return beat.Beater(lb), nil
}

func (lb *Logbeat) Run(b *beat.Beat) error {
	// Judge whether it is running or not
	var err error
	if lb.done != nil {
		err = errors.New("Logbeat is already running!")
		logp.Warn("Logbeat.Run() method is repeatedly executed!")
		return err
	}

	// Initialize
	logp.Info("Logbeat is running! Hit CTRL-C to stop it.")

	lb.done = make(chan struct{})
	defer func() {
		lb.done = nil
	}()

	lb.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(lb.period)

	// Run
	lb.handler.Run()
	defer func(handler Handler) {
		err := handler.Stop()
		if err != nil {
			logp.Err("Failed to stop Handler: %v", err)
		}
	}(lb.handler)

	for {
		select {
		case <-lb.done: // Receive stop signal
			return nil
		case <-ticker.C:
		}

		events, err := lb.handler.Handle2Events()
		if err != nil {
			continue
		}
		lb.client.PublishAll(events)
	}

}

func (lb *Logbeat) Stop() {
	if lb.client != nil {
		err := lb.client.Close()
		if err != nil {
			logp.Err("Stop Logbeat error: %v", err)
		}
	}
	close(lb.done)
}
