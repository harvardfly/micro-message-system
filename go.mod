module micro-message-system

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/garyburd/redigo v1.6.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.2
	github.com/go-acme/lego/v3 v3.3.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.4.0
	github.com/gomodule/redigo v1.8.1
	github.com/gorilla/websocket v1.4.1
	github.com/jinzhu/gorm v1.9.12
	github.com/juju/ratelimit v1.0.1
	github.com/mailru/easyjson v0.7.1 // indirect
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-micro/v2 v2.6.0
	github.com/micro/go-plugins/broker/kafka v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/registry/etcdv3 v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/registry/kubernetes v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/breaker/hystrix v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit v0.0.0-20200119172437-4fe21aa238fd
	github.com/olivere/elastic v6.2.32+incompatible
	github.com/satori/go.uuid v1.2.0
	google.golang.org/protobuf v1.23.0
	gopkg.in/go-playground/validator.v8 v8.18.2
)
