package filter

import (
	"errors"
	"github.com/ClessLi/dingtalk-log-pusher/pkg/logbeat/config"
	"github.com/elastic/beats/libbeat/logp"
	"regexp"
)

type Filters interface {
	Filter(msg string) (map[string]string, bool)
}

type FiltersImpl struct {
	filters []Filter
}

func NewFilters(filtersInfo []config.Filter) (Filters, error) {
	filters := make([]Filter, 0)
	for _, filterInfo := range filtersInfo {
		reg, err := regexp.Compile(`^\s*(` + filterInfo.DateReg + `)\s*(.*` + filterInfo.Keyword + `.*?)\s*$`)
		if err != nil {
			logp.Warn("Failed to initialize filter, cased by: %v", err)
			continue
		}

		filter := Filter{
			reg:          reg,
			timeTemplate: filterInfo.TimeTemplate,
		}

		filters = append(filters, filter)

	}

	if len(filters) < 1 {
		return nil, errors.New("Failed to initialize filters")
	}

	return &FiltersImpl{filters: filters}, nil
}

func (fs FiltersImpl) Filter(msg string) (map[string]string, bool) {
	for _, filter := range fs.filters {
		if result, isMatched := filter.Filter(msg); isMatched {
			return result, isMatched
		}
	}
	return nil, false
}

type Filter struct {
	reg          *regexp.Regexp
	timeTemplate string
}

func (f Filter) Filter(msg string) (map[string]string, bool) {
	// TODO: 筛选出所需匹配字段外，加入日期格式化样板“key: TimeTemplate”，
	matches := f.reg.FindStringSubmatch(msg)
	if len(matches) == 3 {
		result := make(map[string]string)
		result["TimeTemplate"] = f.timeTemplate
		result["Date"] = matches[1]
		result["Msg"] = matches[2]
		return result, true
	}
	return nil, false
}
