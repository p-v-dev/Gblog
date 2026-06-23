package infra

import "github.com/gin-gonic/gin"

func CurrentUserID(c *gin.Context) string {
	id, _ := c.Get("user_id")
	return id.(string)
}
