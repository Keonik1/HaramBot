package tools

import (
	"fmt"
	"log"
	"os"
)

// Цвета ANSI
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[36m" // Бирюзовый / сине-зеленый
	colorGreen  = "\033[32m"
	colorGray   = "\033[90m"
)

var (
	infoLogger    = log.New(os.Stdout, "", log.LstdFlags)
	warningLogger = log.New(os.Stdout, "", log.LstdFlags)
	errorLogger   = log.New(os.Stderr, "", log.LstdFlags)
	debugLogger   = log.New(os.Stdout, "", log.LstdFlags)
	traceLogger   = log.New(os.Stdout, "", log.LstdFlags)
)

func LogInfo(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	infoLogger.Println(colorBlue + "INFO: " + msg + colorReset)
}

func LogWarning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	warningLogger.Println(colorYellow + "WARNING: " + msg + colorReset)
}

func LogError(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	errorLogger.Println(colorRed + "ERROR: " + msg + colorReset)
}

func LogDebug(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	debugLogger.Println(colorGreen + "DEBUG: " + msg + colorReset)
}

func LogTrace(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	traceLogger.Println(colorGray + "TRACE: " + msg + colorReset)
}
