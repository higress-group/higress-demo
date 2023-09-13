package main

import (
    "errors"
    "net/http"
    // "strings"

    "github.com/alibaba/higress/plugins/wasm-go/pkg/wrapper"
    "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
    "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
    "github.com/tidwall/gjson"
)

func main() {
    wrapper.SetCtx(
        "http-call",
        wrapper.ParseConfigBy(parseConfig),
        wrapper.ProcessRequestHeadersBy(onHttpRequestHeaders),
    )
}

type MyConfig struct {
    // 用于发起HTTP调用client
    client      wrapper.HttpClient
    // 请求url
    requestPath string
    // 根据该key取出调用服务的应答头对应字段，再设置到原始请求的请求头，key为此配置项
    tokenHeader string
}

func parseConfig(json gjson.Result, config *MyConfig, log wrapper.Log) error {
    // config.tokenHeader = json.Get("tokenHeader").String()
    // if config.tokenHeader == "" {
    //     return errors.New("missing tokenHeader in config")
    // }
    config.requestPath = json.Get("requestPath").String()
    if config.requestPath == "" {
        return errors.New("missing requestPath in config")
    }
    serviceSource := json.Get("serviceSource").String()
    // 固定地址和DNS类型的serviceName，为控制台中创建服务时指定
    // nacos和k8s来源的serviceName，即服务注册时指定的原始名称
    serviceName := json.Get("serviceName").String()
    servicePort := json.Get("servicePort").Int()
    if serviceName == "" || servicePort == 0 {
        return errors.New("invalid service config")
    }
    switch serviceSource {
    case "k8s":
        namespace := json.Get("namespace").String()
        config.client = wrapper.NewClusterClient(wrapper.K8sCluster{
            ServiceName: serviceName,
            Namespace:   namespace,
            Port:        servicePort,
        })
        return nil
    case "nacos":
        namespace := json.Get("namespace").String()
        config.client = wrapper.NewClusterClient(wrapper.NacosCluster{
            ServiceName: serviceName,
            NamespaceID: namespace,
            Port:        servicePort,
        })
        return nil
    case "ip":
        config.client = wrapper.NewClusterClient(wrapper.StaticIpCluster{
            ServiceName: serviceName,
            Port:        servicePort,
        })
        return nil
    case "dns":
        domain := json.Get("domain").String()
        config.client = wrapper.NewClusterClient(wrapper.DnsCluster{
            ServiceName: serviceName,
            Port:        servicePort,
            Domain:      domain,
        })
        return nil
    default:
        return errors.New("unknown service source: " + serviceSource)
    }
}

func onHttpRequestHeaders(ctx wrapper.HttpContext, config MyConfig, log wrapper.Log) types.Action {
    // 使用client的Get方法发起HTTP Get调用，此处省略了timeout参数，默认超时时间500毫秒
    config.client.Get(config.requestPath, nil,
        // 回调函数，将在响应异步返回时被执行
        func(statusCode int, responseHeaders http.Header, responseBody []byte) {
            // 用defer，在函数返回前恢复原始请求流程，继续往下处理，才能正常转发给后端服务
            defer proxywasm.ResumeHttpRequest()
            // 请求没有返回200状态码，进行处理
            if statusCode != http.StatusOK {
				        log.Errorf("http call failed, status: %d", statusCode)
				        return
			      }
            // 打印响应的HTTP状态码和应答body
            log.Infof("get status: %d, request body: %s", statusCode, responseBody)
            // 从应答头中解析token字段设置到原始请求头中
            token := responseHeaders.Get(config.tokenHeader)
            if token != "" {
                proxywasm.AddHttpRequestHeader(config.tokenHeader, token)
            }
        })
    // 需要等待异步回调完成，返回Pause状态，可以被ResumeHttpRequest恢复
    return types.ActionPause
}