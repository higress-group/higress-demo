package main

import (
	. "github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

func main() {
	SetCtx(
		"my-plugin",
		ParseConfigBy(parseConfig),
		ProcessRequestHeadersBy(onHttpRequestHeaders),
	)
}

type MyConfig struct {
	content string
}

func parseConfig(json gjson.Result, config *MyConfig, log Log) error {
	config.content = json.Get("content").String()
	return nil
}

func onHttpRequestHeaders(ctx HttpContext, config MyConfig, log Log) types.Action {
	proxywasm.SendHttpResponse(200, nil, []byte(config.content), -1)
	log.Warnf("my custom content response: %s, request url: %s://%s%s", config.content,
		ctx.Scheme(), ctx.Host(), ctx.Path())
	return types.ActionContinue
}
