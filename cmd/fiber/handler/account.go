package handler

import (
	"fmt"
	"net/http"

	"github.com/dalpengida/portfolio-api-go-mysql/common"
	"github.com/dalpengida/portfolio-api-go-mysql/common/auth"
	"github.com/dalpengida/portfolio-api-go-mysql/common/redis"
	"github.com/dalpengida/portfolio-api-go-mysql/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Account struct{}

func (Account) Signup(c *fiber.Ctx) error {
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

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Errorf("invalid request").Error())
	}
	if err := validator.New().Struct(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request parameter", err.Error())
	}

	accountServer := service.NewAccountService()
	userId, err := accountServer.Register(c.Context(), req.Provider, req.ProviderId, req.Username)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "account register failed", err.Error())
	}

	// token 발행
	// TODO: user_id 값을 string 으로 할지 int64로 할지 정해야 함
	token, err := auth.NewJwtToken(string(rune(userId)))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	// token cache
	err = redis.SetEX(c.Context(), redis.SessionKey(string(rune(userId))), token, redis.SESSION_EXPIRE)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	res := new(signupResponse)
	res.Token = token

	return c.JSON(res)
}

func (Account) Login(c *fiber.Ctx) error {
	type loginRequest struct {
		ProviderId string `json:"provider_id" validate:"required"`
		Provider   string `json:"provider" validate:"required"`
	}
	type loginResponse struct {
		common.Response
		Token string `json:"token"`
	}

	req := new(loginRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, fmt.Errorf("invalid request").Error())
	}
	if err := validator.New().Struct(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request parameter", err.Error())
	}

	accountService := service.NewAccountService()
	userId, err := accountService.GetUserIdFromIDP(c.Context(), req.Provider, req.ProviderId)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "user_id find failed", err.Error())
	}

	token, err := auth.NewJwtToken(string(rune(userId)))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	err = redis.SetEX(c.Context(), redis.SessionKey(string(rune(userId))), token, redis.SESSION_EXPIRE)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	res := new(loginResponse)
	res.Token = token

	return c.JSON(res)
}

func (Account) Logout(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	res := new(common.Response)

	_, err := redis.Del(c.Context(), redis.SessionKey(userId))
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(res)
}
