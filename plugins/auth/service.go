package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ginger-go/micro"
	"github.com/ginger-go/micro/plugins/apicall"
)

// service will call the auth service to get the api uuid map at the beginning
func initApiMap(engine *micro.Engine) {
	routeInfos := engine.GinEngine.Routes()
	var routes = make([]string, 0)
	for _, routeInfo := range routeInfos {
		routes = append(routes, routeInfo.Method+":"+routeInfo.Path)
	}
	resp, err := apicall.POST[map[string]string](AUTH_SERVICE_IP+"/micro/api-map", routes, map[string]string{
		"Authorization": "Bearer " + SYSTEM_TOKEN,
	}, "", nil)
	if err != nil {
		panic("failed to initialize service")
	}
	for k, v := range *resp.Data {
		API_UUID_MAP[k] = v
	}
}

func GetSystemID() string {
	return SYSTEM_ID
}

func GetApiUUID(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path
	return API_UUID_MAP[method+":"+path]
}
