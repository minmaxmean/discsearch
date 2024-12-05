package common

import (
	"context"
	"fmt"

	"github.com/m-nny/discsearch/lib/spotclient"
	"github.com/m-nny/discsearch/lib/spotclient/token"
	"github.com/m-nny/discsearch/lib/utils"
)

type App struct {
	SpotifyClient *spotclient.SpotifyClient
}

func GetApp(ctx context.Context, username string) (*App, error) {
	if err := utils.LoadFlagsFromEnv(); err != nil {
		return nil, fmt.Errorf("error loading flag: %w", err)
	}
	ts, err := token.GetTokenStorage()
	if err != nil {
		return nil, err
	}
	spotifyClient, err := spotclient.New(ctx, username, ts)
	if err != nil {
		return nil, err
	}
	return &App{
		SpotifyClient: spotifyClient,
	}, nil
}
