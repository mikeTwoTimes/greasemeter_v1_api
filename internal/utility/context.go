package utility

import (
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) int {
	id, exist := c.Get("userId")

	if !exist {
		return 0
	}

	userId, ok := id.(int)

	if !ok {
		return 0
	}

	return userId
}
