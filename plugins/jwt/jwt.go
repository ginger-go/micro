package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Issue issues a JWT token with the given claims, private key and ttl.
func Issue(claims *Claims, privKeyPem string, ttl time.Duration) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyPem))
	if err != nil {
		return "", fmt.Errorf("issue jwt: parse key pem: %w", err)
	}

	now := time.Now().UTC()
	mapClaims := make(jwt.MapClaims)
	mapClaims["uuid"] = claims.UUID
	mapClaims["name"] = claims.Name
	mapClaims["ip"] = claims.IP
	mapClaims["token_type"] = claims.TokenType
	mapClaims["apis"] = claims.APIs
	mapClaims["data_sets"] = claims.DataSets
	mapClaims["data"] = claims.Data
	mapClaims["iat"] = now.Unix()
	mapClaims["exp"] = now.Add(ttl).Unix()
	mapClaims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("issue jwt: sign: %w", err)
	}

	return token, nil
}

// ParseWithPublicKey parses a JWT token with the given claims and public key.
func ParseWithPublicKey(tokenStr string, pubKeyPem string) (*Claims, error) {
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

// ParseWithPrivateKey parses a JWT token with the given claims and private key.
func ParseWithPrivateKey(tokenStr string, privKeyPem string) (*Claims, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyPem))
	if err != nil {
		return nil, fmt.Errorf("parse: parse key: %w", err)
	}

	token, err := verifyToken(tokenStr, key.Public())
	if err != nil {
		return nil, fmt.Errorf("parse: verify token: %w", err)
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
	claims.UUID = mapClaims["uuid"].(string)
	claims.Name = mapClaims["name"].(string)
	claims.IP = mapClaims["ip"].(string)
	claims.TokenType = mapClaims["token_type"].(string)
	claims.APIs = make(map[string]string)
	for _, value := range mapClaims["apis"].(map[string]interface{}) {
		claims.APIs[value.(string)] = value.(string)
	}
	claims.DataSets = make([]string, 0)
	for _, value := range mapClaims["data_sets"].([]interface{}) {
		claims.DataSets = append(claims.DataSets, value.(string))
	}

	return claims, nil
}
