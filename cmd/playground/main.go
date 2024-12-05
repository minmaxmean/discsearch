package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/m-nny/discsearch/cmd/common"
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
	app, err := common.GetApp(ctx, *username)
	if err != nil {
		return err
	}
	tracks, err := app.SpotifyClient.SavedTracks(ctx)
	if err != nil {
		return err
	}
	slog.Debug("saved tracks", "tracks", tracks)

	return nil
}
