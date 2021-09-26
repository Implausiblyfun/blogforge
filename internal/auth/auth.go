// Package auth currently only
package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/implausiblyfun/blogforge/internal/user"
	"github.com/jmoiron/sqlx"
)

func BasicAuther(db *sqlx.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		u, pass, err := decomposeBasic(c.Request.Header.Get("Authorization"))
		fmt.Printf("Attempting to check basic authing %s, %s, %v\n", u, pass, err)
		if err != nil {
			c.Header("WWW-Authenticate", "base")
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		usrs, err := user.SpecificBlogger(c.Request.Context(), u, db)
		if len(usrs) == 1 {
			usr := usrs[0]
			if ok := usr.ValidatePassword(pass); ok {
				c.Set("user", usr)
				c.Next()
				return
			}
		}

		c.Header("WWW-Authenticate", "base")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error"})
		c.AbortWithStatus(http.StatusUnauthorized)

	}
}

func decomposeBasic(header string) (string, string, error) {
	toDecode := strings.TrimPrefix(header, "Basic ")
	if len(toDecode) == len(header) {
		return "", "", fmt.Errorf("Auth type was not Basic")
	}
	baseBytes, err := base64.StdEncoding.DecodeString(toDecode)
	if err != nil {
		return "", "", fmt.Errorf("No auth user stored")
	}
	baseStr := string(baseBytes)
	split := strings.SplitN(baseStr, ":", 2)
	if len(split) != 2 {
		return split[0], "", fmt.Errorf("No Password specified")
	}
	return split[0], split[1], nil
}
