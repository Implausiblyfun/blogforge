package apid

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/implausiblyfun/blogforge/internal/standardroutes"
	"github.com/implausiblyfun/blogforge/internal/user"
	"github.com/jmoiron/sqlx"
)

// WithFire is nice way to ask the server to stop... withfire.
func WithFire(cancel func()) gin.HandlerFunc {
	return func(c *gin.Context) {
		uID := "just nate"
		fmt.Printf("User: %s just nicely asked us to stop the server\n", uID)
		standardroutes.EmptyGoodHandler()(c)
		cancel()

	}
}

// WithFireAndBrimstone just panics... because well its scary.
func WithFireAndBrimstone() gin.HandlerFunc {
	return func(c *gin.Context) {

		standardroutes.EmptyGoodHandler()(c)
		panic("BRIMSTONE")
	}
}

// ListUsers is a intermidate state for db checks and or in future for admin listing.
func ListUsers(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// return all users
		usrs, err := user.Bloggers(c.Request.Context(), "*", db)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, standardroutes.NotFoundBody())
			return
		}

		c.JSON(http.StatusOK, usrs)

	}
}

// AdminConfig is for our testing and sanity checks and would be removed.
func AdminConfig(cfg map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("%v", cfg))
	}

}
