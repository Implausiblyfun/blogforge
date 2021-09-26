package apid

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/implausiblyfun/blogforge/internal/user"
	"github.com/jmoiron/sqlx"
)

// NewUserHandler for creation of new users
func NewUserHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := &user.User{}
		err := c.BindJSON(u)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Must include json object with the following fields filled in username, first_namem, last_name, password"})
			return
		}
		err = u.Store(c.Request.Context(), db, true)
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "created" + u.Username})
	}
}
