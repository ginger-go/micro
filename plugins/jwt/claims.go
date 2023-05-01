package jwt

// Claims is designed for microservice ecosystem.
type Claims struct {
	UUID       string                 `json:"uuid"`
	Name       string                 `json:"name"`
	IP         string                 `json:"ip"`
	IsRoot     bool                   `json:"is_root"`     // root user is the billed user
	IsAdmin    bool                   `json:"is_admin"`    // admin user is the user who can manage the system (non-client)
	TokenType  string                 `json:"token_type"`  // system-token, access-token, refresh-token, api-token
	AuthGroup  []string               `json:"auth_groups"` // []auth_group_uuid...
	Workspaces []string               `json:"Workspace"`   // []workspace_uuid...
	Data       map[string]interface{} `json:"data"`
}

// NewClaims creates a new Claims.
func NewClaims(uuid, name, ip, tokenType string, isRoot, isAdmin bool, authGroups []string, Workspaces []string) *Claims {
	return &Claims{
		UUID:       uuid,
		Name:       name,
		IP:         ip,
		IsRoot:     isRoot,
		IsAdmin:    isAdmin,
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

// HasWorkspace checks if the Claims has the workspace.
func (c *Claims) HasWorkspace(workspaceUUID string) bool {
	for _, workspace := range c.Workspaces {
		if workspace == workspaceUUID {
			return true
		}
	}
	return false
}
