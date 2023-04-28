package logger

import (
	"github.com/ginger-go/env"
	"github.com/ginger-go/micro"
)

// Setup the log service
// Call this function in every service's main.go
// Please set it up after the api service is setup
func SetupLogService(engine *micro.Engine, logFolderPath string) {
	// Setup log service ip
	LOG_SERVICE_IP = env.String("LOG_SERVICE_IP", "")
	if LOG_SERVICE_IP == "" {
		panic("LOG_SERVICE_IP is not set") // must set LOG_SERVICE_IP
	}

	// Setup the log folder path
	LOG_FOLDER = logFolderPath

	// This api is called by logged-in user to get the system logs
	engine.GinEngine.GET("/micro/log", getLog)

	// This cron job runs every minute to send the logs to the log service
	micro.Cron(engine, "0 * * * * *", sendLog)
}
