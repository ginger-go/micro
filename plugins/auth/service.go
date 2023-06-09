package auth

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ginger-go/micro"
	"github.com/ginger-go/micro/plugins/apicall"
)

// service will call the auth service to get the api uuid map at the beginning
func initApiMap(engine *micro.Engine) {
	var routes = make([]string, 0)
	for tag := range API_UUID_MAP {
		routes = append(routes, tag)
	}
	resp, err := apicall.POST[UpdateApiMapResponse](AUTH_SERVICE_IP+"/micro/api-map", &UpdateApiMapRequest{
		SystemInfo: &SystemInfo{
			UUID: SYSTEM_ID,
			Name: SYSTEM_NAME,
		},
		Routes: routes,
	}, map[string]string{
		"Authorization": "Bearer " + SYSTEM_TOKEN,
	}, "", nil)
	if err != nil {
		panic("failed to initialize service")
	}
	for k, v := range resp.Data.ApiUUIDMap {
		API_UUID_MAP[k] = v
	}
}

// service will call the auth service to get the jwt public pem at the beginning
func initPublicPem(engine *micro.Engine) {
	resp, err := apicall.GET[AuthPublicPem](AUTH_SERVICE_IP+"/micro/token", nil, map[string]string{
		"Authorization": "Bearer " + SYSTEM_TOKEN,
	}, "", nil)
	if err != nil {
		panic("failed to initialize service")
	}
	SYSTEM_TOKEN_PUBLIC_PEM = resp.Data.SystemPem
	USER_TOKEN_PUBLIC_PEM = resp.Data.UserPem
}

func checkUserHasRight(authGroup []string, systemID string, apiUUID string) bool {
	for _, authGroup := range authGroup {
		if authGroup == "*" {
			return true
		}
		resp, err := apicall.GET[GetAllowedApisResponse](AUTH_SERVICE_IP+"/micro/allowed-api", map[string]string{
			"system_id":  systemID,
			"auth_group": authGroup,
		}, map[string]string{
			"Authorization": "Bearer " + SYSTEM_TOKEN,
		}, "", nil)
		if err != nil || !resp.Success {
			log.Println("failed to check user has right", err)
			continue
		}
		for _, api := range resp.Data.AllowedApis {
			if api == apiUUID {
				return true
			}
		}
	}
	return false
}

func checkUserHasUsage(userUUID, apiUUID string) (subscriptionUUID string) {
	resp, err := apicall.GET[CheckUserIsAllowedResponse](USAGE_SERVICE_IP+"/micro/usage", map[string]string{
		"userUUID": userUUID,
		"apiUUID":  apiUUID,
	}, map[string]string{}, "", nil)
	if err != nil {
		log.Println("failed to check user is allowed", err)
		return ""
	}
	return resp.Data.SubscriptionUUID
}

func sendUsageCron() {
	m := make(map[string]int64)
	for k, v := range SUBSCRIPTION_USAGE_MAP {
		m[k] = v
	}
	SUBSCRIPTION_USAGE_MAP = make(map[string]int64)
	_, err := apicall.POST[struct{}](USAGE_SERVICE_IP+"/micro/usage", m, map[string]string{
		"Authorization": "Bearer " + SYSTEM_TOKEN,
	}, "", nil)
	if err != nil {
		log.Println("failed to send usage", err)
		return
	}
}

func GetSystemID() string {
	return SYSTEM_ID
}

func GetSystemName() string {
	return SYSTEM_NAME
}

func GetApiUUID(c *gin.Context) string {
	method := c.Request.Method
	path := c.FullPath()
	return API_UUID_MAP[method+":"+path]
}
