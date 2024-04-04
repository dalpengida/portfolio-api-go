package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/dalpengida/portfolio-api-go-mysql/database/my"
)

var ctx = context.Background()

func init() {
	err := my.SetBoilerDatabas()
	if err != nil {
		panic(err)
	}
}

// Test_ExistUsername 는 이미 username 을 사용하고 있는지 검증을 하기 위한 함수
func Test_ExistUsername(t *testing.T) {
	accoutService := NewAccountService()
	exist, err := accoutService.IsExistUsername(ctx, "dalpengida")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(exist)
}

// Test_RegAccount 는 유저 정보를 등록하는 것으 검증 하기 위한 함수
func Test_RegAccount(t *testing.T) {
	accoutService := NewAccountService()
	_, err := accoutService.Register(ctx, "google", "google-uuid-xxxx", "test-username")
	if err != nil {
		t.Fatal(err)
	}
}
