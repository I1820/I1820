package config

// Default represents default configuration in YAML format with 2-space
const Default = `
debug: true
tm: # tm communicates with tm component
  url: http://127.0.0.1:1995
database:
  url: mongodb://127.0.0.1:27017
core: # core broker
  broker:
    addr: tcp://127.0.0.1:1883
`
