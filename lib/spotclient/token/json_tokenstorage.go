package token

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/m-nny/discsearch/lib/utils"
	"golang.org/x/oauth2"
)

var (
	jsonToken = flag.String("spotify_tokenstorage", "spotify_token.json", "json cache of spotify tokens")
)

type JsonTokenStorage struct {
	tokenMap map[string]*oauth2.Token
}

var _ TokenStorage = (*JsonTokenStorage)(nil)

func (ts *JsonTokenStorage) GetSpotifyToken(ctx context.Context, username string) (*oauth2.Token, error) {
	return ts.tokenMap[username], nil
}

func (ts *JsonTokenStorage) StoreSpotifyToken(ctx context.Context, username string, token *oauth2.Token) error {
	ts.tokenMap[username] = token
	return utils.CJsonSave(*jsonToken, ts.tokenMap)
}

func NewJsonTokenStorage() (*JsonTokenStorage, error) {
	ts := &JsonTokenStorage{tokenMap: make(map[string]*oauth2.Token)}
	if err := utils.CJsonLoad(*jsonToken, &ts.tokenMap); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot load json token storage: %w", err)
	}
	return ts, nil
}
