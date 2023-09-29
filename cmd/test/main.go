package main

import (
	"context"
	"fmt"
	"os"
	"github.com/PotterVombad/test/internal/api"
	"github.com/PotterVombad/test/internal/db/mongo"
	"github.com/PotterVombad/test/internal/tokens"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	url := fmt.Sprintf(
		"mongodb://%s:%s@%s",
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_ADDR"),
	)
	m := mongo.MustNew(ctx, url, os.Getenv("MONGO_DB_NAME"), os.Getenv("MONGO_COL"))

	defer func() {
		if err := m.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	tokens := tokens.New(os.Getenv("JWT_SECRET_KEY"))
	api := api.New(m, tokens)

	log.Infof("start app")
	if err := api.Run(); err != nil {
		panic(err)
	}
}
