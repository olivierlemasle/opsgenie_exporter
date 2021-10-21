package config

import (
	"os"

	configReader "github.com/ccding/go-config-reader/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var log *logrus.Logger
var configPath string
var config *configReader.Config

func AddFlags(a *kingpin.Application) {
	a.Flag("config", "").StringVar(&configPath)
}

func LoadConfiguration(logger *logrus.Logger) {
	log = logger
	if configPath == "" {
		log.Debugf("No config path specified using --config; using default path")
		configPath = "conf/lamp.conf"
	}
	log.Debugf("Config path: %v", configPath)

	load(configPath)
}

func load(confPath string) {
	if _, err := os.Stat(confPath); !os.IsNotExist(err) {
		conf := configReader.NewConfig(confPath)
		err = conf.Read()
		if err == nil {
			config = conf
		} else {
			log.Errorf("Could not read config file: %v", err.Error())
		}
	} else {
		log.Errorf("Could not read config file: %v", err.Error())
	}
}

// Get method returns the configuration properties value according to the key.
func Get(key string) string {
	if config != nil {
		return config.Get("", key)
	}
	return ""
}
