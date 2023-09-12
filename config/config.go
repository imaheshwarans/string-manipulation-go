package config

import (
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"gopkg.in/yaml.v2"
)

type Configuration struct {
	LogLevel    string `yaml:"log-level" mapstructure:"LOG_LEVEL"`
	LogFilePath string `yaml:"log-file" mapstructure:"LOG_FILE"`

	LogWriter io.Writer
}

// this function sets the configuration file name and type
func init() {
	viper.AddConfigPath("D:\\Comcast-Delivery-interview-exercise\\string-manipulation-go\\")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
}

func (conf *Configuration) Save(filename string) error {
	configFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "Failed to create config file")
	}
	defer func() {
		derr := configFile.Close()
		if derr != nil {
			log.Fatalln("Error closing config file")
		}
	}()

	err = yaml.NewEncoder(configFile).Encode(conf)
	if err != nil {
		return errors.Wrap(err, "Failed to encode config structure")
	}
	return nil
}

func LoadConfiguration() (*Configuration, error) {
	ret := Configuration{}
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			return &ret, errors.Wrap(err, "Config file not found")
		}
		return &ret, errors.Wrap(err, "Failed to load config")
	}
	if err := viper.Unmarshal(&ret); err != nil {
		return &ret, errors.Wrap(err, "Failed to unmarshal config")
	}
	return &ret, nil
}
