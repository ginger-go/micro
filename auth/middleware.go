package auth

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginger-go/micro"
	"github.com/google/uuid"
)

const (
	ERR_CODE_UNAUTHORIZED = "b97cf20d-42b6-470e-9e08-b4bb852c3811"
	ERR_CODE_FORBIDDEN    = "7792176d-0196-4a57-a959-93062c2b9b41"

	ERR_MSG_UNAUTHORIZED = "Unauthorized"
	ERR_MSG_FORBIDDEN    = "Forbidden"
)

// Only allow system admin or system user to access
// system user should contain system id in the token
// root user & workspace user should not access
func SystemManagerOnly(engine *micro.Engine, authHost string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID, traces := getTracePair(ctx)

		claims := getAuthClaims(ctx, authHost)
		if claims == nil {
			abortUnauthorized(engine, ctx, traceID, traces)
			return
		}

		if claims.IsSystemAdmin() {
			ctx.Next()
			return
		}

		if claims.IsSystemUser() || claims.IsSystemToken() {
			for _, systemID := range claims.AllowedSystems {
				if systemID == engine.SystemID {
					ctx.Next()
					return
				}
			}
			abortForbidden(engine, ctx, traceID, traces)
			return
		}

		abortUnauthorized(engine, ctx, traceID, traces)
	}
}

// Allow root user or workspace user to access
// workspace user should contain workspace id in the token
// workspace user should contain system id in the token
// system admin & system user should not access to client apis
func LoginRequired(engine *micro.Engine, authHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID, traces := getTracePair(c)

		claims := getAuthClaims(c, authHost)
		if claims == nil {
			abortUnauthorized(engine, c, traceID, traces)
			return
		}

		if claims.IsRootUser() {
			c.Next()
			return
		}

		if claims.IsWorkspaceUser() || claims.IsAPIToken() {
			var hasWorkspace bool
			var hasSystem bool
			workspaceID := micro.GetWorkspaceID(c)
			if workspaceID != "" {
				for _, workspace := range claims.AllowedWorkspaces {
					if workspace == workspaceID {
						hasWorkspace = true
						break
					}
				}
			} else {
				hasWorkspace = true // no workspace id in the token, may not required
			}
			for _, system := range claims.AllowedSystems {
				if system == engine.SystemID {
					hasSystem = true
					break
				}
			}
			if hasWorkspace && hasSystem {
				c.Next()
				return
			}
			abortForbidden(engine, c, traceID, traces)
			return
		}

		abortUnauthorized(engine, c, traceID, traces)
	}
}

func abortUnauthorized(engine *micro.Engine, ctx *gin.Context, traceID string, traces []micro.Trace) {
	traces = append(traces, micro.Trace{
		TraceID:    traceID,
		Success:    false,
		Time:       time.Now(),
		SystemID:   engine.SystemID,
		SystemName: engine.SystemName,
		Error: micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
	})
	response := micro.Response{
		Success: false,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_UNAUTHORIZED,
			Message: ERR_MSG_UNAUTHORIZED,
		},
		Traces: traces,
	}
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func abortForbidden(engine *micro.Engine, ctx *gin.Context, traceID string, traces []micro.Trace) {
	traces = append(traces, micro.Trace{
		TraceID:    traceID,
		Success:    false,
		Time:       time.Now(),
		SystemID:   engine.SystemID,
		SystemName: engine.SystemName,
		Error: micro.ResponseError{
			Code:    ERR_CODE_FORBIDDEN,
			Message: ERR_MSG_FORBIDDEN,
		},
	})
	response := micro.Response{
		Success: false,
		Error: &micro.ResponseError{
			Code:    ERR_CODE_FORBIDDEN,
			Message: ERR_MSG_FORBIDDEN,
		},
		Traces: traces,
	}
	ctx.AbortWithStatusJSON(http.StatusForbidden, response)
}

func getAuthClaims(c *gin.Context, authHost string) *Claims {
	token := GetJWTTokenFromHeader(c)
	if token == "" {
		micro.Logger.Info("auth: no jwt token found in header (ip: ", c.ClientIP(), ")")
		return nil
	}

	publicPEM := getPublicKey(authHost)
	claims, err := VerifyJWTWithPublicKey(token, publicPEM)
	if err != nil {
		micro.Logger.Error("auth: failed to verify jwt token: ", err)
		return nil
	}

	return claims
}

var publicKeyPEM = ""
var publicKeyPEMExpired = time.Now()

func getPublicKey(host string) string {
	if publicKeyPEM != "" && time.Now().Before(publicKeyPEMExpired) {
		return publicKeyPEM
	}

	url := host + "/public_key"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		micro.Logger.Error("auth: failed to get public key from auth service: ", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		micro.Logger.Error("auth: failed to get public key from auth service: ", err)
		return ""
	}

	publicKeyPEM = string(body)
	publicKeyPEMExpired = time.Now().Add(time.Hour)
	return publicKeyPEM
}

// get trace id and traces from context
func getTracePair(c *gin.Context) (string, []micro.Trace) {
	traceID := micro.GetTraceID(c)
	traces := micro.GetTraces(c)

	if traceID == "" {
		traceID = uuid.NewString()
		micro.SetTraceID(c, traceID)
	}

	if len(traces) == 0 {
		traces = make([]micro.Trace, 0)
		micro.SetTraces(c, traces)
	}

	return traceID, traces
}
