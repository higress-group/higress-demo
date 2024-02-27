module github.com/alibaba/higress/plugins/wasm-go/extensions/redis-test

go 1.18

require (
	github.com/alibaba/higress/plugins/wasm-go v0.0.0-00010101000000-000000000000
	github.com/tetratelabs/proxy-wasm-go-sdk v0.19.1-0.20220822060051-f9d179a57f8c
	github.com/tidwall/gjson v1.14.3
	github.com/tidwall/resp v0.1.1
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/higress-group/nottinygc v0.0.0-20231101025119-e93c4c2f8520 // indirect
	github.com/magefile/mage v1.14.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
)

replace (
	github.com/alibaba/higress/plugins/wasm-go => ../..
	github.com/tetratelabs/proxy-wasm-go-sdk => github.com/higress-group/proxy-wasm-go-sdk v0.0.0-20240105034322-9a6ac242c3dd
)