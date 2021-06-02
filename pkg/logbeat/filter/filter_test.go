package filter

import (
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/config"
	"reflect"
	"regexp"
	"testing"
)

func TestFilter_Filter(t *testing.T) {
	tem1 := "[2006-01-02 15:04:05,000]"
	reg1, err := regexp.Compile(`^\s*(\[\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2},\d{3}\])\s*(.*<Test1>.*?)\s*$`)
	if err != nil {
		t.Fatal(err)
	}

	msg1 := " [2021-06-01 10:43:50,020] log message for <Test1>, have fun!!     "

	tem2 := "2006/01/02 15:04:05.000"
	reg2, err := regexp.Compile(`^\s*(\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2}.\d{3})\s*(.*<Test2>.*?)\s*$`)
	if err != nil {
		t.Fatal(err)
	}

	msg2 := " 2020/06/01  11:43:50.020    log message for <Test2>, have fun too!!"
	msgOther := "other message for <Test1> <Test2>, have fun..."

	type fields struct {
		reg          *regexp.Regexp
		timeTemplate string
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
		want1  bool
	}{
		{
			name: "normal test 1",
			fields: fields{
				reg:          reg1,
				timeTemplate: tem1,
			},
			args:  args{msg: msg1},
			want:  map[string]string{"TimeTemplate": tem1, "Date": "[2021-06-01 10:43:50,020]", "Msg": "log message for <Test1>, have fun!!"},
			want1: true,
		},
		{
			name: "normal test 2",
			fields: fields{
				reg:          reg2,
				timeTemplate: tem2,
			},
			args:  args{msg: msg2},
			want:  map[string]string{"TimeTemplate": tem2, "Date": "2020/06/01  11:43:50.020", "Msg": "log message for <Test2>, have fun too!!"},
			want1: true,
		},
		{
			name: "filter other message with test 1",
			fields: fields{
				reg:          reg1,
				timeTemplate: tem1,
			},
			args:  args{msg: msgOther},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Filter{
				reg:          tt.fields.reg,
				timeTemplate: tt.fields.timeTemplate,
			}
			got, got1 := f.Filter(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Filter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFilters_Filter(t *testing.T) {
	tem1 := "[2006-01-02 15:04:05,000]"
	reg1, err := regexp.Compile(`^\s*(\[\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2},\d{3}\])\s*(.*<Test1>.*?)\s*$`)
	if err != nil {
		t.Fatal(err)
	}

	msg1 := " [2021-06-01 10:43:50,020] log message for <Test1>, have fun!!     "

	tem2 := "2006/01/02 15:04:05.000"
	reg2, err := regexp.Compile(`^\s*(\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2}.\d{3})\s*(.*<Test2>.*?)\s*$`)
	if err != nil {
		t.Fatal(err)
	}

	msg2 := " 2020/06/01  11:43:50.020    log message for <Test2>, have fun too!!"
	msgOther := "other message for <Test1> <Test2>, have fun..."

	filters := []Filter{{
		reg:          reg1,
		timeTemplate: tem1,
	}, {
		reg:          reg2,
		timeTemplate: tem2,
	}}
	type fields struct {
		filters []Filter
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
		want1  bool
	}{
		{
			name:   "normal test 1",
			fields: fields{filters: filters},
			args:   args{msg: msg1},
			want:   map[string]string{"TimeTemplate": tem1, "Date": "[2021-06-01 10:43:50,020]", "Msg": "log message for <Test1>, have fun!!"},
			want1:  true,
		},
		{
			name:   "normal test 2",
			fields: fields{filters: filters},
			args:   args{msg: msg2},
			want:   map[string]string{"TimeTemplate": tem2, "Date": "2020/06/01  11:43:50.020", "Msg": "log message for <Test2>, have fun too!!"},
			want1:  true,
		},
		{
			name:   "filter other message",
			fields: fields{filters: filters},
			args:   args{msg: msgOther},
			want:   nil,
			want1:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := FiltersImpl{
				filters: tt.fields.filters,
			}
			got, got1 := fs.Filter(tt.args.msg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Filter() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewFilters(t *testing.T) {
	tem1 := "[2006-01-02 15:04:05,000]"
	reg1, err := regexp.Compile(`^\s*(\[\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2},\d{3}\])\s*(.*<Test1>.*?)\s*$`)
	if err != nil {
		t.Fatal(err)
	}

	tem2 := "2006/01/02 15:04:05.000"
	reg2, err := regexp.Compile(`^\s*(\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2}.\d{3})\s*(.*<Test2>.*?)\s*$`)
	if err != nil {
		t.Fatal(err)
	}

	filter1 := Filter{
		reg:          reg1,
		timeTemplate: tem1,
	}

	filter2 := Filter{
		reg:          reg2,
		timeTemplate: tem2,
	}

	filterInfo1 := config.Filter{
		Keyword:      "<Test1>",
		DateReg:      `\[\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2},\d{3}\]`,
		TimeTemplate: tem1,
	}

	filterInfo2 := config.Filter{
		Keyword:      "<Test2>",
		DateReg:      `\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2}.\d{3}`,
		TimeTemplate: tem2,
	}

	type args struct {
		filtersInfo []config.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    Filters
		wantErr bool
	}{
		{
			name: "two infos",
			args: args{filtersInfo: []config.Filter{filterInfo1, filterInfo2}},
			want: &FiltersImpl{filters: []Filter{filter1, filter2}},
		},
		{
			name: "one info",
			args: args{filtersInfo: []config.Filter{filterInfo2}},
			want: &FiltersImpl{filters: []Filter{filter2}},
		},
		{
			name:    "nil info",
			args:    args{filtersInfo: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "null info",
			args:    args{filtersInfo: []config.Filter{}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFilters(tt.args.filtersInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilters() got = %v, want %v", got, tt.want)
			}
		})
	}
}
