package gin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mohaali482/goAuth/auth"
	"github.com/mohaali482/goAuth/internal/http/gin/errors"
	"github.com/mohaali482/goAuth/internal/http/gin/middlewares"
)

func Handlers(s auth.UserService) *gin.Engine {
	r := gin.Default()
	r.Handle("POST", "/accounts/login", Login(s))
	r.Handle("DELETE", "/accounts/logout", Logout(s))
	r.Handle("POST", "/accounts/signup", Signup(s))
	usersGroup := r.Group("/users").Use(middlewares.AuthMiddleware(s))
	{
		usersGroup.Handle("POST", "", Create(s))
		usersGroup.Handle("GET", "", GetAll(s))
		usersGroup.Handle("GET", ":id", GetByID(s))
		usersGroup.Handle("DELETE", ":id", Delete(s))
		usersGroup.Handle("PATCH", ":id", Update(s))
	}

	return r

}

func Create(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user auth.User
		err := c.ShouldBindJSON(&user)
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

func Delete(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is not a valid id"})
		}
		if err := s.Delete(id); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func GetAll(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.GetAll()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetByID(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is not a valid id"})
			return
		}
		user, err := s.GetByID(id)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func Login(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin auth.UserLogin
		err := c.ShouldBindJSON(&userLogin)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}
		user, err := s.Login(userLogin.Username, userLogin.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": auth.ErrWrongCredentials.Error()})
			return
		}
		tokens, err := s.GenerateJWT(user)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		c.SetCookie("refresh_token", tokens["refresh"], 3600, "/", "localhost", false, true)
		c.SetCookie("access_token", tokens["access"], 3600, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func Logout(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("refresh_token", "", 0, "/", "localhost", false, true)
		c.SetCookie("access_token", "", 0, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func Signup(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userForm auth.UserForm
		err := c.ShouldBindJSON(&userForm)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}
		err = userForm.Validate()
		if err != nil {
			errors.ReturnErrorResponse(err, c)
			return
		}

		user, err := s.Create(userForm.ToUserEntity())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}

}

func Update(s auth.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userForm auth.UserForm
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is not a valid id"})
		}
		err = c.ShouldBindJSON(&userForm)
		if err != nil {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		if userForm.Phone != "" {
			err = userForm.ValidatePhone()
			if err != nil {
				errors.ReturnErrorResponse(err, c)
				return
			}
		}

		user, err := s.Update(id, userForm.ToUserEntity())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
