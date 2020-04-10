package config

type (
	RpcConfig struct {
		Version string
		Server  struct {
			Name      string
			RateLimit int64
		}
		Etcd struct {
			Address  []string
			UserName string
			Password string
		}
		Engine struct {
			Name       string
			DataSource string
		}
	}

	ApiConfig struct {
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
		Engine struct {
			Name       string
			DataSource string
		}
	}
)
