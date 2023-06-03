package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/mohaali/goAuth/auth"
)

func Handlers(s auth.UserService) *gin.Engine {
	r := gin.Default()

	return r

}
