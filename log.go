//Package gologging provides a logger implementation based on the github.com/op/go-logging pkg
package gologging

import (
	"fmt"
	"io"
	"os"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	gologging "github.com/op/go-logging"
)

// Namespace is the key to look for extra configuration details
const Namespace = "github_com/devopsfaith/krakend-gologging"

var (
	// ErrEmptyValue is the error returned when there is no config under the namespace
	ErrWrongConfig = fmt.Errorf("getting the extra config for the krakend-gologging module")
	// LoggingPattern is the pattern to use for rendering the logs
	LoggingPattern = ` %{time:2006/01/02 - 15:04:05.000} %{color}▶ %{level:.6s}%{color:reset} %{message}`
)

// NewLogger returns a krakend logger wrapping a gologging logger
func NewLogger(cfg config.ExtraConfig, ws ...io.Writer) (logging.Logger, error) {
	logConfig, ok := ConfigGetter(cfg).(Config)
	if !ok {
		return nil, ErrWrongConfig
	}
	module := "KRAKEND"
	loggr := gologging.MustGetLogger(module)

	backends := []gologging.Backend{}
	var b gologging.Backend
	if logConfig.StdOut {
		b = gologging.NewLogBackend(os.Stdout, logConfig.Prefix, 0)
		backends = append(backends, b)
	}

	for _, w := range ws {
		b = gologging.NewLogBackend(w, logConfig.Prefix, 0)
		backends = append(backends, b)
	}

	if logConfig.Syslog {
		var err error
		b, err = gologging.NewSyslogBackend(logConfig.Prefix)
		if err != nil {
			return nil, err
		}
		backends = append(backends, b)
	}

	for i, b := range backends {
		format := gologging.MustStringFormatter(LoggingPattern)
		backendLeveled := gologging.AddModuleLevel(gologging.NewBackendFormatter(b, format))
		logLevel, err := gologging.LogLevel(logConfig.Level)
		if err != nil {
			return nil, err
		}
		backendLeveled.SetLevel(logLevel, module)
		backends[i] = backendLeveled
	}

	gologging.SetBackend(backends...)
	return Logger{loggr}, nil
}

// ConfigGetter implements the config.ConfigGetter interface
func ConfigGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := Config{}
	if v, ok := tmp["stdout"]; ok {
		cfg.StdOut = v.(bool)
	}
	if v, ok := tmp["syslog"]; ok {
		cfg.Syslog = v.(bool)
	}
	if v, ok := tmp["level"]; ok {
		cfg.Level = v.(string)
	}
	if v, ok := tmp["prefix"]; ok {
		cfg.Prefix = v.(string)
	}
	return cfg
}

// Config is the custom config struct containing the params for the logger
type Config struct {
	Level  string
	StdOut bool
	Syslog bool
	Prefix string
}

// Logger is a wrapper over a github.com/op/go-logging logger
type Logger struct {
	logger *gologging.Logger
}

// Debug implements the logger interface
func (l Logger) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}

// Info implements the logger interface
func (l Logger) Info(v ...interface{}) {
	l.logger.Info(v...)
}

// Warning implements the logger interface
func (l Logger) Warning(v ...interface{}) {
	l.logger.Warning(v...)
}

// Error implements the logger interface
func (l Logger) Error(v ...interface{}) {
	l.logger.Error(v...)
}

// Critical implements the logger interface
func (l Logger) Critical(v ...interface{}) {
	l.logger.Critical(v...)
}

// Fatal implements the logger interface
func (l Logger) Fatal(v ...interface{}) {
	l.logger.Fatal(v...)
}
