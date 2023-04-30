package auth

type AuthPublicPem struct {
	SystemPem string `json:"system_pem"`
	UserPem   string `json:"user_pem"`
}

type SystemInfo struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type CheckUserIsAllowedResponse struct {
	SubscriptionUUID string `json:"subscription_uuid"`
}

type UpdateApiMapRequest struct {
	SystemInfo *SystemInfo `json:"system_info"`
	Routes     []string    `json:"routes"`
}

type UpdateApiMapResponse struct {
	ApiUUIDMap map[string]string `json:"api_uuid_map"`
}
