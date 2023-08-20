package usecase

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase struct {
	storage storage
	logger  *zap.Logger
	key     string
}

type storage interface {
}

func New(st storage, logger *zap.Logger) *UseCase {
	return &UseCase{storage: st, logger: logger}
}

func (u *UseCase) GetTokens(guid string) (access string, refresh string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodES256,
		jwt.MapClaims{
			"guid": guid,
		})

	access, err = t.SignedString(u.key)
	if err != nil {
		u.logger.Error("can't sign token", zap.Error(err))
		return "", "", ErrSign
	}

	uuidObj, err := uuid.NewUUID()
	if err != nil {
		u.logger.Error("can't generate uuid", zap.Error(err))
		return "", "", ErrGenerateUUID
	}

	return access, uuidObj.String(), nil
}

func (u *UseCase) RefreshTokens(oldRefresh string) (access string, refresh string, err error) {
	return "", "", nil
}
