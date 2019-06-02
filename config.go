/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 30-01-2019
 * |
 * | File Name:     config.go
 * +===============================================
 */

package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all link component configurations
type Config struct {
	PM struct {
		URL string
	}
	Database struct {
		URL string
	}
	Core struct {
		Broker struct {
			Addr string
		}
	}
}

// config reads configuration with viper
func config() Config {
	var defaultConfig = []byte(`
### configuration is in the YAML format
### and it use 2-space as tab.
pm: # pm communicates with pm component
  url: http://127.0.0.1:8080
database:
  url: mongodb://127.0.0.1:27017
core: # core broker
  broker:
    addr: tcp://127.0.0.1:1883
`)

	var instance Config

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
		logrus.Fatalf("fatal error loading **default** config array: %s \n", err)
	}

	v.SetConfigName("config")

	if err := v.MergeInConfig(); err != nil {
		switch err.(type) {
		default:
			logrus.Fatalf("fatal error loading config file: %s \n", err)
		case viper.ConfigFileNotFoundError:
			logrus.Infof("no config file found. Using defaults and environment variables")
		}
	}

	v.SetEnvPrefix("i1820_link")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.UnmarshalExact(&instance); err != nil {
		logrus.Infof("configuration: %s", err)
	}
	fmt.Printf("Following configuration is loaded:\n%+v\n", instance)

	return instance
}
