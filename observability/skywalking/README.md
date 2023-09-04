# 部署过程
1. 部署网关：参见 [Higress安装文档](http://higress.io/zh-cn/docs/user/quickstart/)
2. 部署skywalking：`kubectl apply -f skywalking.yaml`
3. 部署springboot服务：`kubectl apply -f app-yamls`