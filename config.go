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
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all dm component configurations
type Config struct {
	Debug    bool
	Database struct {
		URL string
	}
}

// config reads configuration with viper
func config() Config {
	var defaultConfig = []byte(`
### configuration is in the YAML format
### and it use 2-space as tab.
debug: true
database:
  url: mongodb://127.0.0.1:27017
`)
	var instance Config

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config.default")

	if err := v.ReadConfig(bytes.NewReader(defaultConfig)); err != nil {
		log.Fatalf("Fatal error loading **default** config file: %s \n", err)
	}

	v.SetConfigName("config")

	if err := v.MergeInConfig(); err != nil {
		switch err.(type) {
		default:
			log.Fatalf("Fatal error loading config file: %s \n", err)
		case viper.ConfigFileNotFoundError:
			log.Printf("No config file found. Using defaults and environment variables")
		}
	}

	v.SetEnvPrefix("i1820_dm")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.UnmarshalExact(&instance); err != nil {
		log.Printf("configuration: %s", err)
	}
	fmt.Printf("Following configuration is loaded:\n%+v\n", instance)

	return instance
}
