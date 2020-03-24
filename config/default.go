package config

// Default return default configuration
// nolint: gomnd
func Default() Config {
	return Config{
		TM: TM{
			URL: "http://127.0.0.1:1378",
		},
		Database: Database{
			URL:  "mongodb://127.0.0.1:27017",
			Name: "i1820",
		},
		Rabbitmq: Rabbitmq{
			Host:           "127.0.0.1",
			Port:           5672,
			User:           "guest",
			Pass:           "guest",
			RetryThreshold: 10,
		},
		MQTT: MQTT{
			Addr: "tcp://127.0.0.1:1883",
		},
		Docker: Docker{
			Host: "",
			Runner: Runner{
				Database: Database{
					URL:  "mongodb://172.17.0.1",
					Name: "i1820",
				},
				Rabbitmq: Rabbitmq{
					Host:           "172.17.0.1",
					Port:           5672,
					User:           "guest",
					Pass:           "guest",
					RetryThreshold: 10,
				},
			},
		},
	}
}
