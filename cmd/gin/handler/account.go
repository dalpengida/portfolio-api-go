package handler

import (
	"fmt"
	"net/http"

	"github.com/dalpengida/portfolio-api-go-mysql/common"
	"github.com/dalpengida/portfolio-api-go-mysql/common/auth"
	"github.com/dalpengida/portfolio-api-go-mysql/common/redis"
	"github.com/dalpengida/portfolio-api-go-mysql/service"
	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

// gin 에서 context 를 넘겨 줄 때
// https://github.com/gin-gonic/gin/issues/2845

type Account struct{}

func (Account) Signup(c *gin.Context) {
	type signupRequest struct {
		ProviderId string `json:"provider_id" validate:"required"`
		Provider   string `json:"provider" validate:"required"`
		Username   string `json:"username" validate:"required"`
	}
	type signupResponse struct {
		common.Response
		Token string `json:"token"`
	}

	req := new(signupRequest)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request, %w", err))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request parameter, %w", err))
		return
	}

	accountServer := service.NewAccountService()
	userId, err := accountServer.Register(c.Request.Context(), req.Provider, req.ProviderId, req.Username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("account register failed, %w", err))
		return
	}

	// token 발행
	// TODO: user_id 값을 string 으로 할지 int64로 할지 정해야 함
	token, err := auth.NewJwtToken(string(rune(userId)))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// token cache
	err = redis.SetEX(c.Request.Context(), redis.SessionKey(string(rune(userId))), token, redis.SESSION_EXPIRE)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res := new(signupResponse)
	res.Token = token

	c.JSON(http.StatusOK, res)
}

func (Account) Login(c *gin.Context) {
	type loginRequest struct {
		ProviderId string `json:"provider_id" validate:"required"`
		Provider   string `json:"provider" validate:"required"`
	}
	type loginResponse struct {
		common.Response
		Token string `json:"token"`
	}

	req := new(loginRequest)

	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request, %w", err))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid request parameter, %w", err))
		return
	}

	accountService := service.NewAccountService()
	userId, err := accountService.GetUserIdFromIDP(c.Request.Context(), req.Provider, req.ProviderId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("user_id find failed, %w", err))
		return
	}

	token, err := auth.NewJwtToken(string(rune(userId)))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = redis.SetEX(c.Request.Context(), redis.SessionKey(string(rune(userId))), token, redis.SESSION_EXPIRE)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	res := new(loginResponse)
	res.Token = token

	c.JSON(http.StatusOK, res)
}

func (Account) Logout(c *gin.Context) {
	userId := c.GetString("user_id")

	res := new(common.Response)

	_, err := redis.Del(c.Request.Context(), redis.SessionKey(userId))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
