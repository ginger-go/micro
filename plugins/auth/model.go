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
