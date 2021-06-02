package output_formatter

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"regexp"
	"time"
)

type OutputFormatter interface {
	Format(msgs ...map[string]string) ([]beat.Event, error)
}

type OutputFormatterImpl struct {
	formatStr string
	fields    []string
}

func (f OutputFormatterImpl) Format(msgs ...map[string]string) ([]beat.Event, error) {
	if msgs == nil || len(msgs) < 1 {
		return nil, errors.New("messages is null")
	}
	events := make([]beat.Event, 0)
	//defer func() {
	//	err := recover()
	//}()
	for _, msgMap := range msgs {
		values := f.getValues(msgMap)
		if len(values) != len(f.fields) {
			logp.Warn("Format message '%v' error, want fields length %d, got %d", msgMap, len(f.fields), len(values))
			continue
		}
		msg := fmt.Sprintf(f.formatStr, f.getValues(msgMap)...)

		event := beat.Event{
			Timestamp: f.generateTimestamp(msgMap),
			Fields:    common.MapStr{"msg": msg},
		}

		events = append(events, event)
	}
	if len(events) < 1 {
		return nil, fmt.Errorf("Failed to format message: %v", msgs)
	}
	return events, nil
}

func (f OutputFormatterImpl) generateTimestamp(msgMap map[string]string) time.Time {
	template, hasTemplate := msgMap["TimeTemplate"]
	date, hasDate := msgMap["Date"]
	if hasTemplate && hasDate {
		t, err := time.Parse(changeSpecialChar2Dot(template), changeSpecialChar2Dot(date))
		if err == nil {
			return t
		} else {
			logp.Warn("Failed to generate timestamp with 'Date': %v, and 'TimeTemplate': %v, Cased by: %v", date, template, err)
		}
	} else {
		logp.Info("There is no 'Date' or 'TimeTemplate' to resolve to a timestamp")
	}
	return time.Now()
}

func (f OutputFormatterImpl) getValues(msgMap map[string]string) []interface{} {
	values := make([]interface{}, 0)
	for i := 0; i < len(f.fields); i++ {
		if value, has := msgMap[f.fields[i]]; has {
			values = append(values, value)
		}
	}
	return values
}

func NewOutputFormatter(format string) (OutputFormatter, error) {
	var (
		valueStart, valueEnd bool
		formatStrBuff        = bytes.NewBufferString("")
		fieldBuff            = bytes.NewBufferString("")
		fields               = make([]string, 0)
	)

	n := len(format)

	for i := 0; i < n; i++ {

		if !valueStart {
			if i+1 != n && format[i] == '$' && format[i+1] == '{' {
				_, err := formatStrBuff.WriteString("%v")
				if err != nil {
					return nil, err
				}
				valueStart = true
				i++
				continue
			}
			err := formatStrBuff.WriteByte(format[i])
			if err != nil {
				return nil, err
			}
			continue
		}

		if !valueEnd {
			if format[i] == '}' {
				fields = append(fields, fieldBuff.String())
				fieldBuff.Reset()
				valueStart = false
				valueEnd = false
				continue
			}
			err := fieldBuff.WriteByte(format[i])
			if err != nil {
				return nil, err
			}
		}

	}

	formatStr := formatStrBuff.String()

	if len(fields) < 1 {
		return nil, fmt.Errorf("Can not parse format string: %v", format)
	}

	return &OutputFormatterImpl{
		formatStr: formatStr,
		fields:    fields,
	}, nil
}

func changeSpecialChar2Dot(str string) string {
	reg, _ := regexp.Compile(`[^0-9a-zA-Z]+`)
	ret := reg.ReplaceAllString(str, ".")
	return ret
}
