package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ginger-go/micro"
	"github.com/ginger-go/micro/plugins/jwt"
)

// Only allow the auth service to access this api
func AuthServiceOnly(ctx *gin.Context) {
	traceID := micro.GetTraceID(ctx)
	traces := micro.GetTraces(ctx)
	if ctx.ClientIP() != AUTH_SERVICE_IP {
		ctx.AbortWithStatusJSON(401, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_UNAUTHORIZED,
				Message: ERR_MSG_UNAUTHORIZED,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}
	ctx.Next()
}

// Only allow to access with system token, user token or api token
func LoginRequired(ctx *gin.Context) {
	traceID := micro.GetTraceID(ctx)
	traces := micro.GetTraces(ctx)
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.AbortWithStatusJSON(401, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_UNAUTHORIZED,
				Message: ERR_MSG_UNAUTHORIZED,
			},
			TraceID: traceID,
			Traces:  traces,
		})
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
			ctx.AbortWithStatusJSON(401, micro.Response{
				Success: false,
				Error: &micro.ResponseError{
					Code:    ERR_CODE_UNAUTHORIZED,
					Message: ERR_MSG_UNAUTHORIZED,
				},
				TraceID: traceID,
				Traces:  traces,
			})
			return
		}
	}

	requestTag := ctx.Request.Method + ":" + ctx.Request.URL.Path
	apiUUID, ok := API_UUID_MAP[requestTag]
	if apiUUID == "" || !ok {
		ctx.AbortWithStatusJSON(403, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_FORBIDDEN,
				Message: ERR_MSG_FORBIDDEN,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}

	// for refresh token, it is never allowed to access any api
	if claims.TokenType == jwt.TOKEN_TYPE_REFRESH_TOKEN {
		ctx.AbortWithStatusJSON(401, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_UNAUTHORIZED,
				Message: ERR_MSG_UNAUTHORIZED,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}

	// for system token, if the request is from the same ip as the token, then pass
	if claims.TokenType == jwt.TOKEN_TYPE_SYSTEM_TOKEN {
		// the system token must be restricted to a specific ip
		if ctx.ClientIP() != claims.IP {
			ctx.AbortWithStatusJSON(401, micro.Response{
				Success: false,
				Error: &micro.ResponseError{
					Code:    ERR_CODE_UNAUTHORIZED,
					Message: ERR_MSG_UNAUTHORIZED,
				},
				TraceID: traceID,
				Traces:  traces,
			})
			return
		} else {
			ctx.Next()
			return
		}
	}

	// for api token
	if claims.TokenType == jwt.TOKEN_TYPE_API_TOKEN {
		// if the api token has been restricted to a specific ip, then check the ip
		if claims.IP != "" && claims.IP != ctx.ClientIP() {
			ctx.AbortWithStatusJSON(401, micro.Response{
				Success: false,
				Error: &micro.ResponseError{
					Code:    ERR_CODE_UNAUTHORIZED,
					Message: ERR_MSG_UNAUTHORIZED,
				},
				TraceID: traceID,
				Traces:  traces,
			})
			return
		}
		if claims.HasAPIRight(SYSTEM_ID, apiUUID) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(403, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_FORBIDDEN,
				Message: ERR_MSG_FORBIDDEN,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}

	// for access token
	if claims.TokenType == jwt.TOKEN_TYPE_ACCESS_TOKEN {
		// the access token must be restricted to a specific ip
		// it is supposed to refresh the access token if the ip is changed
		if claims.IP != ctx.ClientIP() {
			ctx.AbortWithStatusJSON(401, micro.Response{
				Success: false,
				Error: &micro.ResponseError{
					Code:    ERR_CODE_UNAUTHORIZED,
					Message: ERR_MSG_UNAUTHORIZED,
				},
				TraceID: traceID,
				Traces:  traces,
			})
			return
		}
		if claims.HasAPIRight(SYSTEM_ID, apiUUID) {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(403, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_FORBIDDEN,
				Message: ERR_MSG_FORBIDDEN,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}

	// unknown token type, should not happen
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

// Only allow to access with refresh token
func RefreshTokenOnly(ctx *gin.Context) {
	traceID := micro.GetTraceID(ctx)
	traces := micro.GetTraces(ctx)
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.AbortWithStatusJSON(401, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_UNAUTHORIZED,
				Message: ERR_MSG_UNAUTHORIZED,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}

	claims, err := jwt.ParseWithPublicKey(token, USER_TOKEN_PUBLIC_PEM)
	if err != nil || claims == nil {
		ctx.AbortWithStatusJSON(401, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_UNAUTHORIZED,
				Message: ERR_MSG_UNAUTHORIZED,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}

	if claims.TokenType != jwt.TOKEN_TYPE_REFRESH_TOKEN {
		ctx.AbortWithStatusJSON(401, micro.Response{
			Success: false,
			Error: &micro.ResponseError{
				Code:    ERR_CODE_UNAUTHORIZED,
				Message: ERR_MSG_UNAUTHORIZED,
			},
			TraceID: traceID,
			Traces:  traces,
		})
		return
	}
}
