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

package parse_vehicle_tracelog

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/processors"
	"github.com/elastic/beats/v7/libbeat/processors/util"
)

const (
	procName   = "parse_vehicle_tracelog"
	logName    = "processor." + procName
	patternStr = "^(\\d{4}\\-\\d{2}\\-\\d{2}\\s\\d{2}:\\d{2}:\\d{2}\\.\\d{3})\\s+(\\d+)\\s+(\\d+)\\s+([a-zA-Z]+)\\s+(.*):\\s*##MSG##\\s*\\[(\\w*)\\]\\s*\\[(\\w*)\\]\\s*\\[(\\w*)\\]\\s*\\[([^\\[\\]]*)\\]\\s*\\[([^\\[\\]]*)\\]\\s+"
)

func init() {
	processors.RegisterPlugin(procName, NewParseVehicleTracelog)
	// jsprocessor.RegisterPlugin(strings.Title(procName), New)
}

type parseVehicleTrace2trace struct {
	config  Config
	logger  *logp.Logger
	pattern *regexp.Regexp
}

// NewParseVehicleTrace2trace constructs a new parse_vehicle_trace2trace processor.
func NewParseVehicleTracelog(cfg *common.Config) (processors.Processor, error) {
	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		return nil, makeErrConfigUnpack(err)
	}

	logger := logp.NewLogger(logName)

	p := &parseVehicleTrace2trace{
		config:  config,
		logger:  logger,
		pattern: regexp.MustCompile(patternStr),
	}

	return p, nil
}

// Run processing parser
func (p *parseVehicleTrace2trace) Run(event *beat.Event) (*beat.Event, error) {
	// get the content of log
	message, err := event.GetValue(p.config.Field)
	if err != nil {
		if p.config.IgnoreMissing && errors.Cause(err) == common.ErrKeyNotFound {
			return event, nil
		}
		return nil, makeErrMissingField(p.config.Field, err)
	}

	path, err := event.GetValue("log.file.path")
	if err != nil {
		return nil, makeErrMissingField("log.file.path", err)
	}

	/* parse */
	items := strings.Split(path.(string), "@")
	if len(items) == 6 {
		event.Fields["x-header_filename"] = items[0][strings.LastIndex(items[0], "/")+1 : strings.LastIndex(items[0], ".")]
		event.Fields["x-header_ecu"] = items[1]
		event.Fields["x-header_vid"] = items[2]
		event.Fields["x-header_log_type"] = items[3]
		event.Fields["x-header_created_at"] = items[4]
		event.Fields["x-header_uploaded_at"] = items[5]
	}

	msg := message.(string)
	// override
	event.Fields["message"] = msg
	matches := p.pattern.FindStringSubmatch(msg)
	if len(matches) < 11 || len(matches[6]) < 0 {
		// Drop event
		return nil, nil
	}

	// the time field is served for trace collector
	event.Fields["time"] = matches[1]

	pid, _ := strconv.ParseInt(matches[2], 10, 64)
	event.Fields["pid"] = pid
	tid, _ := strconv.ParseInt(matches[3], 10, 64)
	event.Fields["tid"] = tid
	if value, ok := util.LevelMap[matches[4]]; ok {
		event.Fields["level"] = value
	} else {
		event.Fields["level"] = strings.ToUpper(matches[4])
	}
	event.Fields["tag"] = matches[5]
	event.Fields["trace_id"] = matches[6]

	event.Fields["span_id"] = matches[7]
	event.Fields["parent_span_id"] = matches[8]
	event.Fields["network"] = matches[9]
	event.Fields["user_id"] = matches[10]
	if endIdx := strings.LastIndex(msg, "##MSG##"); endIdx > len(matches[0]) {
		event.Fields["message"] = msg[len(matches[0]):endIdx]
	} else {
		event.Fields["message"] = msg[len(matches[0]):]
	}

	return event, nil
}

func (p *parseVehicleTrace2trace) String() string {
	conf, _ := json.Marshal(p.config)
	return procName + "=" + string(conf)
}
