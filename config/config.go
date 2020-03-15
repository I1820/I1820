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
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	// Config holds all link component configurations
	Config struct {
		Debug    bool
		TM       TM
		Database Database
		MQTT     MQTT
	}

	// TM holds I1820 Things Manager configuration
	TM struct {
		URL string
	}

	// Database holds database configuration
	Database struct {
		URL  string
		Name string
	}

	// MQTT holds MQTT configuration
	MQTT struct {
		Addr string
	}
)

// New reads configuration with viper and create configuration instance
func New() Config {
	var instance Config

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config.default")

	if err := v.ReadConfig(bytes.NewBufferString(Default)); err != nil {
		logrus.Fatalf("fatal error loading **default** config file: %s \n", err)
	}

	v.SetConfigName("config")

	if err := v.MergeInConfig(); err != nil {
		logrus.Warnf("no config file found, using defaults and environment variables")
	}

	v.SetEnvPrefix("i1820_link")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if err := v.UnmarshalExact(&instance); err != nil {
		logrus.Fatalf("unmarshal error: %s", err)
	}

	logrus.Infof("following configuration is loaded:\n%+v", instance)

	return instance
}