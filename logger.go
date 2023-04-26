package micro

import (
	"log"
	"os"
	"time"
)

var Logger = &loggerImp{
	folder:    "./logs",
	stdLogger: getStdLogger(),
}

type loggerImp struct {
	folder    string
	stdLogger *log.Logger
}

type loggerInf interface {
	SetPrefix(prefix string)
	Println(v ...any)
}

func (l *loggerImp) Info(v ...any) {
	print(l.stdLogger, "INFO", v...)
	print(getFileLogger(l.folder, getLogFileName()), "INFO", v...)
}

func (l *loggerImp) Error(v ...any) {
	print(l.stdLogger, "ERROR", v...)
	print(getFileLogger(l.folder, getLogFileName()), "ERROR", v...)
}

func (l *loggerImp) Debug(v ...any) {
	print(l.stdLogger, "DEBUG", v...)
}

func (l *loggerImp) Fatal(v ...any) {
	print(l.stdLogger, "FATAL", v...)
	print(getFileLogger(l.folder, getLogFileName()), "FATAL", v...)
}

func (l *loggerImp) TraceInfo(traceID string, v ...any) {
	printTrace(l.stdLogger, traceID, v...)
	printTrace(getFileLogger(l.folder, getLogFileName()), traceID, v...)
}

func (l *loggerImp) TraceError(traceID string, v ...any) {
	printTrace(l.stdLogger, traceID, v...)
	printTrace(getFileLogger(l.folder, getLogFileName()), traceID, v...)
}

func (l *loggerImp) TraceDebug(traceID string, v ...any) {
	printTrace(l.stdLogger, traceID, v...)
}

func (l *loggerImp) TraceFatal(traceID string, v ...any) {
	printTrace(l.stdLogger, traceID, v...)
	printTrace(getFileLogger(l.folder, getLogFileName()), traceID, v...)
}

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

func print(logger loggerInf, level string, v ...any) {
	logger.SetPrefix("[" + level + "] ")
	log.Println(v...)
}

func printTrace(logger loggerInf, traceID string, v ...any) {
	logger.SetPrefix("[TRACE] (trace: " + traceID + ")")
	log.Println(v...)
}
