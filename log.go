package main

import (
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
	fmt.Printf("%s [%s] %s\n", ts, level.String(), fullmessage)
}
