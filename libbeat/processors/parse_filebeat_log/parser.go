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

package parse_filebeat_log

import (
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/processors"
)

const (
	procName = "parse_filebeat_log"
	logName  = "processor." + procName
)

func init() {
	processors.RegisterPlugin(procName, New)
	// jsprocessor.RegisterPlugin(strings.Title(procName), New)
}

type parseFilebeatLog struct {
	config Config
	logger *logp.Logger
}

// New constructs a new fingerprint processor.
func New(cfg *common.Config) (processors.Processor, error) {
	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		return nil, makeErrConfigUnpack(err)
	}

	log := logp.NewLogger(logName)

	p := &parseFilebeatLog{
		config: config,
		logger: log,
	}

	return p, nil
}

// Run parse filebeat's log
func (p *parseFilebeatLog) Run(event *beat.Event) (*beat.Event, error) {
	// event filter
	processor := event.Fields[processors.FieldProcessor]
	if processor != procName {
		return event, nil
	}

	// get the content of log
	msg, err := event.GetValue(p.config.Field)
	if err != nil {
		if p.config.IgnoreMissing && errors.Cause(err) == common.ErrKeyNotFound {
			return event, nil
		}

		return nil, makeErrMissingField(p.config.Field, err)
	}

	message, ok := msg.(string)
	if !ok {
		return nil, makeErrFieldType(p.config.Field, "string", fmt.Sprintf("%T", msg))
	}

	// Parse log message
	terms := strings.SplitN(message, "\t", 4)
	// Drop logs with incorrect format
	if len(terms) != 4 {
		if p.config.IgnoreMalformed {
			return event, nil
		}

		return nil, makeErrLogFormat("[datetime]\t[LEVEL]\t[hostname]\t[message]")
	}

	// Parse log time
	for _, layout := range p.config.Layouts {
		ts, err := time.ParseInLocation(layout, terms[0], p.config.Timezone.Location())
		if err == nil {
			_, err = event.PutValue(p.config.TimeField, ts.UTC())
			if err != nil {
				return nil, makeErrCompute(err)
			}

			break
		}
	}

	// Padding fields
	event.Fields["level"] = strings.ToUpper(terms[1])
	event.Fields["hostname"] = terms[2]
	event.Fields["message"] = terms[3]

	return event, nil
}

func (p *parseFilebeatLog) String() string {
	conf, _ := json.Marshal(p.config)
	return procName + "=" + string(conf)
}
