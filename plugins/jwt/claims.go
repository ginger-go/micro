package jwt

// Claims is designed for microservice ecosystem.
type Claims struct {
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	IP        string                 `json:"ip"`
	TokenType string                 `json:"token_type"`  // system-token, access-token, refresh-token, api-token
	AuthGroup []string               `json:"auth_groups"` // map[api_uuid]system_uuid
	DataSets  []string               `json:"data_sets"`   // [data_set_uuid...]
	Data      map[string]interface{} `json:"data"`
}

// NewClaims creates a new Claims.
func NewClaims(uuid, name, ip, tokenType string, authGroups []string, dataSets []string) *Claims {
	return &Claims{
		UUID:      uuid,
		Name:      name,
		IP:        ip,
		TokenType: tokenType,
		AuthGroup: authGroups,
		DataSets:  dataSets,
		Data:      make(map[string]interface{}),
	}
}

// Set sets a key-value pair to Claims.Data.
func (c *Claims) Set(key string, value interface{}) {
	c.Data[key] = value
}

// Get gets a value from Claims.Data.
func (c *Claims) Get(key string) interface{} {
	return c.Data[key]
}
