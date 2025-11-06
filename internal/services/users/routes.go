package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/utility"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store     types.UserStore
	jwtSecret string
	mailer    *Mailer
}

func NewHandler(store types.UserStore, jwtSecret string, mailer *Mailer) *Handler {
	return &Handler{
		store:     store,
		jwtSecret: jwtSecret,
		mailer:    mailer,
	}
}

func (h *Handler) RegisterRoutes(v1, auth *gin.RouterGroup) {
	v1.POST("/users/register", h.createUser)
	v1.POST("/users/login", h.login)
	v1.POST("/users/forgot-password", h.forgotPassword)
	v1.POST("/users/reset-password/:token", h.resetPassword)

	auth.DELETE("/users", h.deleteUser)
}

// @Summary	    Registers a new user
// @Description	Registers a new user given unique credentials
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body types.RegisterPayload true "User"
// @Success	    204	
// @Router      /v1/users/register [post]
func (h *Handler) createUser(c *gin.Context) {
	req, err := utility.ParseRegister(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	req.Password = string(hashedPassword)

	if err = h.store.CreateUser(req); err != nil {
		c.JSON(utility.MapError(err))
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}

// @Summary	    Signs in a returning user
// @Description	Signs in a user given valid credentials
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body types.LoginPayload true "User"
// @Success	    200 {object} types.Login
// @Router      /v1/users/login [post]
func (h *Handler) login(c *gin.Context) {
	req, err := utility.ParseLogin(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := h.store.GetUserCredentials(req.Name)

	if existingUser.Id == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(existingUser.Password),
		[]byte(req.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
		"type":   "auth",
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
	} else {
		c.JSON(http.StatusOK, types.Login{Token: tokenString})
	}
}

// @Summary	    Sends a password reset email
// @Description	Sends a password reset email to a verified address
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       email body types.ForgotPasswordPayload true "Email"
// @Success	    204
// @Router      /v1/users/forgot-password [post]
func (h *Handler) forgotPassword(c *gin.Context) {
	email, err := utility.ParseForgotPassword(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} 

	userId, err := h.store.GetUserFromEmail(email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	} else if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email not registered",
		})
		return
	}

	tokenString, err := h.store.CreateResetToken(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(h.mailer.SendPasswordReset(tokenString, email))
}

// @Summary	    Resets a user's password
// @Description	Resets a user's password given a valid reset token
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       token path string true "Reset token"
// @Success	    204
// @Router      /v1/users/reset-password/{token} [post]
func (h *Handler) resetPassword(c *gin.Context) {
	data, err := h.store.GetDataFromResetToken(c.Param("token"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if data.UserId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	} else if data.Expiration.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	} 
	
	password, err := utility.ParseResetPassword(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	err = h.store.UpdateUserPassword(data.UserId, string(hashedPassword))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}

// @Summary	    Deletes a user
// @Description	Deletes a user given they have valid access token
// @Tags        users
// @Accept      json
// @Produce     json
// @Success	    204
// @Router      /v1/users [delete]
// @Security    BearerAuth
func (h *Handler) deleteUser(c *gin.Context) {
	userId := c.MustGet("userId").(int)

	if err := h.store.DeleteUser(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
        c.JSON(http.StatusNoContent, nil)
    }
}
