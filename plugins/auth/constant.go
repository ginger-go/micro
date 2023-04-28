package auth

import "github.com/ginger-go/env"

// AUTH_SERVICE_IP is the IP address of the auth service
// Please set it to the environment variable AUTH_SERVICE_IP
var AUTH_SERVICE_IP string

// USAGE_SERVICE_IP is the IP address of the usage service
// Please set it to the environment variable USAGE_SERVICE_IP
var USAGE_SERVICE_IP string

// These public pem are used to verify the jwt token
var USER_TOKEN_PUBLIC_PEM string
var SYSTEM_TOKEN_PUBLIC_PEM string

// This is the system token for calling the api internally
// Please set it to the environment variable SYSTEM_TOKEN
var SYSTEM_TOKEN string

// These are the system basic info
// Please set it to the environment variable SYSTEM_UUIDs and SYSTEM_NAME
var SYSTEM_ID string
var SYSTEM_NAME string

// These are the auth related error code and message
const (
	ERR_CODE_UNAUTHORIZED = "b97cf20d-42b6-470e-9e08-b4bb852c3811"
	ERR_CODE_FORBIDDEN    = "7792176d-0196-4a57-a959-93062c2b9b41"
	ERR_MSG_UNAUTHORIZED  = "Unauthorized"
	ERR_MSG_FORBIDDEN     = "Forbidden"
)

// This map is recorded the api and its uuid, which is determined by the auth service
var API_UUID_MAP = make(map[string]string)

// This map is recorded the subscription usage
// This map will send to the usage service periodically
var SUBSCRIPTION_USAGE_MAP = make(map[string]int64)

func init() {
	SYSTEM_ID = env.String("SYSTEM_ID", "")
	if SYSTEM_ID == "" {
		panic("SYSTEM_ID is empty")
	}
	SYSTEM_NAME = env.String("SYSTEM_NAME", "")
	if SYSTEM_NAME == "" {
		panic("SYSTEM_NAME is empty")
	}
}
