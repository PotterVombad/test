package memory

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type (
	MemoryDB struct {
		db map[string]string
	}
)

func (m MemoryDB) Disconnect(ctx context.Context) error {
	return nil
}

func (m MemoryDB) IsTokensExist(ctx context.Context, uid string) (bool, error) {
	_, ok := m.db[uid]
	return ok, nil
}

func (m MemoryDB) SaveRefreshToken(ctx context.Context, uid, token string) error {
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate hash | %w", err)
	}

	m.db[uid] = string(hashedRefreshToken)

	return nil
}

func (m MemoryDB) GetUserByRefreshToken(ctx context.Context, token string) (string, error) {
	for uid, t := range m.db {
		if err := bcrypt.CompareHashAndPassword([]byte(t), []byte(token)); err == nil {
			return uid, nil
		}
	}

	return "", nil
}

func (m MemoryDB) DeleteRefreshTokenByUser(ctx context.Context, uid string) error {
	delete(m.db, uid)
	return nil
}

func New() MemoryDB {
	return MemoryDB{
		db: make(map[string]string),
	}
}
