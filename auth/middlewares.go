package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	verifier "github.com/okta/okta-jwt-verifier-golang"
	"os"
	"strings"
)

func isAuthenticated(c *gin.Context) bool {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}
	tokenParts := strings.Split(authHeader, "Bearer ")
	bearerToken := tokenParts[1]

	tv := map[string]string{}
	tv["aud"] = "api://default"
	tv["cid"] = os.Getenv("SPA_CLIENT_ID")
	jv := verifier.JwtVerifier{
		Issuer:           os.Getenv("ISSUER"),
		ClaimsToValidate: tv,
	}

	myJwt, err := jv.New().VerifyAccessToken(bearerToken)
	if err != nil {
		return false
	}
	c.Set("user_email", myJwt.Claims["sub"])
	return true
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c) {
			err := errors.New("auth error")
			c.AbortWithError(401, err)
		}
	}
}
