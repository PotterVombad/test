package tokens

import (
	"strings"
	"github.com/PotterVombad/test/internal/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type (
	Repo interface {
		GetAccessToken(uid string) (string, error)
		GetRefreshToken() string
		GetPairs(uid string) (models.TokensPairs, error)
	}

	Tokens struct {
		secretKey string
	}
)

func (t Tokens) GetPairs(uid string) (models.TokensPairs, error) {
	accessToken, err := t.GetAccessToken(uid)
	if err != nil {
		return models.TokensPairs{}, err
	}

	return models.TokensPairs{
		AccessToken:  accessToken,
		RefreshToken: t.GetRefreshToken(),
	}, nil
}

func (t Tokens) GetAccessToken(uid string) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS512)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["uid"] = uid
	claims["aud"] = "test"
	claims["nbf"] = time.Now().UTC()
	claims["exp"] = time.Now().Add(time.Hour * 24).UTC()

	return accessToken.SignedString([]byte(t.secretKey))
}
func (t Tokens) GetRefreshToken() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func New(secretKey string) Tokens {
	return Tokens{
		secretKey: secretKey,
	}
}
