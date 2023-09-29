package db

import "golang.org/x/net/context"

type (
	Store interface {
		SaveRefreshToken(ctx context.Context, uid, token string) error
		GetUserByRefreshToken(ctx context.Context, token string) (string, error)
		IsTokensExist(ctx context.Context, uid string) (bool, error)
		DeleteRefreshTokenByUser(ctx context.Context, uid string) error
		Disconnect(ctx context.Context) error
	}
)
