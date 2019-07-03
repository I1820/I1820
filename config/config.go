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

package config

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all tm component configurations
type Config struct {
	Debug    bool
	Database struct {
		URL string
	}
}

// New reads configuration with viper
func New() Config {
	var instance Config

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config.default")

	if err := v.ReadConfig(bytes.NewBufferString(Default)); err != nil {
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

	v.SetEnvPrefix("i1820_tm")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.UnmarshalExact(&instance); err != nil {
		log.Printf("configuration: %s", err)
	}
	fmt.Printf("Following configuration is loaded:\n%+v\n", instance)

	return instance
}
