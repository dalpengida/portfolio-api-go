package service

import (
	"context"
	"fmt"

	"github.com/dalpengida/portfolio-api-go-mysql/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Account struct{}

func NewAccountService() Account {
	return Account{}
}

func (Account) IsExistUsername(ctx context.Context, username string) (bool, error) {
	//db := boil.GetContextDB()
	count, err := models.Usernames(
		qm.Where("username = ?", username),
	).CountG(ctx)
	if err != nil {
		return false, err
	}

	if count != 0 {
		return true, nil
	}

	return false, nil
}

// Register 는 idp 정보 및 account 정보를 등록
func (Account) Register(ctx context.Context, provider, providerId, username string) (int64, error) {
	tx, err := boil.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("transaction begin faeild, %w", err)
	}

	userId := int64(uuid.New().ID())

	idp := models.Idp{
		Provider:   provider,
		ProviderID: providerId,
		UserID:     userId,
	}
	err = idp.Insert(ctx, tx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err) // rollback 오류가 났을 경우 어떻게 할 수 있는 게 없어서, 로그만 일단 남김
		}

		return 0, fmt.Errorf("idp insert failed, %w", err)
	}

	account := models.Account{
		UserID:   userId,
		Username: username,
	}
	err = account.Insert(ctx, tx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Error().Err(err) // rollback 오류가 났을 경우 어떻게 할 수 있는 게 없어서, 로그만 일단 남김
		}

		return 0, fmt.Errorf("account insert failed, %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("transaction commit error , %w", err)
	}

	return userId, nil
}

// GetUserIdFromIDP 는 idp 정보를 가지고 userid 값을 조회
func (Account) GetUserIdFromIDP(ctx context.Context, provider, providerId string) (int64, error) {
	idp, err := models.Idps(
		qm.Select("user_id"),
		qm.Where("provider=? AND provider_id = ?", provider, providerId),
	).OneG(ctx)
	if err != nil {
		return 0, fmt.Errorf("find provider info failed, %w", err)
	}

	return idp.UserID, nil
}
