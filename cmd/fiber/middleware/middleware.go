package middleware

import (
	"github.com/dalpengida/portfolio-api-go-mysql/config"
	"github.com/golang-jwt/jwt/v5"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// jwtError 는 jwt 에러 처리
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

// RequireJWT jwt 토큰이 있어야만 엑세스 가능한 api 미들웨어
func RequireJWT() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config(config.JWT_SIGNING_KEY))},
		ErrorHandler: jwtError,
	})
}

// RequireJWTWithClaims 는 JWT토큰 검사와 함께 claim 정보를 locals 에 넣어줌, handler 에서 사용을 할 수 있도록 하기 위함
func RequireJWTWithClaims() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config(config.JWT_SIGNING_KEY))},
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			userId := claims["user_id"].(string)

			// userId 값을 추출해서 저장
			c.Locals("user_id", userId)

			log.Info().Msgf("user '%s' accessing to '%s'", userId, c.Request().URI().String())
			return c.Next()
		},

		ErrorHandler: jwtError,
	})
}
