# 基于go micro的分布式即时通信系统


## 技术架构
```$xslt
1.微服务框架：go micro 基于go-plugins可插拔模式
2.服务发现：etcd
3.服务端限流：令牌桶(ratelimit)
4.熔断机制：hystrix
5.消息中间件：kafka
6.web框架：gin  gorm
7.数据库：mysql
8.token认证：jwt-go
```

## 系统模块
```$xslt
1. 用户服务:
    用户API服务 - 用户注册登录
    用户RPC服务 - 用户鉴权 网关gate调用
2. IM服务：
    IM消息发送服务 - produce消息到kafka 网关gate调用
    IM消息订阅服务 - 订阅kafka消息供websocket展示
3. 网关：
    调用userrpc验证用户token是否存在
    调用imrpc发送消息到kafka
    获取并绑定用户与IM websocket服务地址
4. Upload服务：
    文件分块上传  api没有放到网关，避免网关流量过大
```
## 系统架构图：
![Image text](https://github.com/harvardfly/micro-message-system/blob/master/docs/IM%E9%80%9A%E4%BF%A1%E7%B3%BB%E7%BB%9F%E6%9E%B6%E6%9E%84%E5%9B%BE.jpg) 

## 系统环境要求
```$xslt
golang >= 1.14
protoc >= 3.6.1
```

## 用户服务-userservice
### 生成pb
```$xslt
cd userserver/protos
protoc --proto_path=. --micro_out=. --go_out=. user.proto
```
### 启动user rpc服务
```$xslt
cd /userserver/cmd/rpc
go run user_rpc.go -f ../config/config_rpc.json
```
### 启动user api
```$xslt
cd /userserver/cmd/api
go run user_api.go -f ../config/config_api.json
```

## IM通信服务-imservice
### 生成pb
```$xslt
cd imserver/protos
protoc --proto_path=. --micro_out=. --go_out=. im.proto
```
### 启动kafka消息推送rpc服务 供网关gateway调用
```$xslt
cd imserver/cmd/
go run rpcproducer/im_rpc.go
```

### 启动kafka消息订阅消费者服务(进程) 订阅kafka消息推送到websocket
```$xslt
cd imserver/cmd/
go run imconsumer/im_server.go -f ./config/config_im_1.json
go run imconsumer/im_server.go -f ./config/config_im_2.json
```

## 网关API-gateway 需token认证并调用userservice和imservice
### 调用rpc发送消息到kafka(/send)
```$xslt
cd /userserver/cmd/rpc
go run user_rpc.go -f ../config/config_rpc.json

cd imserver/cmd/
go run rpcproducer/im_rpc.go

cd gateway/cmd
go run api/gateway_api.go
```
### 获取并绑定用户与IM服务地址(/address)
```$xslt
cd /userserver/cmd/rpc
go run user_rpc.go -f ../config/config_rpc.json

cd gateway/cmd
go run api/gateway_api.go
```

## docker 方式运行
```$xslt
详见 docker目录
```

## k8s 方式运行
```$xslt
详见 k8s目录
sh batch_deploy.sh
```