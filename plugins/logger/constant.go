package logger

// LOG_SERVICE_IP is the IP address of the auth service
// Please set it to the environment variable LOG_SERVICE_IP
var LOG_SERVICE_IP string

// This is the folder where the logs will be stored
var LOG_FOLDER = "./logs"

const (
	LOG_LEVEL_INFO  = "INFO"
	LOG_LEVEL_ERROR = "ERROR"
	LOG_LEVEL_FATAL = "FATAL"
	LOG_LEVEL_DEBUG = "DEBUG"
)
