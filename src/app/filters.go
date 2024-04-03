package app

import (
	"encoding/json"
	"regexp"
	"strings"
)

type LogFilter interface {
	// filter 过滤消息, 处理消息字段内容.
	filter(message []*LogMessage) ([]*LogMessage, error)
}

type RegularFilter struct {
	regular *regexp.Regexp
	fields  []string
}

func (f *RegularFilter) filter(message []*LogMessage) ([]*LogMessage, error) {
	for _, item := range message {
		text, dataOk := item.Attribute["log"]
		if !dataOk {
			continue
		}
		matches := f.regular.FindStringSubmatch(text.(string))
		if matches == nil {
			continue
		}
		for i, key := range f.fields {
			if len(matches) > i+1 {
				item.Attribute[key] = strings.Trim(matches[i+1], "\n")
			} else {
				continue
			}
		}
	}
	return message, nil
}

type JsonFilter struct{}

func (f *JsonFilter) filter(message []*LogMessage) ([]*LogMessage, error) {
	for _, item := range message {
		var attribute map[string]interface{}
		err := json.Unmarshal(item.Body, &attribute)
		if err != nil {
			return nil, err
		}
		item.Attribute = attribute
	}
	return message, nil
}
