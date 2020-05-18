package output

import (
    "errors"
    "strings"
)

type Level int8

const (
    DEBUG Level = iota
    INFO
    WARN
    ERROR
    PANIC
    FATAL
)

func ParseLevel(levelStr string) (l Level, err error) {
    levelStr = strings.ToUpper(levelStr)
    switch levelStr {
    case "DEBUG":
        return DEBUG, nil
    case "INFO":
        return INFO, nil
    case "WARN":
        return WARN, nil
    case "ERROR":
        return ERROR, nil
    case "PANIC":
        return PANIC, nil
    case "FATAL":
        return FATAL, nil
    }
    return Level(-1), errors.New("invalid level string")
}
