package mongo

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type (
	MongoDB struct {
		client *mongo.Client
		col    *mongo.Collection
	}

	user struct {
		UID   string `bson:"uid"`
		Token string `bson:"refresh_token"`
	}
)

func (m MongoDB) IsTokensExist(ctx context.Context, uid string) (bool, error) {
	res := m.col.FindOne(ctx, bson.M{
		"uid": uid,
	})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, res.Err()
	}

	return true, nil
}

func (m MongoDB) SaveRefreshToken(ctx context.Context, uid, token string) error {
	hashedRefreshToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate hash | %w", err)
	}

	if _, err := m.col.InsertOne(ctx,
		bson.M{
			"uid":           uid,
			"refresh_token": hashedRefreshToken,
		}); err != nil {
		return fmt.Errorf("col.InsertOne | %w", err)
	}

	return nil
}

func (m MongoDB) GetUserByRefreshToken(ctx context.Context, token string) (string, error) {
	cursor, err := m.col.Find(context.Background(), bson.M{})
	if err != nil {
		return "", err
	}

	defer cursor.Close(context.Background())

	if err := cursor.Err(); err != nil {
		return "", err
	}

	for cursor.Next(context.Background()) {
		var u user
		if err := cursor.Decode(&u); err != nil {
			log.Errorf("error decoding document: %s", err)
			continue
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Token), []byte(token)); err != nil {
			continue
		}

		return u.UID, nil
	}
	return "", nil
}

func (m MongoDB) DeleteRefreshTokenByUser(ctx context.Context, uid string) error {
	if _, err := m.col.DeleteOne(ctx, bson.M{
		"uid": uid,
	}); err != nil {
		return err
	}

	return nil
}

func (m MongoDB) Disconnect(ctx context.Context) error {
	if err := m.client.Disconnect(context.Background()); err != nil {
		return err
	}

	return nil
}

func MustNew(
	ctx context.Context, url, dbName, col string,
) MongoDB {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	if res := client.Database(dbName).RunCommand(
		ctx, bson.D{{
			Key:   "ping",
			Value: 1,
		}},
	); res.Err() != nil {
		panic(res.Err())
	}

	return MongoDB{
		col:    client.Database(dbName).Collection(col),
		client: client,
	}
}
