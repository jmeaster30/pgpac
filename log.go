package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type MyLogLevel int

const (
	MyLogDebug MyLogLevel = iota
	MyLogInfo
	MyLogWarn
	MyLogError
)

func (lvl MyLogLevel) String() string {
	if lvl == MyLogDebug {
		return "DEBUG"
	} else if lvl == MyLogInfo {
		return "INFO"
	} else if lvl == MyLogWarn {
		return "WARN"
	} else if lvl == MyLogError {
		return "ERROR"
	}
	return "----"
}

func (lvl MyLogLevel) ShouldPrint(minimumLogLevel string) bool {
	if lvl == MyLogDebug && (strings.ToLower(minimumLogLevel) == "info" || strings.ToLower(minimumLogLevel) == "warn" || strings.ToLower(minimumLogLevel) == "error") {
		return false
	} else if lvl == MyLogInfo && (strings.ToLower(minimumLogLevel) == "warn" || strings.ToLower(minimumLogLevel) == "error") {
		return false
	} else if lvl == MyLogWarn && strings.ToLower(minimumLogLevel) == "error" {
		return false
	}
	return true
}

func LogDebug(minimumLogLevel string, message string, args ...any) {
	Log(minimumLogLevel, MyLogDebug, message, args...)
}

func LogInfo(minimumLogLevel string, message string, args ...any) {
	Log(minimumLogLevel, MyLogInfo, message, args...)
}

func LogWarning(minimumLogLevel string, message string, args ...any) {
	Log(minimumLogLevel, MyLogWarn, message, args...)
}

func LogError(minimumLogLevel string, message string, args ...any) {
	Log(minimumLogLevel, MyLogError, message, args...)
}

func LogAssert(minimumLogLevel string, shouldPrint bool, message string, args ...any) {
	if shouldPrint {
		Log(minimumLogLevel, MyLogError, "\033[1;31mASSERTION FAILED!!!\033[0m "+message, args...)
	}
}

func Log(minimumLogLevel string, level MyLogLevel, message string, args ...any) {
	fullmessage := fmt.Sprintf(message, args...)
	ts := time.Now().Format("2006-01-02 15:04:05")
	color := "\033[0m"
	if level == MyLogDebug {
		color = "\033[0;35m"
	} else if level == MyLogInfo {
		color = "\033[0;36m"
	} else if level == MyLogWarn {
		color = "\033[0;33m"
	} else if level == MyLogError {
		color = "\033[0;31m"
	}

	if level.ShouldPrint(minimumLogLevel) {
		fmt.Printf("\033[0;32m%s %s[%s]\033[0m %s\n", ts, color, level.String(), fullmessage)
	}
}

func PrettyPrint(s interface{}) string {
	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		LogError("error", "INTERNAL ERROR: Failed to pretty print object :(")
		return "<ERROR>"
	}
	return string(json)
}
