package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ashishkhuraishy/go_forum/utils/token"
	"github.com/gin-gonic/gin"
)

const (
	authHeader     = "authorization"
	authTypeBearer = "bearer"
	authPayLoadKey = "auth+payload_key"
)

func authMiddleWare(token token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authHeader)
		if len(authHeader) == 0 {
			err := fmt.Errorf("request doesn't contain any auth headers")
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := fmt.Errorf("invalid auth header format")
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != authTypeBearer {
			err := fmt.Errorf("unsupported auth type")
			c.AbortWithStatusJSON(http.StatusBadRequest, err)
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		c.Set(authPayLoadKey, payload)
		c.Next()
	}
}
