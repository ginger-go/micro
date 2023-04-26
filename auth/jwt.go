package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateRSAKeyPair() (privKeyStr string, pubKeyStr string, err1 error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	))

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	))

	return privKeyPem, pubKeyPem, nil
}

type Claims struct {
	ID                uint                   `json:"id"`
	Username          string                 `json:"username"`
	Role              string                 `json:"role"` // system_admin, system_user, system_token, root_user, workspace_user, api_token
	AllowedSystems    []string               `json:"allowed_systems"`
	AllowedWorkspaces []string               `json:"allowed_workspaces"`
	Msg               string                 `json:"msg"`
	Nonce             string                 `json:"nonce"`
	Sig               string                 `json:"sig"`
	Data              map[string]interface{} `json:"data"`
}

func NewClaims(id uint, username string, role string, allowedSystems []string, allowedWorkspaces []string) *Claims {
	return &Claims{
		ID:                id,
		Username:          username,
		Role:              role,
		AllowedSystems:    allowedSystems,
		AllowedWorkspaces: allowedWorkspaces,
		Data:              make(map[string]interface{}),
	}
}

func (c *Claims) Set(key string, value interface{}) {
	c.Data[key] = value
}

func (c *Claims) Get(key string) interface{} {
	return c.Data[key]
}

func (c *Claims) IsSystemAdmin() bool {
	return c.Role == "system_admin"
}

func (c *Claims) IsSystemUser() bool {
	return c.Role == "system_user"
}

func (c *Claims) IsSystemToken() bool {
	return c.Role == "system_token"
}

func (c *Claims) IsRootUser() bool {
	return c.Role == "root_user"
}

func (c *Claims) IsWorkspaceUser() bool {
	return c.Role == "workspace_user"
}

func (c *Claims) IsAPIToken() bool {
	return c.Role == "api_token"
}

func NewJWT(claims *Claims, privKeyPem string, ttl time.Duration) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyPem))
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	mapClaims := make(jwt.MapClaims)
	mapClaims["id"] = claims.ID
	mapClaims["username"] = claims.Username
	mapClaims["role"] = claims.Role
	mapClaims["allowed_systems"] = claims.AllowedSystems
	mapClaims["allowed_workspaces"] = claims.AllowedWorkspaces
	mapClaims["dat"] = claims.Data
	mapClaims["iat"] = now.Unix()
	mapClaims["exp"] = now.Add(ttl).Unix()
	mapClaims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign: %w", err)
	}

	return token, nil
}

func VerifyJWTWithPublicKey(tokenStr string, pubKeyPem string) (*Claims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKeyPem))
	if err != nil {
		return nil, fmt.Errorf("parse: parse key: %w", err)
	}

	token, err := verifyToken(tokenStr, key)
	if err != nil {
		return nil, fmt.Errorf("parse: verify token: %w", err)
	}

	return jwtTokenToClaims(token)
}

func VerifyJWTWithPrivateKey(tokenStr string, privKeyPem string) (*Claims, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyPem))
	if err != nil {
		return nil, fmt.Errorf("verify: parse key: %w", err)
	}

	token, err := verifyToken(tokenStr, key.Public())
	if err != nil {
		return nil, fmt.Errorf("verify: verify token: %w", err)
	}

	return jwtTokenToClaims(token)
}

func verifyToken(tokenStr string, pubKey interface{}) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func jwtTokenToClaims(token *jwt.Token) (*Claims, error) {
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("parse: invalid claims")
	}

	var claims = new(Claims)
	claims.ID = uint(mapClaims["id"].(float64))
	claims.Username = mapClaims["username"].(string)
	claims.Role = mapClaims["role"].(string)
	claims.AllowedSystems = make([]string, 0)
	for _, value := range mapClaims["allowed_systems"].([]interface{}) {
		claims.AllowedSystems = append(claims.AllowedSystems, value.(string))
	}
	claims.AllowedWorkspaces = make([]string, 0)
	for _, value := range mapClaims["allowed_workspaces"].([]interface{}) {
		claims.AllowedWorkspaces = append(claims.AllowedWorkspaces, value.(string))
	}
	claims.Data = mapClaims["dat"].(map[string]interface{})

	return claims, nil
}
