package jwt

// Claims is designed for microservice ecosystem.
type Claims struct {
	UUID       string                 `json:"uuid"`
	Name       string                 `json:"name"`
	IP         string                 `json:"ip"`
	IsRoot     bool                   `json:"is_root"`
	TokenType  string                 `json:"token_type"`  // system-token, access-token, refresh-token, api-token
	AuthGroup  []string               `json:"auth_groups"` // []auth_group_uuid...
	Workspaces []string               `json:"Workspace"`   // []workspace_uuid...
	Data       map[string]interface{} `json:"data"`
}

// NewClaims creates a new Claims.
func NewClaims(uuid, name, ip, tokenType string, isRoot bool, authGroups []string, Workspaces []string) *Claims {
	return &Claims{
		UUID:       uuid,
		Name:       name,
		IP:         ip,
		IsRoot:     isRoot,
		TokenType:  tokenType,
		AuthGroup:  authGroups,
		Workspaces: Workspaces,
		Data:       make(map[string]interface{}),
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
