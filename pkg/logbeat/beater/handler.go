package beater

import (
	"errors"
	"fmt"
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/config"
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/filter"
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/output_formatter"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

type Handler interface {
	Run()
	Handle2Events() ([]beat.Event, error)
	Stop() error
}

type HandlerImpl struct {
	units []handlerUnit
}

func NewHandler(cfg config.Config) (Handler, error) {
	// check inputs info
	inputsLen := len(cfg.Inputs)
	if inputsLen < 1 {
		return nil, errors.New("Config.Inputs is null")
	}

	// initialize log handler units
	units := make([]handlerUnit, 0)
	for _, inputInfo := range cfg.Inputs {
		unit, err := NewHandlerUnit(inputInfo, cfg.OutputFormat)
		if err != nil {
			logp.Warn("Failed to initialize log handler unit for log: %v", inputInfo.Path())
			continue
		}
		units = append(units, unit)
	}
	if len(units) < 1 {
		err := errors.New("Failed to initialize log handler due to unsuccessful initialization of log handler unit")
		return nil, err
	}

	return newHandler(units), nil
}

func newHandler(units []handlerUnit) Handler {
	return &HandlerImpl{
		units: units,
	}
}

func (h *HandlerImpl) Run() {
	for _, unit := range h.units {
		unit.Run()
	}
}

func (h *HandlerImpl) Handle2Events() ([]beat.Event, error) {
	events := make([]beat.Event, 0)
	for _, unit := range h.units {
		es, err := unit.Handle()
		if err != nil || es == nil || len(es) < 1 {
			continue
		}
		events = append(events, es...)
	}

	if len(events) < 1 {
		return nil, errors.New("None event")
	}

	return events, nil
}

func (h *HandlerImpl) Stop() error {
	unitStopFailedList := make([]string, 0)
	for _, unit := range h.units {
		err := unit.Stop()
		if err != nil {
			unitStopFailedList = append(unitStopFailedList, unit.UnitName())
		}

	}
	if len(unitStopFailedList) < 1 {
		return nil
	}
	return fmt.Errorf("Stop log handler error due to units '%v' stop timeout", unitStopFailedList)
}

type handlerUnit interface {
	Run()
	Handle() ([]beat.Event, error)
	Stop() error
	UnitName() string
}

type handlerUnitImpl struct {
	name      string
	inputPath string
	//tail *tail.Tail
	fields    map[string]string
	filters   filter.Filters
	formatter output_formatter.OutputFormatter
	cache     []map[string]string
	lock      *sync.Mutex
	done      chan struct{}
}

func NewHandlerUnit(inputInfo config.Input, formatStr string) (handlerUnit, error) {

	filters, err := filter.NewFilters(inputInfo.Filters)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	for _, fieldKey := range inputInfo.GetFields() {
		fieldValue, err := inputInfo.String(fieldKey, -1)
		if err != nil {
			return nil, err
		}
		fields[fieldKey] = fieldValue
	}

	// initialize output formatter
	formatter, err := output_formatter.NewOutputFormatter(formatStr)
	if err != nil {
		err = fmt.Errorf("Failed to initialize log handler: %v", err)
		return nil, err
	}

	return newHandlerUnit(inputInfo.Path(), fields, formatter, filters), nil
	//unit := newHandlerUnit(inputInfo.Path(), fields, formatter, filters)
	//return unit, nil
}

func newHandlerUnit(inputPath string, fields map[string]string, formatter output_formatter.OutputFormatter, filters filter.Filters) handlerUnit {
	return &handlerUnitImpl{
		name:      "handler_unit_of_" + inputPath,
		inputPath: inputPath,
		//tail:      t,
		fields:    fields,
		formatter: formatter,
		filters:   filters,
		cache:     make([]map[string]string, 0),
		lock:      new(sync.Mutex),
	}
}

func (u *handlerUnitImpl) Run() {
	if u.done != nil {
		return
	}
	u.done = make(chan struct{})
	defer func() {
		close(u.done)
		u.done = nil
	}()
	u.watchFile()
}

func (u *handlerUnitImpl) Handle() ([]beat.Event, error) {
	u.lock.Lock()
	defer u.lock.Unlock()

	if len(u.cache) < 1 {
		return nil, fmt.Errorf("none message")
	}

	msgs := make([]map[string]string, 0)
	//copy(msgs, u.cache)
	msgs = append(msgs, u.cache...)
	u.cache = u.cache[:0]

	events, err := u.formatter.Format(msgs...)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (u *handlerUnitImpl) Stop() error {
	select {
	case u.done <- struct{}{}:
	case <-time.After(10 * time.Second):
		err := fmt.Errorf("Stop log handler unit '%v' timeout", u.name)
		logp.Warn(err.Error())
		return err
	}
	return nil
}

func (u handlerUnitImpl) UnitName() string {
	return u.name
}

func (u *handlerUnitImpl) watchFile() {
	t, err := tail.TailFile(u.inputPath, tail.Config{
		//tails, err := tail.TailFile(filepath, tail.Config{
		ReOpen: true,
		Follow: true,
		// Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		logp.Err("Read log file '%v' error: %v", u.inputPath, err)
		return
	}
	defer func(t *tail.Tail) {
		err := t.Stop()
		if err != nil {
			logp.Err("Stop read log file '%v' error: %v", u.inputPath, err)
		}
	}(t)

	u.cache = u.cache[:0]
	var (
		line *tail.Line
	)

	for {
		select {
		case <-u.done:
			return
		case line = <-t.Lines:
		}
		// 筛选该行数据
		msg, isMatch := u.filters.Filter(line.Text)
		if !isMatch {
			continue
		}
		// 可能会存在配置的fields字段覆盖筛选结果
		for key, value := range u.fields {
			msg[key] = value
		}
		u.lock.Lock()
		u.cache = append(u.cache, msg)
		u.lock.Unlock()
	}
}
