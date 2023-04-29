package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginger-go/micro"
	"github.com/ginger-go/micro/plugins/jwt"
)

// Only allow the auth service to access this api
func AuthServiceOnly(ctx *gin.Context) {
	if ctx.ClientIP() != AUTH_SERVICE_IP {
		abortUnauthorized(ctx)
		return
	}
	ctx.Next()
}

// Only allow to access with system token, user token or api token
func LoginRequired(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		abortUnauthorized(ctx)
		return
	}

	var claims *jwt.Claims
	var err error
	// try to parse token with user token public key
	claims, err = jwt.ParseWithPublicKey(token, USER_TOKEN_PUBLIC_PEM)
	if err != nil || claims == nil {
		// try to parse token with system token public key
		claims, err = jwt.ParseWithPublicKey(token, SYSTEM_TOKEN_PUBLIC_PEM)
		if err != nil || claims == nil {
			abortUnauthorized(ctx)
			return
		}
	}

	apiUUID := GetApiUUID(ctx)
	if apiUUID == "" {
		abortForbidden(ctx)
		return
	}

	// for refresh token, it is never allowed to access any api
	if claims.TokenType == jwt.TOKEN_TYPE_REFRESH_TOKEN {
		abortUnauthorized(ctx)
		return
	}

	// for system token, if the request is from the same ip as the token, then pass
	if claims.TokenType == jwt.TOKEN_TYPE_SYSTEM_TOKEN {
		// the system token must be restricted to a specific ip
		if ctx.ClientIP() != claims.IP {
			abortUnauthorized(ctx)
			return
		} else {
			ctx.Next()
			return
		}
	}

	subscriptionUUID := checkUserIsAllowed(claims.UUID, GetApiUUID(ctx))
	if subscriptionUUID == "" {
		abortForbidden(ctx)
		return
	}

	// for api token
	if claims.TokenType == jwt.TOKEN_TYPE_API_TOKEN {
		// if the api token has been restricted to a specific ip, then check the ip
		if claims.IP != "" && claims.IP != ctx.ClientIP() {
			abortUnauthorized(ctx)
			return
		}
		if claims.HasAPIRight(SYSTEM_ID, apiUUID) {
			_, ok := SUBSCRIPTION_USAGE_MAP[subscriptionUUID]
			if !ok {
				SUBSCRIPTION_USAGE_MAP[subscriptionUUID] = 1
			} else {
				SUBSCRIPTION_USAGE_MAP[subscriptionUUID] += 1
			}
			ctx.Next()
			return
		}
		abortForbidden(ctx)
		return
	}

	// for access token
	if claims.TokenType == jwt.TOKEN_TYPE_ACCESS_TOKEN {
		// the access token must be restricted to a specific ip
		// it is supposed to refresh the access token if the ip is changed
		if claims.IP != ctx.ClientIP() {
			abortUnauthorized(ctx)
			return
		}
		if claims.HasAPIRight(SYSTEM_ID, apiUUID) {
			_, ok := SUBSCRIPTION_USAGE_MAP[subscriptionUUID]
			if !ok {
				SUBSCRIPTION_USAGE_MAP[subscriptionUUID] = 1
			} else {
				SUBSCRIPTION_USAGE_MAP[subscriptionUUID] += 1
			}
			ctx.Next()
			return
		}
		abortForbidden(ctx)
		return
	}

	// unknown token type, should not happen
	abortUnauthorized(ctx)
}

// Only allow to access with refresh token
func RefreshTokenOnly(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		abortUnauthorized(ctx)
		return
	}

	claims, err := jwt.ParseWithPublicKey(token, USER_TOKEN_PUBLIC_PEM)
	if err != nil || claims == nil {
		abortUnauthorized(ctx)
		return
	}

	if claims.TokenType != jwt.TOKEN_TYPE_REFRESH_TOKEN {
		abortUnauthorized(ctx)
		return
	}
}

func abortUnauthorized(ctx *gin.Context) {
	traceID := micro.GetTraceID(ctx)
	traces := micro.GetTraces(ctx)
	traces = append(traces, micro.Trace{
		Success:    false,
		Time:       time.Now(),
		SystemID:   SYSTEM_ID,
		SystemName: SYSTEM_NAME,
		TraceID:    traceID,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
	})
	ctx.AbortWithStatusJSON(401, micro.Response{
		Success: false,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
		TraceID: traceID,
		Traces:  traces,
	})
}

func abortForbidden(ctx *gin.Context) {
	traceID := micro.GetTraceID(ctx)
	traces := micro.GetTraces(ctx)
	traces = append(traces, micro.Trace{
		Success:    false,
		Time:       time.Now(),
		SystemID:   SYSTEM_ID,
		SystemName: SYSTEM_NAME,
		TraceID:    traceID,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_FORBIDDEN,
			Message: ERR_MSG_FORBIDDEN,
		},
	})
	ctx.AbortWithStatusJSON(403, micro.Response{
		Success: false,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_FORBIDDEN,
			Message: ERR_MSG_FORBIDDEN,
		},
		TraceID: traceID,
		Traces:  traces,
	})
}
