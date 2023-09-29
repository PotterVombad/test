package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"test/internal/db"
	"test/internal/tokens"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type API struct {
	db     db.Store
	tokens tokens.Repo
}
func (a API) getTokenPairs(c *gin.Context) {
	uid := c.Query("uid")
	exist, err := a.db.IsTokensExist(c.Request.Context(), uid)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if exist {
		c.String(http.StatusOK, "pairs already exists")
		return
	}
	pairs, err := a.tokens.GetPairs(uid)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}


	if err := a.db.SaveRefreshToken(c.Request.Context(), uid, pairs.RefreshToken); err != nil {
		log.Error(err)
		return
	}

	pairs.RefreshToken = string(base64.StdEncoding.EncodeToString(
		[]byte(pairs.RefreshToken),
	))

	c.JSON(http.StatusOK, pairs)
}

func (a API) refresh(c *gin.Context) {
	token, err := base64.StdEncoding.DecodeString(c.Query("token"))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	uid, err := a.db.GetUserByRefreshToken(c.Request.Context(), string(token))
	if err != nil {
		log.Error(fmt.Errorf("refresh by token %s | %s", token, err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if uid == "" {
		c.String(http.StatusOK, "user not found")
		return
	}

	if err := a.db.DeleteRefreshTokenByUser(c.Request.Context(), uid); err != nil {
		log.Error(err)
		return
	}

	c.String(http.StatusOK, "use '/tokens_pairs' to get new pairs")
}

func (a API) Run() error {
	r := gin.Default()
	r.GET("/tokens_pairs", a.getTokenPairs)
	r.PUT("/refresh", a.refresh)

	if err := r.Run(); err != nil {
		panic(err)
	}

	return nil
}

func New(db db.Store, tokens tokens.Repo) API {
	return API{
		db:     db,
		tokens: tokens,
	}
}
