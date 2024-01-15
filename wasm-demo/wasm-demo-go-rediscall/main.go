// Copyright (c) 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/resp"

	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
)

func main() {
	wrapper.SetCtx(
		"redis-demo",
		wrapper.ParseConfigBy(parseConfig),
		wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
		wrapper.ProcessResponseHeadersBy(onHttpResponseHeaders),
	)
}

type RedisCallConfig struct {
	client wrapper.RedisClient
	qpm    int
}

func parseConfig(json gjson.Result, config *RedisCallConfig, log wrapper.Log) error {
	serviceSource := json.Get("serviceSource").String()
	serviceName := json.Get("serviceName").String()
	servicePort := json.Get("servicePort").Int()
	username := json.Get("username").String()
	password := json.Get("password").String()
	timeout := json.Get("timeout").Int()
	qpm := json.Get("qpm").Int()
	config.qpm = int(qpm)
	switch serviceSource {
	case "k8s":
		namespace := json.Get("namespace").String()
		config.client = wrapper.NewRedisClusterClient(wrapper.K8sCluster{
			ServiceName: serviceName,
			Namespace:   namespace,
			Port:        servicePort,
		})
	case "nacos":
		namespace := json.Get("namespace").String()
		config.client = wrapper.NewRedisClusterClient(wrapper.NacosCluster{
			ServiceName: serviceName,
			NamespaceID: namespace,
			Port:        servicePort,
		})
	case "ip":
		config.client = wrapper.NewRedisClusterClient(wrapper.StaticIpCluster{
			ServiceName: serviceName,
			Port:        servicePort,
		})
	case "dns":
		domain := json.Get("domain").String()
		config.client = wrapper.NewRedisClusterClient(wrapper.DnsCluster{
			ServiceName: serviceName,
			Port:        servicePort,
			Domain:      domain,
		})
	default:
		return errors.New("unknown service source: " + serviceSource)
	}
	config.client.Init(username, password, timeout)
	return nil
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config RedisCallConfig, log wrapper.Log) types.Action {
	now := time.Now()
	minuteAligned := now.Truncate(time.Minute)
	timeStamp := strconv.FormatInt(minuteAligned.Unix(), 10)
	config.client.Get(timeStamp, func(status int, response resp.Value) {
		if status != 0 {
			log.Errorf("Error occured while calling redis")
			defer proxywasm.ResumeHttpRequest()
			return
		}
		ctx.SetContext("timeStamp", timeStamp)
		ctx.SetContext("CallTimeLeft", strconv.Itoa(config.qpm-response.Integer()-1))
		if response.Integer() >= config.qpm {
			proxywasm.SendHttpResponse(429, [][2]string{{"timeStamp", timeStamp}, {"CallTimeLeft", "0"}}, []byte("Too many requests"), -1)
		}
		config.client.Incr(timeStamp, func(status int, response resp.Value) {
			defer proxywasm.ResumeHttpRequest()
			if status != 0 {
				log.Errorf("Error occured while calling redis")
			}
		})
	})

	return types.ActionPause
}

func onHttpResponseHeaders(ctx wrapper.HttpContext, config RedisCallConfig, log wrapper.Log) types.Action {
	proxywasm.AddHttpResponseHeader("timeStamp", ctx.GetContext("timeStamp").(string))
	proxywasm.AddHttpResponseHeader("CallTimeLeft", ctx.GetContext("CallTimeLeft").(string))
	return types.ActionContinue
}
