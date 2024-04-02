package auth

import (
	"time"

	"github.com/dalpengida/portfolio-api-go/config"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: log 관련 구조체를 정리를 해야 함, 현재는 마구잡이

var (
	VERSION             string        = "1.0.0"
	TOKEN_EXPIRE_PERIOD time.Duration = time.Hour * 72
)

type JWTClaims struct {
	UserId  string `json:"user_id"`
	Version string `json:"version"`
	Expire  int64  `json:"expire"`

	jwt.RegisteredClaims
}

// NewJwtToken 는 jwt 토큰을 만들고 signing 후 전달
func NewJwtToken(userId string) (string, error) {
	signingKey := []byte(config.Config(config.JWT_SIGNING_KEY)) // 나중에 변수 처리를 해야 함

	claims := JWTClaims{
		UserId:  userId,
		Version: VERSION,
		Expire:  time.Now().Add(TOKEN_EXPIRE_PERIOD).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signingToken, nil
}

// // GㅁetClaims 는 claim 정보를 한번에 파싱을 해서 전달을 하기 위하여 무식하게 json marshaling 했다가 unmarshaling 을 함
// // jwt v5 로 넘어 오면서, map 값으로 전달을 함
// // 속도 이슈로 아마 해당 함수를 안 쓰게 되리 것 같음, bench 확인 필요
// func GetClaims(c *fiber.Ctx) (*JWTClaims, bool) {
// 	u := c.Locals("user")
// 	if u == nil {
// 		return nil, false
// 	}
// 	t, ok := u.(*jwt.Token)
// 	if !ok {
// 		return nil, false
// 	}

// 	// map 으로 되어 있는 데이터를 struct 에 넣어 주기 위하여 json marshaling 후 unmarshaling 을 진행
// 	j, err := json.Marshal(t.Claims)
// 	if err != nil {
// 		log.Error().Err(err).Msg("json marshal failed")
// 		return nil, false
// 	}

// 	claims := new(JWTClaims)
// 	err = json.Unmarshal(j, &claims)
// 	if err != nil {
// 		log.Err(err).Msg("json unmarshal failed")
// 	}

// 	return claims, true
// }
