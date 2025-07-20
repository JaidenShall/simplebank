package api

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	db "github.com/JaidenShall/simplebank/db/sqlc"
	"github.com/JaidenShall/simplebank/token"
	"github.com/JaidenShall/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func createAndSetAuthToken(t *testing.T, request *http.Request, tokenMaker token.Maker, username string) {
	if len(username) == 0 {
		return
	}

	token, err := tokenMaker.CreateToken(username, time.Minute)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
