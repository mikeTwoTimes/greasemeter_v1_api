package users

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/types"
	"github.com/mikeTwoTimes/greasemeter_v1_api/internal/utility"
	"github.com/sendgrid/sendgrid-go"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store     types.UserStore
	jwtSecret string
	mailer    mailClient
}

func NewHandler(store types.UserStore, jwtSecret string, client *sendgrid.Client) *Handler {
	return &Handler{
		store:     store,
		jwtSecret: jwtSecret,
		mailer:    mailClient{client: client},
	}
}

func (h *Handler) RegisterRoutes(v1, auth *gin.RouterGroup) {
	v1.POST("/users/register", h.createUser)
	v1.POST("/users/login", h.login)
	v1.POST("/users/forgot-password", h.forgotPassword)
	// v1.POST("/users/reset-password", h.resetPassword)

	auth.DELETE("/users", h.deleteUser)
}

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
			"error": "Error generating token",
		})
	} else {
		c.JSON(http.StatusOK, types.Login{Token: tokenString})
	}
}

func (h *Handler) forgotPassword(c *gin.Context) {
	email, err := utility.ParseForgotPassword(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} 

	userId, err := h.store.GetUserByEmail(email)

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
		"type":   "password_reset",
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error generating token",
		})
		return
	}

	err = h.mailer.SendPasswordReset(tokenString, email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
/*
func (h *Handler) resetPassword(c *gin.Context) {
	claims, err := utility.GetClaimsFromParam(c, a.JWTSecret)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	userId, ok := claims["userId"].(float64)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
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

	_, err = h.store.UpdateUserPassword(int(userId), string(hashedPassword))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusNoContent, nil)
	}
}
*/
func (h *Handler) deleteUser(c *gin.Context) {
	userId := utility.GetUserFromContext(c)

	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
	} else if err := h.store.DeleteUser(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    } else {
        c.JSON(http.StatusNoContent, nil)
    }
}
