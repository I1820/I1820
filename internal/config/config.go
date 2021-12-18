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
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/sirupsen/logrus"
)

const (
	// Namespace of I1820.
	Namespace = "I1820"

	// Prefix indicates environment variables prefix.
	Prefix = "I1820_"
)

// Component ports are defined here.
const (
	// LinkPort is a port of link component.
	LinkPort = 0
	// TMPort is a port of thing manager component.
	TMPort = 1378
	// DMPort is port of data manager component.
	DMPort = 1373
	// PMPort is port of project manager component.
	PMPort = 1999
)

type (
	// Config holds all link component configurations.
	Config struct {
		TM       TM       `koanf:"tm"`
		Database Database `koanf:"database"`
		NATS     NATS     `koanf:"nats"`
		MQTT     MQTT     `mapstrcuture:"mqtt"`
		Docker   Docker   `koanf:"docker"`
	}

	// TM holds I1820 Things Manager configuration.
	TM struct {
		URL string `koanf:"url"`
	}

	// Database holds database configuration.
	Database struct {
		URL  string `koanf:"url"`
		Name string `koanf:"name"`
	}

	// MQTT holds MQTT configuration.
	MQTT struct {
		Addr string `koanf:"addr"`
	}

	// NATS hodls NATS configuration.
	NATS struct {
		URL string `koanf:"url"`
	}

	// Docker holds Docker Host configuration for running the runners.
	Docker struct {
		Host   string `koanf:"host"`
		Runner Runner `koanf:"runner"`
	}

	// Runner contains the information that are required in runners for get and store the data.
	Runner struct {
		Database Database `koanf:"database"`
		NATS     NATS     `koanf:"nats"`
	}
)

// New reads configuration with koanf and create configuration instance.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from its struct
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider("config.yml"), yaml.Parser()); err != nil {
		logrus.Errorf("error loading config.yml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "__", ".")
	}), nil); err != nil {
		logrus.Errorf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	logrus.Infof("following configuration is loaded:\n%+v", instance)

	return instance
}
