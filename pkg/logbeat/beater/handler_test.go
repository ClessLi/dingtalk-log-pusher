package beater

import (
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/filter"
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/output_formatter"
	"github.com/elastic/beats/libbeat/beat"
	"reflect"
	"sync"
	"testing"
)

func TestHandlerImpl_Handle2Events(t *testing.T) {
	type fields struct {
		units []handlerUnit
	}
	tests := []struct {
		name    string
		fields  fields
		want    []beat.Event
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{units: []handlerUnit{new(mockNormalHandlerUnit)}},
			want: []beat.Event{{
				Timestamp: mockTimestamp,
				Fields:    mockMapStr,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerImpl{
				units: tt.fields.units,
			}
			got, err := h.Handle2Events()
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle2Events() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle2Events() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestHandlerImpl_Run(t *testing.T) {
//	type fields struct {
//		units []handlerUnit
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &HandlerImpl{
//				units: tt.fields.units,
//			}
//		})
//	}
//}

func TestHandlerImpl_Stop(t *testing.T) {
	type fields struct {
		units []handlerUnit
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{units: []handlerUnit{new(mockNormalHandlerUnit)}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerImpl{
				units: tt.fields.units,
			}
			if err := h.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestNewHandler(t *testing.T) {
//	type args struct {
//		cfg config.Config
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    Handler
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NewHandler(tt.args.cfg)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NewHandler() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewHandler() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestNewHandlerUnit(t *testing.T) {
//	type args struct {
//		inputInfo config.Input
//		formatStr string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    handlerUnit
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NewHandlerUnit(tt.args.inputInfo, tt.args.formatStr)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NewHandlerUnit() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewHandlerUnit() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func Test_handlerUnitImpl_Handle(t *testing.T) {
	type fields struct {
		name      string
		inputPath string
		fields    map[string]string
		filters   filter.Filters
		formatter output_formatter.OutputFormatter
		cache     []map[string]string
		lock      *sync.Mutex
		done      chan struct{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    []beat.Event
		wantErr bool
	}{
		{
			name: "normal test",
			fields: fields{
				name:      "test",
				inputPath: ".",
				fields:    mockFields,
				filters:   new(mockNormalFilters),
				formatter: new(mockOutputFormatter),
				cache:     []map[string]string{mockFilteredMsgs[0], mockFilteredMsgs[1]},
				lock:      new(sync.Mutex),
				done:      make(chan struct{}),
			},
			want: []beat.Event{mockEvent1, mockEvent2},
		},
		{
			name: "one message",
			fields: fields{
				name:      "test",
				inputPath: ".",
				fields:    mockFields,
				filters:   new(mockNormalFilters),
				formatter: new(mockOutputFormatter),
				cache:     []map[string]string{mockFilteredMsgs[0]},
				lock:      new(sync.Mutex),
				done:      make(chan struct{}),
			},
			want: []beat.Event{mockEvent1},
		},
		{
			name: "no message cached",
			fields: fields{
				name:      "test",
				inputPath: ".",
				fields:    mockFields,
				filters:   new(mockNormalFilters),
				formatter: new(mockOutputFormatter),
				cache:     []map[string]string{},
				lock:      new(sync.Mutex),
				done:      make(chan struct{}),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &handlerUnitImpl{
				name:      tt.fields.name,
				inputPath: tt.fields.inputPath,
				fields:    tt.fields.fields,
				filters:   tt.fields.filters,
				formatter: tt.fields.formatter,
				cache:     tt.fields.cache,
				lock:      tt.fields.lock,
				done:      tt.fields.done,
			}
			got, err := u.Handle()
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_handlerUnitImpl_Run(t *testing.T) {
//	type fields struct {
//		name      string
//		inputPath string
//		fields    map[string]string
//		filters   filter.Filters
//		formatter output_formatter.OutputFormatter
//		cache     []map[string]string
//		lock      *sync.Mutex
//		done      chan struct{}
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := &handlerUnitImpl{
//				name:      tt.fields.name,
//				inputPath: tt.fields.inputPath,
//				fields:    tt.fields.fields,
//				filters:   tt.fields.filters,
//				formatter: tt.fields.formatter,
//				cache:     tt.fields.cache,
//				lock:      tt.fields.lock,
//				done:      tt.fields.done,
//			}
//		})
//	}
//}

//func Test_handlerUnitImpl_Stop(t *testing.T) {
//	type fields struct {
//		name      string
//		inputPath string
//		fields    map[string]string
//		filters   filter.Filters
//		formatter output_formatter.OutputFormatter
//		cache     []map[string]string
//		lock      *sync.Mutex
//		done      chan struct{}
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		{
//			name: "normal test",
//			fields: fields{
//				name:      "test",
//				inputPath: ".",
//				fields:    mockFields,
//				filters:   new(mockNormalFilters),
//				formatter: new(mockOutputFormatter),
//				cache:     []map[string]string{},
//				lock:      new(sync.Mutex),
//				done:      make(chan struct{}),
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := &handlerUnitImpl{
//				name:      tt.fields.name,
//				inputPath: tt.fields.inputPath,
//				fields:    tt.fields.fields,
//				filters:   tt.fields.filters,
//				formatter: tt.fields.formatter,
//				cache:     tt.fields.cache,
//				lock:      tt.fields.lock,
//				done:      tt.fields.done,
//			}
//			if err := u.Stop(); (err != nil) != tt.wantErr {
//				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

//func Test_handlerUnitImpl_UnitName(t *testing.T) {
//	type fields struct {
//		name      string
//		inputPath string
//		fields    map[string]string
//		filters   filter.Filters
//		formatter output_formatter.OutputFormatter
//		cache     []map[string]string
//		lock      *sync.Mutex
//		done      chan struct{}
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   string
//	}{
//		{
//			name: "normal test",
//
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := handlerUnitImpl{
//				name:      tt.fields.name,
//				inputPath: tt.fields.inputPath,
//				fields:    tt.fields.fields,
//				filters:   tt.fields.filters,
//				formatter: tt.fields.formatter,
//				cache:     tt.fields.cache,
//				lock:      tt.fields.lock,
//				done:      tt.fields.done,
//			}
//			if got := u.UnitName(); got != tt.want {
//				t.Errorf("UnitName() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_handlerUnitImpl_watchFile(t *testing.T) {
//	type fields struct {
//		name      string
//		inputPath string
//		fields    map[string]string
//		filters   filter.Filters
//		formatter output_formatter.OutputFormatter
//		cache     []map[string]string
//		lock      *sync.Mutex
//		done      chan struct{}
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			u := &handlerUnitImpl{
//				name:      tt.fields.name,
//				inputPath: tt.fields.inputPath,
//				fields:    tt.fields.fields,
//				filters:   tt.fields.filters,
//				formatter: tt.fields.formatter,
//				cache:     tt.fields.cache,
//				lock:      tt.fields.lock,
//				done:      tt.fields.done,
//			}
//		})
//	}
//}

//func Test_newHandler(t *testing.T) {
//	type args struct {
//		units []handlerUnit
//	}
//	tests := []struct {
//		name string
//		args args
//		want Handler
//	}{
//		{
//			name: "normal test",
//			args: args{units: []handlerUnit{new(mockNormalHandlerUnit)}},
//			want: &HandlerImpl{units: []handlerUnit{new(mockNormalHandlerUnit)}},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := newHandler(tt.args.units); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("newHandler() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_newHandlerUnit(t *testing.T) {
//	type args struct {
//		inputPath string
//		fields    map[string]string
//		formatter output_formatter.OutputFormatter
//		filters   filter.Filters
//	}
//	tests := []struct {
//		name string
//		args args
//		want handlerUnit
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := newHandlerUnit(tt.args.inputPath, tt.args.fields, tt.args.formatter, tt.args.filters); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("newHandlerUnit() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
