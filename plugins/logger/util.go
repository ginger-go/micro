package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func getStdLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags)
}

func getFileLogger(logFolderPath, filename string) *log.Logger {
	err := os.MkdirAll(logFolderPath, 0755)
	if err != nil {
		log.Fatalf("error creating log folder: %v", err)
	}
	path := logFolderPath + "/" + filename
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return log.New(file, "", log.LstdFlags)
}

func getLogFileName() string {
	return time.Now().Format("2006-01-02") + ".log"
}

type loggerInterface interface {
	SetPrefix(prefix string)
	Println(v ...any)
}

func print(logger loggerInterface, systemID, apiUUID, traceID string, level string, v ...any) {
	logger.SetPrefix(fmt.Sprintf("[%s] (%s/%s@%s) ", level, systemID, apiUUID, traceID))
	logger.Println(v...)
}
