package output_formatter

import (
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"reflect"
	"testing"
	"time"
)

func TestNewOutputFormatter(t *testing.T) {
	defaultFormatStr := "${Env}环境-${Title}日志：\n\t日志时间：${Date}\n\t日志信息：${Msg}"
	wrongFormatStr := "$Env环境-$Title日志\n\ttimestamp: $Date\n\tmessage: $Msg"
	type args struct {
		format string
	}
	tests := []struct {
		name    string
		args    args
		msg     map[string]string
		want    OutputFormatter
		wantErr bool
	}{
		{
			name: "default test",
			args: args{format: defaultFormatStr},
			msg:  map[string]string{"Env": "UAT", "Title": "测试应用", "Date": "[2021-06-01 13:50:30,253]", "Msg": "log message for <Test1>, have fun!!", "TimeTemplate": "[2006-01-02 15:04:05,000]"},
			want: &OutputFormatterImpl{
				formatStr: "%v环境-%v日志：\n\t日志时间：%v\n\t日志信息：%v",
				fields:    []string{"Env", "Title", "Date", "Msg"},
			},
		},
		{
			name:    "wrong format test",
			args:    args{format: wrongFormatStr},
			msg:     map[string]string{"Env": "UAT", "Title": "测试应用", "Date": "[2021-06-01 13:50:30,253]", "Msg": "log message for <Test1>, have fun!!", "TimeTemplate": "[2006-01-02 15:04:05,000]"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOutputFormatter(tt.args.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOutputFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOutputFormatter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatter_Format(t *testing.T) {
	defaultFormatStr := "%v环境-%v日志：\n\t日志时间：%v\n\t日志信息：%v"
	defaultFields := []string{"Env", "Title", "Date", "Msg"}
	defaultMsg := map[string]string{"Env": "UAT", "Title": "测试应用", "Date": "[2021-06-01 13:50:30,253]", "Msg": "log message for <Test1>, have fun!!", "TimeTemplate": "[2006-01-02 15:04:05,000]"}
	defaultTimestamp, err := time.Parse(changeSpecialChar2Dot(defaultMsg["TimeTemplate"]), changeSpecialChar2Dot(defaultMsg["Date"]))
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		formatStr string
		fields    []string
	}
	type args struct {
		msgs []map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []beat.Event
		wantErr bool
	}{
		{
			name: "default test",
			fields: fields{
				formatStr: defaultFormatStr,
				fields:    defaultFields,
			},
			args: args{msgs: []map[string]string{defaultMsg}},
			want: []beat.Event{{
				Timestamp: defaultTimestamp,
				Fields:    common.MapStr{"msg": fmt.Sprintf(defaultFormatStr, defaultMsg["Env"], defaultMsg["Title"], defaultMsg["Date"], defaultMsg["Msg"])},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := OutputFormatterImpl{
				formatStr: tt.fields.formatStr,
				fields:    tt.fields.fields,
			}
			got, err := f.Format(tt.args.msgs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatter_generateTimestamp(t *testing.T) {
	defaultFormatStr := "%v环境-%v日志：\n\t日志时间：%v\n\t日志信息：%v"
	defaultFields := []string{"Env", "Title", "Date", "Msg"}
	defaultMsg := map[string]string{"Env": "UAT", "Title": "测试应用", "Date": "[2021-06-01 13:50:30,253]", "Msg": "log message for <Test1>, have fun!!", "TimeTemplate": "[2006-01-02 15:04:05,000]"}
	defaultTimestamp, err := time.Parse(changeSpecialChar2Dot(defaultMsg["TimeTemplate"]), changeSpecialChar2Dot(defaultMsg["Date"]))
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		formatStr string
		fields    []string
	}
	type args struct {
		msgMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{
			name: "default test",
			fields: fields{
				formatStr: defaultFormatStr,
				fields:    defaultFields,
			},
			args: args{msgMap: defaultMsg},
			want: defaultTimestamp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := OutputFormatterImpl{
				formatStr: tt.fields.formatStr,
				fields:    tt.fields.fields,
			}
			if got := f.generateTimestamp(tt.args.msgMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_changeSpecialChar2dot(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal test",
			args: args{"[2021-06-01 13:50:30, 253]"},
			want: ".2021.06.01.13.50.30.253.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := changeSpecialChar2Dot(tt.args.str); got != tt.want {
				t.Errorf("changeSpecialChar2Dot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatterImpl_getValues(t *testing.T) {
	defaultFormatStr := "%v环境-%v日志：\n\t日志时间：%v\n\t日志信息：%v"
	defaultFields := []string{"Env", "Title", "Date", "Msg"}
	defaultMsg := map[string]string{"Env": "UAT", "Title": "测试应用", "Date": "[2021-06-01 13:50:30,253]", "Msg": "log message for <Test1>, have fun!!", "TimeTemplate": "[2006-01-02 15:04:05,000]"}

	type fields struct {
		formatStr string
		fields    []string
	}
	type args struct {
		msgMap map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []interface{}
	}{
		{
			name: "default test",
			fields: fields{
				formatStr: defaultFormatStr,
				fields:    defaultFields,
			},
			args: args{msgMap: defaultMsg},
			want: []interface{}{"UAT", "测试应用", "[2021-06-01 13:50:30,253]", "log message for <Test1>, have fun!!"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := OutputFormatterImpl{
				formatStr: tt.fields.formatStr,
				fields:    tt.fields.fields,
			}
			if got := f.getValues(tt.args.msgMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
