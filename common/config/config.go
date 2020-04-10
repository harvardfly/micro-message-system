package config

// 通用配置

type (
	UserRpcServer struct {
		ClientName string
		ServerName string
	}
	Kafka struct {
		Address []string
		Topic   string
	}

	// im_address,topic,server_name
	ImRpc struct {
		Address      string
		KafkaAddress []string
		Topic        string
		ServerName   string
	}
)
