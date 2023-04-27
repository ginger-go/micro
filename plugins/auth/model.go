package auth

type AuthPublicPem struct {
	SystemPem string `json:"token_pem"`
	UserPem   string `json:"user_pem"`
}

type SystemInfo struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
