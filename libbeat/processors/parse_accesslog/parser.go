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

package parse_accesslog

import (
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/processors"
	"github.com/elastic/beats/v7/libbeat/processors/util"
)

const (
	procName = "parse_accesslog"
	logName  = "processor." + procName
	spliter  = ` `
)

func init() {
	processors.RegisterPlugin(procName, NewParseAccesslog)
	// jsprocessor.RegisterPlugin(strings.Title(procName), New)
}

type parseAccesslog struct {
	config Config
	logger *logp.Logger
}

// NewParseAccesslog constructs a new parse_accesslog processor.
func NewParseAccesslog(cfg *common.Config) (processors.Processor, error) {
	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		return nil, makeErrConfigUnpack(err)
	}

	logger := logp.NewLogger(logName)

	p := &parseAccesslog{
		config: config,
		logger: logger,
	}

	return p, nil
}

// Run parse log
func (p *parseAccesslog) Run(event *beat.Event) (*beat.Event, error) {
	// ignore panic
	defer func() {
		if err := recover(); err != nil {
			p.logger.Warnf("accesslog parse panic: %v", err)
			return
		}
	}()
	/* filter */
	// event filter
	processor, err := event.GetValue(processors.FieldProcessor)
	if err != nil {
		return event, nil
	}
	if processor != procName {
		return event, nil
	}

	collector, err := event.GetValue(processors.FieldCollector)
	if err != nil {
		return nil, err
	}
	err = processors.LogPreprocessing(event, processors.LogFormat(collector.(string)))
	// test
	// err := processors.LogPreprocessing(event, processors.LogFormat("ilogtail"))
	if err != nil {
		return event, err
	}

	msg := event.Fields["message"].(string)
	// log filter
	tagId := strings.Index(msg, util.MsgTag)
	if strings.Index(msg, "request=/misc/ping") != -1 && strings.Index(msg, "request=/actuator/health") != -1 && tagId != -1 {
		return nil, nil
	}

	/* parse */
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

	preInfo := msg[24:tagId]
	items := strings.Split(preInfo, " [")
	subPreItems := strings.Split(items[0], " ")
	event.Fields["jiduservicename"] = subPreItems[0]

	// filter benchmark log
	if strings.HasPrefix(util.Trim(items[3]), util.BenchmarkPrefix) {
		return nil, nil
	}

	level := ""
	if len(items) == 5 && len(subPreItems) == 3 {
		event.Fields["hostname"] = subPreItems[1]
		level = strings.ToUpper(subPreItems[2])
		event.Fields["level"] = level
		fileThread := strings.Split(items[1], "] ")
		if len(fileThread) == 2 {
			event.Fields["file"] = fileThread[0]
			tmId := strings.Index(fileThread[1], " ")
			event.Fields["thread"] = fileThread[1][0:tmId]
			event.Fields["method"] = fileThread[tmId+1:]
		}
		event.Fields["line"] = items[2][0 : len(items[2])-1]
		event.Fields["trace_id"] = items[3][0 : len(items[3])-1]
		event.Fields["span_id"] = items[4][0 : len(items[4])-2]
	}
	startId := strings.Index(msg, spliter)
	if startId == -1 {
		return event, nil
	}
	tailMsg := msg[startId+2:]
	kvs := strings.Split(tailMsg, spliter)
	for _, kv := range kvs {
		idx := strings.Index(kv, "=")
		event.Fields[kv[0:idx]] = kv[idx+1:]
	}
	if event.Fields["level"] != level {
		event.Fields["level"] = level
	}

	status, err := strconv.ParseInt(event.Fields["status"].(string), 10, 64)
	if err != nil {
		status = 0
	}
	event.Fields["status"] = status

	requestLength, err := strconv.ParseInt(event.Fields["request_length"].(string), 10, 64)
	if err != nil {
		requestLength = 0
	}
	event.Fields["request_length"] = requestLength

	latencyMs, err := strconv.ParseInt(event.Fields["latency-ms"].(string), 10, 64)
	if err != nil {
		latencyMs = 0
	}
	event.Fields["latency-ms"] = latencyMs

	msg = strings.ReplaceAll(msg, `\u001f`, "")

	event.Fields["message"] = msg

	return event, nil
}

func (p *parseAccesslog) String() string {
	conf, _ := json.Marshal(p.config)
	return procName + "=" + string(conf)
}
