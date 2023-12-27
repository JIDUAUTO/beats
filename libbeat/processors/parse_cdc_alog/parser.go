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

package parse_cdc_alog

import (
	"strconv"
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
	procName = "parse_cdc_alog"
	logName  = "processor." + procName
)

const (
	UsbMounted   = "UsbDeviceService:     state=MOUNTED"
	UsbUnmounted = "UsbDeviceService:     state=EJECTING"
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

	duration, err := time.ParseDuration(config.AllowOld)
	if err != nil {
		return nil, makeErrCompute(errors.New("duration configuration error"))
	}
	config.AllowOldDuration = duration

	p := &parseServerlog{
		config: config,
		logger: logp.NewLogger(logName),
	}

	return p, nil
}

// Run parse log
func (p *parseServerlog) Run(event *beat.Event) (*beat.Event, error) {
	// event filter
	processor, err := event.GetValue(processors.FieldProcessor)
	if err != nil {
		return nil, err
	}
	if processor != procName {
		return event, nil
	}

	collector, err := event.GetValue(processors.FieldCollector)
	if err != nil {
		return nil, err
	}
	err = processors.LogPreprocessing(event, processors.LogFormat(collector.(string)))
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

	const dateLayout = "01-02 15:04:05.999999"
	msg := message.(string)
	if len(msg) <= len(dateLayout) {
		return nil, nil
	}
	// 日志时间
	logtime, err := time.ParseInLocation(dateLayout, msg[:len(dateLayout)], time.Local)
	if err != nil {
		return nil, makeErrCompute(errors.New("invalid log time: " + msg[:len(dateLayout)]))
	}

	// 解析文件名称中的信息
	path, err := event.GetValue(processors.LogFilename)
	if err != nil {
		return nil, makeErrMissingField("log.file.path", err)
	}

	/* parse */
	items := strings.Split(path.(string), "@")
	if len(items) < 6 {
		return nil, nil
	}
	event.Fields["filename"] = items[0][:strings.LastIndex(items[0], ".")]
	event.Fields["ecu"] = items[1]
	event.Fields["vid"] = items[2]
	event.Fields["log_type"] = items[3]
	event.Fields["modified_at"] = items[4]
	event.Fields["uploaded_at"] = items[5]

	// 移除file信息
	delete(event.Fields, processors.LogFilename)

	// 解析日期
	lastModifiedAt, err := strconv.Atoi(items[4])
	if err != nil {
		return nil, makeErrCompute(errors.New("invalid file modify time"))
	}
	mt := time.UnixMilli(int64(lastModifiedAt)).In(time.Local)
	logtime = logtime.AddDate(mt.Year(), 0, 0)
	if time.Now().Sub(logtime) > p.config.AllowOldDuration || logtime.Sub(time.Now()) > p.config.AllowOldDuration {
		// 过滤日期差异很大的数据
		return nil, nil
	}
	event.Fields[p.config.TimeField] = logtime

	delete(event.Fields, "fields")

	// 业务属性
	// if strings.Index(msg, UsbMounted) > 0 {
	// 	event.Fields["usb_mounted"] = 1
	// } else if strings.Index(msg, UsbUnmounted) > 0 {
	// 	event.Fields["usb_mounted"] = 2
	// }

	return event, nil
}

func (p *parseServerlog) String() string {
	conf, _ := json.Marshal(p.config)
	return procName + "=" + string(conf)
}
