package beater

import (
	"errors"
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"reflect"
	"regexp"
	"sync"
	"time"
)

////// mock for testing Handler and handlerUnit
var (
	changeSpecialChar2Dot = func(str string) string {
		reg, _ := regexp.Compile(`[^0-9a-zA-Z]+`)
		ret := reg.ReplaceAllString(str, ".")
		return ret
	}
	mockTimestamp    = time.Now()
	mockMapStr       = common.MapStr{"Msg": "normal"}
	mockFields       = map[string]string{"Env": "UAT", "Title": "测试应用"}
	mockFilteredMsgs = []map[string]string{
		{"Env": "UAT", "Title": "测试应用", "Date": "[2021-06-01 10:43:50,253]", "Msg": "log message for <Test1>, have fun!!", "TimeTemplate": "[2006-01-02 15:04:05,000]"},
		{"Env": "UAT", "Title": "测试应用", "Date": "2020/06/01 13:50:30.020", "Msg": "log message for <Test2>, have fun too!!", "TimeTemplate": "2006/01/02 15:04:05.000"},
	}
	mockMsg1           = " [2021-06-01 10:43:50,253] log message for <Test1>, have fun!!     "
	mockMsg2           = " 2020/06/01  13:50:30.020    log message for <Test2>, have fun too!!"
	mockFiltersResult1 = map[string]string{"Date": "[2021-06-01 10:43:50,253]", "Msg": "log message for <Test1>, have fun!!", "TimeTemplate": "[2006-01-02 15:04:05,000]"}
	mockFiltersResult2 = map[string]string{"Date": "2020/06/01 13:50:30.020", "Msg": "log message for <Test2>, have fun too!!", "TimeTemplate": "2006/01/02 15:04:05.000"}

	mockFormatStr     = "%v环境-%v日志：\n\t日志时间：%v\n\t日志信息：%v"
	mockTimestamp1, _ = time.Parse(changeSpecialChar2Dot(mockFilteredMsgs[0]["TimeTemplate"]), changeSpecialChar2Dot(mockFilteredMsgs[0]["Date"]))
	mockTimestamp2, _ = time.Parse(changeSpecialChar2Dot(mockFilteredMsgs[1]["TimeTemplate"]), changeSpecialChar2Dot(mockFilteredMsgs[1]["Date"]))
	mockEvent1        = beat.Event{
		Timestamp: mockTimestamp1,
		Fields:    common.MapStr{"Msg": fmt.Sprintf(mockFormatStr, mockFilteredMsgs[0]["Env"], mockFilteredMsgs[0]["Title"], mockFilteredMsgs[0]["Date"], mockFilteredMsgs[0]["Msg"])},
	}
	mockEvent2 = beat.Event{
		Timestamp: mockTimestamp2,
		Fields:    common.MapStr{"Msg": fmt.Sprintf(mockFormatStr, mockFilteredMsgs[1]["Env"], mockFilteredMsgs[1]["Title"], mockFilteredMsgs[1]["Date"], mockFilteredMsgs[1]["Msg"])},
	}
)

type mockNormalHandlerUnit struct {
}

func (m mockNormalHandlerUnit) Run() {
}

func (m mockNormalHandlerUnit) Handle() ([]beat.Event, error) {
	return []beat.Event{{
		Timestamp: mockTimestamp,
		Fields:    mockMapStr,
	}}, nil
}

func (m mockNormalHandlerUnit) Stop() error {
	return nil
}

func (m mockNormalHandlerUnit) UnitName() string {
	return "mock_normal_handler_unit"
}

type mockNormalFilters struct {
}

func (m mockNormalFilters) Filter(msg string) (map[string]string, bool) {
	switch msg {
	case mockMsg1:
		return mockFiltersResult1, true
	case mockMsg2:
		return mockFiltersResult2, true
	default:
		return nil, false
	}
}

type mockOutputFormatter struct {
}

func (m mockOutputFormatter) Format(msgs ...map[string]string) ([]beat.Event, error) {
	if msgs == nil || len(msgs) < 1 {
		return nil, errors.New("messages is null")
	}
	events := make([]beat.Event, 0)
	for _, msg := range msgs {
		switch {
		case reflect.DeepEqual(msg, mockFilteredMsgs[0]):
			events = append(events, mockEvent1)
		case reflect.DeepEqual(msg, mockFilteredMsgs[1]):
			events = append(events, mockEvent2)
		}
	}
	if len(events) < 1 {
		return nil, fmt.Errorf("Failed to format message: %v", msgs)
	}
	return events, nil
}

////// mock for testing Logbeat

var (
	mockBeat = beat.Beat{
		Info:      beat.Info{},
		Publisher: new(mockBeatPublisher),
	}
)

type mockBeatPublisher struct {
}

func (m mockBeatPublisher) ConnectWith(_ beat.ClientConfig) (beat.Client, error) {
	return m.Connect()
}

func (m mockBeatPublisher) Connect() (beat.Client, error) {
	return new(mockBeatClient), nil
}

func (m mockBeatPublisher) SetACKHandler(_ beat.PipelineACKHandler) error {
	fmt.Println("do nothing for mockBeatPublisher.SetACKHandler()")
	return nil
}

type mockBeatClient struct {
}

func (m mockBeatClient) Publish(event beat.Event) {
	fmt.Printf("publish <- '{\"@timestamp\": \"%v\", \"Fields\": %v}'\n", event.Timestamp, event.Fields)
}

func (m mockBeatClient) PublishAll(events []beat.Event) {
	if events != nil && len(events) > 0 {
		for _, event := range events {
			m.Publish(event)
		}
	}
}

func (m mockBeatClient) Close() error {
	fmt.Println("beat client stopping, plz wait one second")
	time.Sleep(time.Second)
	fmt.Println("beat client is stopped")
	return nil
}

type mockHandler struct {
	once *sync.Once
}

func newMockHandler() Handler {
	return Handler(&mockHandler{once: new(sync.Once)})
}

func (m mockHandler) Run() {
	fmt.Println("Handler is running!")
}

func (m mockHandler) Handle2Events() ([]beat.Event, error) {
	events := make([]beat.Event, 0)
	m.once.Do(func() {
		events = append(events, mockEvent1, mockEvent2)
	})
	return events, nil
}

func (m mockHandler) Stop() error {
	fmt.Println("Handler is stopped")
	return nil
}
