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
1. 用户服务
2. 网关
3. IM服务
```

## 系统环境要求
```$xslt
golang >= 1.14
protoc >= 3.6.1
```

## userservice
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