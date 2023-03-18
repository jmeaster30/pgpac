package main

import (
	"encoding/json"
	"fmt"
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

func LogDebug(message string, args ...any) {
	Log(MyLogDebug, message, args...)
}

func LogInfo(message string, args ...any) {
	Log(MyLogInfo, message, args...)
}

func LogWarning(message string, args ...any) {
	Log(MyLogWarn, message, args...)
}

func LogError(message string, args ...any) {
	Log(MyLogError, message, args...)
}

func Log(level MyLogLevel, message string, args ...any) {
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

	fmt.Printf("\033[0;32m%s %s[%s]\033[0m %s\n", ts, color, level.String(), fullmessage)
}

func PrettyPrint(s interface{}) string {
	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		LogError("INTERNAL ERROR: Failed to pretty print object :(")
		return "<ERROR>"
	}
	return string(json)
}
