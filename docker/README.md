# Docker 部署

## 镜像
下面命令在根目录执行：
```
$ sudo docker build -t micro_userserver -f docker/Dockerfile-userserver .
$ sudo docker build -t micro_imserver -f docker/Dockerfile-imserver .
$ sudo docker build -t micro_websocket -f docker/Dockerfile-websocket .
$ sudo docker build -t micro_gateway -f docker/Dockerfile-gateway .
```

# 部署
下面命令在根目录执行：
```
# 启动container
$ sudo docker-compose -f docker-compose.yml up -d

# 重启，启动，关闭
$ docker-compose -f docker-compose.yml restart
$ docker-compose -f docker-compose.yml start
$ docker-compose -f docker-compose.yml stop
```