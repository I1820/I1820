package config

// Default represents default configuration in YAML format with 2-space
const Default = `
tm:
  url: http://127.0.0.1:1995
database:
  url: mongodb://127.0.0.1:27017
  name: i1820
mqtt:
  addr: tcp://127.0.0.1:1883
rabbitmq:
  host: 127.0.0.1
  port: 5672
  user: guest
  pass: guest
docker:
  host: 127.0.0.1
`
