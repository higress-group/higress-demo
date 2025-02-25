package main

import (
	"strconv"
	"time"

	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm"
	"github.com/higress-group/proxy-wasm-go-sdk/proxywasm/types"
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
	serviceName := json.Get("serviceName").String()
	servicePort := json.Get("servicePort").Int()
	domain := json.Get("domain").String()
	username := json.Get("username").String()
	password := json.Get("password").String()
	timeout := json.Get("timeout").Int()
	qpm := json.Get("qpm").Int()
	config.qpm = int(qpm)
	config.client = wrapper.NewRedisClusterClient(wrapper.DnsCluster{
		ServiceName: serviceName,
		Port:        servicePort,
		Domain:      domain,
	})
	return config.client.Init(username, password, timeout)
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config RedisCallConfig, log wrapper.Log) types.Action {
	now := time.Now()
	minuteAligned := now.Truncate(time.Minute)
	timeStamp := strconv.FormatInt(minuteAligned.Unix(), 10)
	// 如果 redis api 返回的 err != nil，一般是由于网关找不到 redis 后端服务，请检查是否误删除了 redis 后端服务
	err := config.client.Incr(timeStamp, func(response resp.Value) {
		if err := response.Error(); err != nil {
			log.Errorf("Error occured while calling redis")
			proxywasm.SendHttpResponse(430, nil, []byte("Error while calling redis"), -1)
		} else {
			ctx.SetContext("timeStamp", timeStamp)
			ctx.SetContext("CallTimeLeft", strconv.Itoa(config.qpm-response.Integer()))
			if response.Integer() == 1 {
				config.client.Expire(timeStamp, 60, func(response resp.Value) {
					proxywasm.ResumeHttpRequest()
				})
			} else {
				if response.Integer() > config.qpm {
					proxywasm.SendHttpResponse(429, [][2]string{{"timeStamp", timeStamp}, {"CallTimeLeft", "0"}}, []byte("Too many requests"), -1)
				} else {
					proxywasm.ResumeHttpRequest()
				}
			}
		}
	})
	if err != nil {
		// 由于调用redis失败，放行请求，记录日志
		log.Errorf("Error occured while calling redis, it seems cannot find the redis cluster.")
		return types.ActionContinue
	} else {
		// 请求hold住，等待redis调用完成
		return types.ActionPause
	}
}

func onHttpResponseHeaders(ctx wrapper.HttpContext, config RedisCallConfig, log wrapper.Log) types.Action {
	if ctx.GetContext("timeStamp") != nil {
		proxywasm.AddHttpResponseHeader("timeStamp", ctx.GetContext("timeStamp").(string))
	}
	if ctx.GetContext("CallTimeLeft") != nil {
		proxywasm.AddHttpResponseHeader("CallTimeLeft", ctx.GetContext("CallTimeLeft").(string))
	}
	return types.ActionContinue
}
