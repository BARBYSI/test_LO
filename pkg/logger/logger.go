package logger

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	loggerOnce sync.Once
	logger     *Logger
)

func Init(cfg *LoggerConfig) {
	loggerOnce.Do(func() {
		if cfg.Enabled {
			logger = NewLogger(cfg)
		} else {
			logger = &Logger{cfg: cfg}
		}
	})
}

type LoggerConfig struct {
	Enabled bool
}

type Logger struct {
	cfg   *LoggerConfig
	logCh chan *logMessage
}

type logMessage struct {
	timestamp string
	level     string
	message   string
}

func NewLogger(cfg *LoggerConfig) *Logger {
	logCh := make(chan *logMessage, 100)
	go func() {
		for msg := range logCh {
			logLine := fmt.Sprintf("[%s] [%s] %s", msg.timestamp, msg.level, msg.message)
			log.Println(logLine)
		}
	}()

	return &Logger{
		cfg:   cfg,
		logCh: logCh,
	}
}

func (l *Logger) log(logMsg *logMessage) {
	if l.cfg.Enabled && logMsg != nil {
		l.logCh <- logMsg
	}
}

func LogError(msg string) {
	if logger == nil {
		return
	}
	logger.log(&logMessage{
		timestamp: time.Now().Format(time.RFC3339),
		level:     "ERROR",
		message:   msg,
	})
}

func LogInfo(msg string) {
	if logger == nil {
		return
	}
	logger.log(&logMessage{
		timestamp: time.Now().Format(time.RFC3339),
		level:     "INFO",
		message:   msg,
	})
}
