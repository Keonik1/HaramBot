package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type LogLevelType int

const (
	TRACE LogLevelType = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

var logLevelNames = map[string]LogLevelType{
	"trace":   TRACE,
	"debug":   DEBUG,
	"info":    INFO,
	"warning": WARNING,
	"error":   ERROR,
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[36m"
	colorGreen  = "\033[32m"
	colorGray   = "\033[90m"
)

var (
	LogLevel       LogLevelType
	LogDestination string // "console", "file", or "both"
	logFile        *os.File
	lineCount      = 0
	maxLines       = 10000
	maxLogFiles    = 10
)

func InitLogger(level string, destination string) error {
	level = strings.ToLower(level)
	LogDestination = destination

	if lvl, ok := logLevelNames[level]; ok {
		LogLevel = lvl
	} else {
		LogLevel = INFO
	}

	if destination == "file" || destination == "both" {
		err := setupLogFile()
		if err != nil {
			return err
		}
	}

	return nil
}

func setupLogFile() error {
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		return err
	}

	return rotateLogFile()
}

func rotateLogFile() error {
	if logFile != nil {
		logFile.Close()
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("logs/%s.log", timestamp)
	var err error
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	lineCount = 0
	cleanupOldLogs()
	return nil
}

func cleanupOldLogs() {
	files, err := filepath.Glob("logs/*.log")
	if err != nil || len(files) <= maxLogFiles {
		return
	}

	// Sort files by modification time
	type fileInfo struct {
		path string
		mod  time.Time
	}
	var infos []fileInfo
	for _, f := range files {
		stat, err := os.Stat(f)
		if err != nil {
			continue
		}
		infos = append(infos, fileInfo{f, stat.ModTime()})
	}

	// Sort by time
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].mod.Before(infos[j].mod)
	})

	for i := 0; i < len(infos)-maxLogFiles; i++ {
		os.Remove(infos[i].path)
	}
}

func shouldLog(level LogLevelType) bool {
	return level >= LogLevel
}

func logMessage(level LogLevelType, levelStr, color, format string, a ...interface{}) {
	if !shouldLog(level) {
		return
	}
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	msg := fmt.Sprintf(format, a...)

	if LogDestination == "console" || LogDestination == "both" {
		lines := strings.Split(msg, "\n")
		for _, line := range lines {
			final := fmt.Sprintf("%s%s %s: %s%s", color, timestamp, strings.ToUpper(levelStr), line, colorReset)
			fmt.Println(final)
		}
	}

	if (LogDestination == "file" || LogDestination == "both") && logFile != nil {
		lines := strings.Split(msg, "\n")
		for _, line := range lines {
			fmt.Fprintf(logFile, "%s %s: %s\n", timestamp, strings.ToUpper(levelStr), line)
		}

		lineCount += len(lines)
		if lineCount >= maxLines {
			rotateLogFile()
		}
	}
}

func LogTrace(format string, a ...interface{}) {
	logMessage(TRACE, "trace", colorGray, format, a...)
}

func LogDebug(format string, a ...interface{}) {
	logMessage(DEBUG, "debug", colorGreen, format, a...)
}

func LogInfo(format string, a ...interface{}) {
	logMessage(INFO, "info", colorBlue, format, a...)
}

func LogWarning(format string, a ...interface{}) {
	logMessage(WARNING, "warning", colorYellow, format, a...)
}

func LogError(format string, a ...interface{}) {
	logMessage(ERROR, "error", colorRed, format, a...)
}

func LogFatal(format string, a ...interface{}) {
	LogError(format, a...)
	os.Exit(1)
}
