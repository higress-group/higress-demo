package main

import (
	"github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

func main() {
	wrapper.SetCtx(
		"my-plugin",
		wrapper.ParseConfigBy(parseConfig),
		wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
	)
}

type MyConfig struct {
	mockEnable bool
}

func parseConfig(json gjson.Result, config *MyConfig, log wrapper.Log) error {
	config.mockEnable = json.Get("mockEnable").Bool()
	return nil
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config MyConfig, log wrapper.Log) types.Action {
	proxywasm.AddHttpRequestHeader("hello", "world")
	if config.mockEnable {
		proxywasm.SendHttpResponse(200, nil, []byte("hello world"), -1)
	}
	return types.ActionContinue
}
