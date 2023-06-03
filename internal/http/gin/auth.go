package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/goAuth/auth"
	"github.com/mohaali482/goAuth/errors"
)

func Handlers(s auth.UserService) *gin.Engine {
	r := gin.Default()
	r.Handle("POST", "/users", Create(s))

	return r

}

func Create(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user auth.User
		err := c.ShouldBind(&user)
		if err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		err = user.Validate()
		if err != nil {
			errors.ReturnErrorResponse(err, c)
			return
		}

		user, err = s.Create(user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}

}
