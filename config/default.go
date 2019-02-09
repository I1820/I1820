package config

var defaultConfig = []byte(`
### configuration is in the YAML format
### and it use 2-space as tab.

core: # recieves data from rabbitmq and stores them into database
  database:
    url: mongodb://127.0.0.1:27017
  broker:
    host: 127.0.0.1:5672
    user: admin
    pass: admin
`)
