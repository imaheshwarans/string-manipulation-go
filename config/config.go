package config

import "io"

type Configuration struct {
	LogLevel    string
	LogFilePath string

	LogWriter io.Writer
}
