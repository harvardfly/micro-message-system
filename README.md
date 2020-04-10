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
1. 网关
2. 用户服务
3. IM服务
```