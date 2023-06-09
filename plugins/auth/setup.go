package auth

import (
	"time"

	"github.com/ginger-go/env"
	"github.com/ginger-go/micro"
	"github.com/ginger-go/micro/plugins/midware"
)

// Setup the auth service
// Call this function in every service's main.go
// Please set it up after the api service is setup
func SetupAuthService(engine *micro.Engine) {
	// Setup auth service ip
	AUTH_SERVICE_IP = env.String("AUTH_SERVICE_IP", "")
	if AUTH_SERVICE_IP == "" {
		panic("AUTH_SERVICE_IP is not set") // must set AUTH_SERVICE_IP
	}

	// This api is called by public to get the system info
	engine.GinEngine.GET("/micro/info", getSystemInfoHandler, midware.RateLimited(time.Minute, 30))

	// This api is called by the auth service to update the public pem
	// The public pem is used to verify the jwt token
	engine.GinEngine.POST("/micro/token", updatePublicPemHandler, midware.RateLimited(time.Minute, 30), AuthServiceOnly)

	// This cron will send the usage to the usage service every minute
	micro.Cron(engine, "0 * * * * *", sendUsageCron)

	// This is init func for initialize the api uuid map
	initApiMap(engine)

	// This is init func for initialize the public pem
	initPublicPem(engine)
}
