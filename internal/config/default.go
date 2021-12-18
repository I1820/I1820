package config

import "github.com/I1820/I1820/internal/db"

// Default return default configuration

func Default() Config {
	return Config{
		TM: TM{
			URL: "http://127.0.0.1:1378",
		},
		Database: db.Config{
			URL:  "mongodb://127.0.0.1:27017",
			Name: "i1820",
		},
		NATS: NATS{
			URL: "nats://127.0.0.1:4222",
		},
		MQTT: MQTT{
			Addr: "tcp://127.0.0.1:1883",
		},
		Docker: Docker{
			Host: "",
			Runner: Runner{
				Database: db.Config{
					URL:  "mongodb://172.17.0.1",
					Name: "i1820",
				},
				NATS: NATS{
					URL: "nats://127.0.0.1:4222",
				},
			},
		},
	}
}
