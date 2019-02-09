package config

var defaultConfig = []byte(`
### configuration is in the YAML format
### and it use 2-space as tab.
debug: true
database:
  url: mongodb://127.0.0.1:27017
core: # recieves data from rabbitmq and stores them into database
  broker:
    host: 127.0.0.1:5672
    user: admin
    pass: admin
`)
