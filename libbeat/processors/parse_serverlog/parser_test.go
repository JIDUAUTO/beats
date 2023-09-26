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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
)

func TestServerLogWithData(t *testing.T) {
	input := common.MapStr{
		"message": `{"contents":{"content":"2023-09-18 11:32:58.511 ai-repair-common ai-repair-common-69685c846c-kr47m INFO [http-nio-8080-exec-1]
com.jidu.postsale.config.LogAspect doAround [66] [4652dc92fb8240777ad468f1623aaaff] [f9567a128ed25419]
【智能维修】【响应日志】{\"code\":0,\"msg\":\"请求成功\"}##JIDU##{\"conts\":{\"cont\":\"123\"},\"ta\":{\"ip\":\"10.90.33.11\",\"name\":\"10.90.33.11\"},
\"time-test\":1695007978}##JIDU## time=2023-09-22T16:42:02+08:00 level=info msg=Stats In One Minute. AVGPT=0 SUM=0 TPS=0.00 statsKey=topic_passport_c_token@apisix-token-clean-passportc statsName=PULL_RT"},
"tags":{"container.image.name":"docker.jidudev.com/tech/ai-repair-common:s.95f66.57.1904","container.ip":"10.90.44.137",
"container.name":"ai-repair-common","host.ip":"10.90.162.80","host.name":"log-collector-6s7vk","k8s.namespace.name":"develop","k8s.node.ip":"10.90.33.11",
"k8s.node.name":"10.90.33.11","k8s.pod.name":"ai-repair-common-69685c846c-kr47m","k8s.pod.uid":"fc75c40f-f5b1-4e64-8ef9-0557c7ceca82",
"log.file.path":"/app/logs/ai-repair-common/serverlog.ai-repair-common-69685c846c-kr47m.log"},"time":1695007978}`,
	}
	testConfig, _ := common.NewConfigFrom(map[string]interface{}{
		"Field":           "message",
		"TimeField":       "logtime",
		"IgnoreMissing":   true,
		"IgnoreMalformed": true,
		"DropOrigin":      true,
	})
	actual := getActualValue(t, testConfig, input)
	expected := map[string]interface{}{
		"logtime":         "2023-09-18 11:32:58.511",
		"jiduservicename": "ai-repair-common",
		"hostname":        "ai-repair-common-69685c846c-kr47m",
		"level":           "INFO",
		"thread":          "http-nio-8080-exec-1",
		"class":           "com.jidu.postsale.config.LogAspect",
		"method":          "doAround",
		"line":            int64(66),
		"trace_id":        "4652dc92fb8240777ad468f1623aaaff",
		"span_id":         "f9567a128ed25419",
		"message":         `{"conts":{"cont":"123"},"ta":{"ip":"10.90.33.11","name":"10.90.33.11"},"time-test":1695007978}##JIDU## time=2023-09-22T16:42:02+08:00 level=info msg=Stats In One Minute. AVGPT=0 SUM=0 TPS=0.00 statsKey=topic_passport_c_token@apisix-token-clean-passportc statsName=PULL_RT`,
		"conts.cont":      "123",
		"ta.name":         "10.90.33.11",
		"ta.ip":           "10.90.33.11",
		"time-test":       float64(1695007978),
	}

	assert.Equal(t, expected["logtime"], actual["logtime"])
	assert.Equal(t, expected["jiduservicename"], actual["jiduservicename"])
	assert.Equal(t, expected["hostname"], actual["hostname"])
	assert.Equal(t, expected["level"], actual["level"])
	assert.Equal(t, expected["thread"], actual["thread"])
	assert.Equal(t, expected["class"], actual["class"])
	assert.Equal(t, expected["method"], actual["method"])
	assert.Equal(t, expected["line"], actual["line"])
	assert.Equal(t, expected["trace_id"], actual["trace_id"])
	assert.Equal(t, expected["span_id"], actual["span_id"])
	assert.Equal(t, expected["message"], actual["message"])

	contsCont, err := actual.GetValue("conts.cont")
	assert.Nil(t, err)
	assert.Equal(t, expected["conts.cont"], contsCont)

	taName, err := actual.GetValue("ta.name")
	assert.Nil(t, err)
	assert.Equal(t, expected["ta.name"], taName)

	taIp, err := actual.GetValue("ta.ip")
	assert.Nil(t, err)
	assert.Equal(t, expected["ta.ip"], taIp)

	assert.Equal(t, expected["time-test"], actual["time-test"])
}

func TestServerLogNoData(t *testing.T) {
	message := `{"contents":{"content":"2023-09-22 16:45:15.806 apisix-token-clean apisix-token-clean-5dcd4464b8-qnqzn INFO [/go/pkg/mod/jidudev.com/tech/rocketmq-client-go/v2@v2.1.6/internal/client.go] - github.com/apache/rocketmq-client-go/v2/internal.GetOrNewRocketMQClient.func3 [266] [] [] ##JIDU####JIDU## time=2023-09-22T16:45:15+08:00 level=info msg=receive get consumer running info request..."},"tags":{"container.image.name":"docker.jidudev.com/tech/apisix-token-clean:71df020e","container.ip":"10.80.224.116","container.name":"apisix-token-clean","host.ip":"10.80.225.102","host.name":"log-collector-2h277","k8s.namespace.name":"develop","k8s.node.ip":"10.80.11.20","k8s.node.name":"10.80.11.20","k8s.pod.name":"apisix-token-clean-5dcd4464b8-qnqzn","k8s.pod.uid":"de60fb4e-5536-4d4b-be56-9cff739f7d7a","log.file.path":"/app/logs/apisix-token-clean/serverlog.apisix-token-clean-5dcd4464b8-qnqzn.log"},"time":1695372315}`
}

func TestServerLogNoTag(t *testing.T) {
	message := `2023-09-22 16:20:46.531 fota-gateway 10.90.42.20 WARN [grpc-default-executor-1598] c.j.j.c.watch.CustomServicesWatch
getInstanceByHost [61] [10676155e01c165f732129fd9c2fa341] [c88dac5f5f59d7b6] jns_warn: Cannot query instance when host is empty,
host null, instances [DefaultServiceInstance{instanceId='develop/fota-gateway-665cbdf9bf-zxjfr', serviceId='develop/fota-gateway-665cbdf9bf-zxjfr',
host='10.90.42.20', port=8080, secure=false, metadata={color=default, control-by=develop/fota-gateway, grpc_port=9090, idc=bj, weight=100}}]`
}

func getActualValue(t *testing.T, config *common.Config, input common.MapStr) common.MapStr {
	p, err := NewParseServerlog(config)
	if err != nil {
		t.Fatal(err)
	}

	actual, _ := p.Run(&beat.Event{Fields: input})
	return actual.Fields
}
