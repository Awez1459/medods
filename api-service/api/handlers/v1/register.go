package v1

import (
	"medods/api-service/api/models"
	pbu "medods/api-service/genproto/user-proto"
	l "medods/api-service/internal/pkg/logger"
	tokens "medods/api-service/internal/pkg/token"
	"net/http"
	"gopkg.in/gomail.v2"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func sendWarningEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "avazbekbekmurodov1459@example.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.example.com", 587, "avazbekbekmurodov1459@example.com", "Awez1459")

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

// RECLAIM TOKENS ...
// @Security BearerAuth
// @Router /v1/users/:id [POST]
// @Summary TOKEN
// @Description Api for tokens of user
// @Tags TOKENS
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.TokenResp
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) Token(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	user, err := h.Service.UserService().Get(c, &pbu.Filter{
		Filter: map[string]string{"id": id},
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		h.Logger.Error("error while get user", l.Error(err))
		return
	}

	clientIP := c.ClientIP()

	h.JwtHandler = tokens.JwtHandler{
		Sub:       user.User.Id,
		Role:      user.User.Role,
		SigninKey: h.Config.Token.SignInKey,
		Log:       h.Logger,
		Timeout:   int(h.Config.Token.AccessTTL),
		Iss:       clientIP,
	}

	access, refresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate tokens",
		})
		h.Logger.Error("error while generate JWT", l.Error(err))
		return
	}

	hashRefresh, err := bcrypt.GenerateFromPassword([]byte(refresh), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash refresh token"})
		h.Logger.Error("error while hash refresh", l.Error(err))
		return
	}

	_, err = h.Service.UserService().Update(c, &pbu.User{
		Id:           user.User.Id,
		RefreshToken: string(hashRefresh),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user with refresh token",
		})
		h.Logger.Error("error while update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.TokenResp{
		Access:  access,
		Refresh: refresh,
	})
}

// UPDATE TOKEN
// @Security BearerAuth
// @Router /v1/token/{refresh} [GET]
// @Summary UPDATE TOKEN
// @Description Api for updated acces token
// @Tags TOKEN
// @Accept json
// @Produce json
// @Param refresh path string true "Refresh Token"
// @Success 200 {object} models.TokenResp
// @Failure 400 {object} models.StandartError
// @Failure 500 {object} models.StandartError
func (h HandlerV1) UpdateToken(c *gin.Context) {
	refresh := c.Param("refresh")
	if refresh == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	user, err := h.Service.UserService().Get(c, &pbu.Filter{
		Filter: map[string]string{"refresh_token": refresh},
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		h.Logger.Error("Failed to get user with refresh token", l.Error(err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.User.RefreshToken), []byte(refresh))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token not found",
		})
		h.Logger.Error("refresh token not found", l.Error(err))
		return
	}

	clientIP := c.ClientIP()
	resClaim, err := tokens.ExtractClaim(refresh, []byte(h.Config.Token.SignInKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Reload Page",
		})
		h.Logger.Error("Failed to extract token update token", l.Error(err))
		return
	}
	if resClaim["iss"] != clientIP {
		h.Logger.Warn("IP address mismatch")
		err := sendWarningEmail(user.User.Email, "IP address mismatch", "Your IP address mismatched.")
		if err != nil {
			h.Logger.Error("Failed to send warning email", l.Error(err))
		}
	}

	h.JwtHandler = tokens.JwtHandler{
		Sub:       user.User.Id,
		Role:      user.User.Role,
		SigninKey: h.Config.Token.SignInKey,
		Log:       h.Logger,
		Timeout:   int(h.Config.Token.AccessTTL),
		Iss:       clientIP,
	}

	newAccess, newRefresh, err := h.JwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new tokens"})
		h.Logger.Error("error while generating new tokens", l.Error(err))
		return
	}

	hashNewR, err := bcrypt.GenerateFromPassword([]byte(newRefresh), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash new refresh token",
		})
		h.Logger.Error("error while hash new refresh token", l.Error(err))
		return
	}

	_, err = h.Service.UserService().Update(c, &pbu.User{
		Id:           user.User.Id,
		RefreshToken: string(hashNewR),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update user with new refresh token",
		})
		h.Logger.Error("error while update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.TokenResp{
		Access:  newAccess,
		Refresh: newRefresh,
	})
}
