package logger

// Info logs the information should be logged
func Info(systemID, apiUUID, traceID string, v ...any) {
	stdLogger := getStdLogger()
	fileLogger := getFileLogger(LOG_FOLDER, getLogFileName())
	print(stdLogger, systemID, apiUUID, traceID, LOG_LEVEL_INFO, v...)
	print(fileLogger, systemID, apiUUID, traceID, LOG_LEVEL_INFO, v...)
}

// Error logs the expected error
func Error(systemID, apiUUID, traceID string, v ...any) {
	stdLogger := getStdLogger()
	fileLogger := getFileLogger(LOG_FOLDER, getLogFileName())
	print(stdLogger, systemID, apiUUID, traceID, LOG_LEVEL_ERROR, v...)
	print(fileLogger, systemID, apiUUID, traceID, LOG_LEVEL_ERROR, v...)
}

// Fatal logs the unexpected error
func Fatal(systemID, apiUUID, traceID string, v ...any) {
	stdLogger := getStdLogger()
	fileLogger := getFileLogger(LOG_FOLDER, getLogFileName())
	print(stdLogger, systemID, apiUUID, traceID, LOG_LEVEL_FATAL, v...)
	print(fileLogger, systemID, apiUUID, traceID, LOG_LEVEL_FATAL, v...)
}

// Debug logs the debug information, should be disabled in production
func Debug(systemID, apiUUID, traceID string, v ...any) {
	stdLogger := getStdLogger()
	print(stdLogger, systemID, apiUUID, traceID, LOG_LEVEL_DEBUG, v...)
}
