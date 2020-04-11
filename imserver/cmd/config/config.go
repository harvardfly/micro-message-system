package config

import "micro-message-system/common/config"

type (
	ImConfig struct {
		Version string
		Port    string
		Server  struct {
			Name      string
			RateLimit int64
		}
		Etcd struct {
			Address  []string
			UserName string
			Password string
		}
		Kafka *config.Kafka
	}

	ImRpcConfig struct {
		Version string
		Topic   string
		Server  struct {
			Name      string
			RateLimit int64
		}
		Etcd struct {
			Address  []string
			UserName string
			Password string
		}

		ImServerList []*config.ImRpc
	}
)
