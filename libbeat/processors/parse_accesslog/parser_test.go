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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
)

func TestServerLogWithData(t *testing.T) {
	input := common.MapStr{
		// "message": `2023-09-20 15:46:18.052 jidu-mesh-be jidu-mesh-be-74bcf7cdd6-z7z5b INFO [/home/jenkins/agent/workspace/all-project-golang-build-job/internal/transport/http/middleware/logger.go] - jidudev.com/tech/jidu-mesh-be/internal/transport/http/middleware.Logger.func1 [82] [3bbebdfd1ed58ae86c386a27acd2ebe4] [9de7d5accf3efb4d] ##JIDU####JIDU##\u001f time=2023-09-20T15:46:18+08:00\u001f level=info\u001f content-type=""\u001f http_referer=https://localhost:8001/mesh-fe/serveDetail/serviceMonitor?env=dev&jns=cn.pe.vi.remote-monitor-wti.mixed&cluster=dev&var-apis=all&var-pods-api=10.80.232.65&relative_api=1&from_api=now-5m&to_api=now&refresh_api=0\u001f http_user_agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36\u001f latency-ms=27\u001f method=GET\u001f remote_addr=172.16.22.54\u001f request=/v1/prom/pod/cpu/limit\u001f request_body=""\u001f request_length=0\u001f request_query=pod_name[]=remote-monitor-wti-55f757f8-rx5tv&start=2023-09-20T07:41:00.000Z&end=2023-09-20T07:46:00.000Z&step=30\u001f response_body={"code":0,"data":[{"namespace":"develop","pod":"remote-monitor-wti-55f757f8-rx5tv","container":"remote-monitor-wti","values":[[1695195660000,"2"],[1695195690000,"2"],[1695195720000,"2"],[1695195750000,"2"],[1695195780000,"2"],[1695195810000,"2"],[1695195840000,"2"],[1695195870000,"2"],[1695195900000,"2"],[1695195930000,"2"],[1695195960000,"2"]]}],"msg":"ok"}\u001f status=200\u001f transport=http`,
		"message": `{"contents":{"content":"2023-10-09 19:35:31.411 jidu-mesh-be jidu-mesh-be-74bcf7cdd6-hxdt9 INFO [/home/jenkins/agent/workspace/all-project-golang-build-job/internal/transport/http/middleware/logger.go] - jidudev.com/tech/jidu-mesh-be/internal/transport/http/middleware.Logger.func1 [82] [e9b51911b97ce763a871da84907aa0ec] [dec73c11670be967] ##JIDU####JIDU##\u001f time=2023-10-09T19:35:31+08:00\u001f level=info\u001f content-type=\"\"\u001f http_referer=https://mesh.jidudev.com/mesh-fe/serveDetail/serviceMonitor?env=dev&jns=cn.pe.vi.remote-monitor-wti.mixed&cluster=dev\u001f http_user_agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36\u001f latency-ms=29\u001f method=GET\u001f remote_addr=10.55.2.93\u001f request=/v1/prom/api/blocked_qps\u001f request_body=\"\"\u001f request_length=0\u001f request_query=jns_name=cn.pe.vi.remote-monitor-wti.mixed&start=2023-10-09T11:30:00.000Z&end=2023-10-09T11:35:00.000Z&step=20&uri[]=all&pod_name[]=10.80.249.141\u001f response_body={\"code\":0,\"data\":[],\"msg\":\"ok\"}\u001f status=200\u001f transport=http"},"tags":{"container.image.name":"docker.jidudev.com/tech/jidu-mesh-be:d.95b2c.c0.1214","container.ip":"10.80.225.69","container.name":"jidu-mesh-be","host.ip":"10.80.249.106","host.name":"log-collector-6zjvm","k8s.namespace.name":"develop","k8s.node.ip":"10.80.11.26","k8s.node.name":"10.80.11.26","k8s.pod.name":"jidu-mesh-be-74bcf7cdd6-hxdt9","k8s.pod.uid":"4bf669b7-d597-45c2-910e-1bf970b25c22","log.file.path":"/app/logs/jidu-mesh-be/serveraccesslog.jidu-mesh-be-74bcf7cdd6-hxdt9.log"},"time":1696851331}`,
	}
	testConfig, _ := common.NewConfigFrom(map[string]interface{}{
		"Field":         "message",
		"TimeField":     "@timestamp",
		"IgnoreMissing": true,
		//	定义Layouts
	})
	actual := getActualValue(t, testConfig, input)
	expected := map[string]interface{}{
		"@timestamp":      "2023-10-09 19:35:31.411",
		"content-type":    "\"\"",
		"file":            "/home/jenkins/agent/workspace/all-project-golang-build-job/internal/transport/http/middleware/logger.go",
		"hostname":        "jidu-mesh-be-74bcf7cdd6-hxdt9",
		"http_referer":    "https://mesh.jidudev.com/mesh-fe/serveDetail/serviceMonitor?env=dev&jns=cn.pe.vi.remote-monitor-wti.mixed&cluster=dev",
		"http_user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
		"jiduservicename": "jidu-mesh-be",
		"latency-ms":      int64(29),
		"level":           "INFO",
		"line":            "82",
		"message":         "2023-10-09 19:35:31.411 jidu-mesh-be jidu-mesh-be-74bcf7cdd6-hxdt9 INFO [/home/jenkins/agent/workspace/all-project-golang-build-job/internal/transport/http/middleware/logger.go] - jidudev.com/tech/jidu-mesh-be/internal/transport/http/middleware.Logger.func1 [82] [e9b51911b97ce763a871da84907aa0ec] [dec73c11670be967] ##JIDU####JIDU##\x1f time=2023-10-09T19:35:31+08:00\x1f level=info\x1f content-type=\"\"\x1f http_referer=https://mesh.jidudev.com/mesh-fe/serveDetail/serviceMonitor?env=dev&jns=cn.pe.vi.remote-monitor-wti.mixed&cluster=dev\x1f http_user_agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36\x1f latency-ms=29\x1f method=GET\x1f remote_addr=10.55.2.93\x1f request=/v1/prom/api/blocked_qps\x1f request_body=\"\"\x1f request_length=0\x1f request_query=jns_name=cn.pe.vi.remote-monitor-wti.mixed&start=2023-10-09T11:30:00.000Z&end=2023-10-09T11:35:00.000Z&step=20&uri[]=all&pod_name[]=10.80.249.141\x1f response_body={\"code\":0,\"data\":[],\"msg\":\"ok\"}\x1f status=200\x1f transport=http",
		"method":          "GET",
		"remote_addr":     "10.55.2.93",
		"request":         "/v1/prom/api/blocked_qps",
		"request_body":    "\"\"",
		"request_length":  int64(0),
		"request_query":   "jns_name=cn.pe.vi.remote-monitor-wti.mixed&start=2023-10-09T11:30:00.000Z&end=2023-10-09T11:35:00.000Z&step=20&uri[]=all&pod_name[]=10.80.249.141",
		"response_body":   "{\"code\":0,\"data\":[],\"msg\":\"ok\"}",
		"span_id":         "dec73c11670be967",
		"status":          int64(200),
		"thread":          "-",
		"time":            "2023-10-09T19:35:31+08:00",
		"trace_id":        "e9b51911b97ce763a871da84907aa0ec",
		"transport":       "http",
	}

	// assert.Equal(t, expected["@timestamp"], actual["@timestamp"])
	assert.Equal(t, expected["content-type"], actual["content-type"])
	assert.Equal(t, expected["file"], actual["file"])
	assert.Equal(t, expected["hostname"], actual["hostname"])
	assert.Equal(t, expected["http_referer"], actual["http_referer"])
	assert.Equal(t, expected["http_user_agent"], actual["http_user_agent"])
	assert.Equal(t, expected["jiduservicename"], actual["jiduservicename"])
	assert.Equal(t, expected["latency-ms"], actual["latency-ms"])
	assert.Equal(t, expected["level"], actual["level"])
	assert.Equal(t, expected["line"], actual["line"])
	assert.Equal(t, expected["message"], actual["message"])
	assert.Equal(t, expected["method"], actual["method"])
	assert.Equal(t, expected["remote_addr"], actual["remote_addr"])
	assert.Equal(t, expected["request"], actual["request"])
	assert.Equal(t, expected["request_body"], actual["request_body"])
	assert.Equal(t, expected["request_length"], actual["request_length"])
	assert.Equal(t, expected["request_query"], actual["request_query"])
	assert.Equal(t, expected["response_body"], actual["response_body"])
	assert.Equal(t, expected["span_id"], actual["span_id"])
	assert.Equal(t, expected["status"], actual["status"])
	assert.Equal(t, expected["thread"], actual["thread"])
	assert.Equal(t, expected["time"], actual["time"])
	assert.Equal(t, expected["trace_id"], actual["trace_id"])
	assert.Equal(t, expected["transport"], actual["transport"])

}

func getActualValue(t *testing.T, config *common.Config, input common.MapStr) common.MapStr {
	p, err := NewParseAccesslog(config)
	if err != nil {
		t.Fatal(err)
	}

	actual, _ := p.Run(&beat.Event{Fields: input})
	// assert.Equal(t, "2023-10-09 19:35:31.411", actual.Timestamp.Format("2006-01-02 15:04:05.000"))
	return actual.Fields
}
