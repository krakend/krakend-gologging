package gologging

import (
	"bytes"
	"regexp"
	"testing"

	gologging "github.com/op/go-logging"
)

const (
	debugMsg    = "Debug msg"
	infoMsg     = "Info msg"
	warningMsg  = "Warning msg"
	errorMsg    = "Error msg"
	criticalMsg = "Critical msg"
)

func TestNewLogger(t *testing.T) {
	levels := []string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}
	regexps := []*regexp.Regexp{
		regexp.MustCompile(debugMsg),
		regexp.MustCompile(infoMsg),
		regexp.MustCompile(warningMsg),
		regexp.MustCompile(errorMsg),
		regexp.MustCompile(criticalMsg),
	}

	for i, level := range levels {
		output, err := logSomeStuff(level)
		if err != nil {
			t.Error(err)
			return
		}
		for j := i; j < len(regexps); j++ {
			if !regexps[j].MatchString(output) {
				t.Errorf("The output doesn't contain the expected msg for the level: %s. [%s]", level, output)
			}
		}
	}
}

func TestNewLogger_unknownLevel(t *testing.T) {
	_, err := NewLogger(newExtraConfig("UNKNOWN"), bytes.NewBuffer(make([]byte, 1024)))
	if err == nil {
		t.Error("The factory didn't return the expected error")
		return
	}
	if err != gologging.ErrInvalidLogLevel {
		t.Errorf("The factory didn't return the expected error. Got: %s", err.Error())
	}
}

func newExtraConfig(level string) map[string]interface{} {
	return map[string]interface{}{
		Namespace: map[string]interface{}{
			"level":  level,
			"prefix": "pref",
			"syslog": false,
			"stdout": true,
		},
	}
}

func logSomeStuff(level string) (string, error) {
	buff := bytes.NewBuffer(make([]byte, 1024))
	logger, err := NewLogger(newExtraConfig(level), buff)
	if err != nil {
		return "", err
	}

	logger.Debug(debugMsg)
	logger.Info(infoMsg)
	logger.Warning(warningMsg)
	logger.Error(errorMsg)
	logger.Critical(criticalMsg)

	return buff.String(), nil
}
