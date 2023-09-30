// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package parse_serverlog

import (
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/processors"
	"github.com/elastic/beats/v7/libbeat/processors/util"
)

const (
	procName = "parse_serverlog"
	logName  = "processor." + procName
)

func init() {
	processors.RegisterPlugin(procName, New)
	// jsprocessor.RegisterPlugin(strings.Title(procName), New)
}

type parseServerlog struct {
	config Config
	logger *logp.Logger
}

// New constructs a new parse_serverlog processor.
func New(cfg *common.Config) (processors.Processor, error) {
	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		return nil, makeErrConfigUnpack(err)
	}

	logger := logp.NewLogger(logName)

	p := &parseServerlog{
		config: config,
		logger: logger,
	}

	return p, nil
}

// Run parse log
func (p *parseServerlog) Run(event *beat.Event) (*beat.Event, error) {
	// event filter
	processor, err := event.GetValue(processors.FieldProcessor)
	if err != nil {
		return event, nil
	}
	if processor != procName {
		return event, nil
	}

	err = processors.LogPreprocessing(event, processors.LogFormat(p.config.Collector))
	if err != nil {
		return event, err
	}

	message, err := event.GetValue(p.config.Field)
	if err != nil {
		if p.config.IgnoreMissing && errors.Cause(err) == common.ErrKeyNotFound {
			return event, nil
		}
		return nil, makeErrMissingField(p.config.Field, err)
	}
	msg := message.(string)
	event.Fields["message"] = msg

	// Parse time field
	for _, layout := range p.config.Layouts {
		ts, err := time.ParseInLocation(layout, msg[0:23], p.config.Timezone.Location())
		if err == nil {
			_, err = event.PutValue(p.config.TimeField, ts)
			if err != nil {
				return nil, makeErrCompute(err)
			}

			break
		}
	}

	items := strings.SplitN(msg, " ", 12)
	if len(items) < 12 {
		// Drop event<malformed log>
		return nil, nil
	}

	// filter benchmark log
	if strings.HasPrefix(util.Trim(items[9]), util.BenchmarkPrefix) {
		// Drop event<benchmark log>
		return nil, nil
	}

	event.Fields["jiduservicename"] = items[2]
	event.Fields["hostname"] = items[3]
	event.Fields["level"] = strings.ToUpper(items[4])

	var beginIdx, endIdx int
	line, err := strconv.ParseInt(util.Trim(items[8]), 10, 64)
	if err == nil {
		event.Fields["thread"] = util.Trim(items[5])
		event.Fields["class"] = items[6]
		event.Fields["method"] = items[7]
		event.Fields["line"] = line
		event.Fields["trace_id"] = util.Trim(items[9])
		event.Fields["span_id"] = util.Trim(items[10])

		if idx := strings.Index(msg, util.MsgTagConcatenated); idx > 0 {
			beginIdx = idx
			event.Fields["message"] = msg[idx+len(util.MsgTagConcatenated):]
		} else if idx = strings.Index(msg, util.MsgTag); idx > 0 {
			beginIdx = idx
			event.Fields["message"] = msg[idx+len(util.MsgTag):]
		} else {
			event.Fields["message"] = items[11]
		}
	}

	// 含有json数据
	endIdx = strings.LastIndex(msg, util.MsgTag)
	if beginIdx > 0 && beginIdx+len(util.MsgTag) != endIdx {
		var obj map[string]interface{}
		err = sonic.UnmarshalString(msg[beginIdx+len(util.MsgTag):endIdx], &obj)
		if err != nil {
			event.Fields["json_error"] = err
		} else {
			for k, v := range obj {
				event.Fields[k] = v
			}
		}
	}

	return event, nil
}

func (p *parseServerlog) String() string {
	conf, _ := json.Marshal(p.config)
	return procName + "=" + string(conf)
}
