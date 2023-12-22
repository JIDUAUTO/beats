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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/processors"
)

func TestServerLogWithData(t *testing.T) {
	input := common.MapStr{
		"message": `{"@timestamp":"2023-09-27T10:55:53.798Z","@metadata":{"beat":"filebeat","type":"_doc","version":"7.9.3"},
"log":{"file":{"path":"/vlog/cdc/A_log_3292_20231221_195338.gz.1737806441831505920@cdc@6c9b10c6fd944651f6c8a22fa376ec13@logcat@1703159620000@1703160317000"},
"offset":6984921},"message":"12-21 20:34:38.005963  3810  6369 D BTS     : 67380602 [cockpit_perception_proxy.cc][34938]user callback elapsed_time_ms:0.header id:3500872286,sn:1874249","fields":{"servicetype":"syslogcdc"}}`,
		"fields": common.MapStr{
			"handler":   procName,
			"collector": string(processors.LogFormatFilebeat),
		},
	}

	testConfig, _ := common.NewConfigFrom(map[string]interface{}{
		"Field":         "message",
		"TimeField":     "@timestamp",
		"IgnoreMissing": false,
	})
	actual := getActualValue(t, testConfig, input)
	expected := map[string]interface{}{
		"filename":    "A_log_3292_20231221_195338.gz",
		"ecu":         "cdc",
		"vid":         "6c9b10c6fd944651f6c8a22fa376ec13",
		"log_type":    "logcat",
		"modified_at": "1703159620000",
		"uploaded_at": "1703160317000",
		"@timestamp":  time.Date(2023, 12, 21, 20, 34, 38, int(5963*time.Microsecond), time.Local),
		"message":     "12-21 20:34:38.005963  3810  6369 D BTS     : 67380602 [cockpit_perception_proxy.cc][34938]user callback elapsed_time_ms:0.header id:3500872286,sn:1874249",
	}

	assert.Equal(t, expected["filename"], actual["filename"])
	assert.Equal(t, expected["ecu"], actual["ecu"])
	assert.Equal(t, expected["vid"], actual["vid"])
	assert.Equal(t, expected["log_type"], actual["log_type"])
	assert.Equal(t, expected["modified_at"], actual["modified_at"])
	assert.Equal(t, expected["uploaded_at"], actual["uploaded_at"])
	assert.Equal(t, expected["@timestamp"], actual["@timestamp"])
	assert.Equal(t, expected["message"], actual["message"])
}

func getActualValue(t *testing.T, config *common.Config, input common.MapStr) common.MapStr {
	p, err := New(config)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := p.Run(&beat.Event{Fields: input})
	if err != nil {
		t.Fatal(err)
	}

	return actual.Fields
}
