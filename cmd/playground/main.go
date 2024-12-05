package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/m-nny/discsearch/lib/spotclient"
	"github.com/m-nny/discsearch/lib/spotclient/token"
	"github.com/m-nny/discsearch/lib/utils"
	"golang.org/x/exp/slog"
)

var (
	username = flag.String("username", "m-nny", "")
)

func main() {
	ctx := context.Background()
	if err := runApp(ctx); err != nil {
		fmt.Printf("Error running app: %q\n", err)
		return
	}
	fmt.Printf("App finished\n")
}

func runApp(ctx context.Context) error {
	if err := utils.LoadFlagsFromEnv(); err != nil {
		return fmt.Errorf("error loading flag: %w", err)
	}
	ts, err := token.GetTokenStorage()
	if err != nil {
		return err
	}
	spotifyClient, err := spotclient.New(ctx, *username, ts)
	if err != nil {
		return err
	}
	if err := spotifyClient.Ping(ctx); err != nil {
		return err
	}
	slog.Debug("runApp", "spotifyClient", spotifyClient)

	return nil
}
