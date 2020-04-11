module micro-message-system

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.2
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.3
	github.com/gorilla/websocket v1.4.1
	github.com/jinzhu/gorm v1.9.12
	github.com/juju/ratelimit v1.0.1
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/broker/kafka v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/registry/etcdv3 v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit v0.0.0-20200119172437-4fe21aa238fd
	github.com/satori/go.uuid v1.2.0
	golang.org/x/tools v0.0.0-20191029190741-b9c20aec41a5 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
)
