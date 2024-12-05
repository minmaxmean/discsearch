package main_test

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/m-nny/discsearch/cmd/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	username = "test_username"
)

func Test_Spotify(t *testing.T) {
	ctx := context.Background()
	app := prepApp(t, ctx)
	t.Run("Ping", func(t *testing.T) {
		err := app.SpotifyClient.Ping(ctx)
		assert.NoError(t, err)
	})
	t.Run("SavedTracks", func(t *testing.T) {
		res, err := app.SpotifyClient.SavedTracks(ctx)
		assert.NoError(t, err)
		assert.Greater(t, len(res), 0)
	})
}

var (
	testWithoutCache = flag.Bool("test_nocache", false, "")
	runIntegration   = flag.Bool("run_int_tests", false, "")
)

func prepApp(t *testing.T, ctx context.Context) *common.App {
	if !*runIntegration {
		t.Skip("run_int_tests is false")
	}
	if *testWithoutCache {
		flag.Set("spotify_tstype", "inmemory")
		flag.Set("cache_enabled", "false")
	} else {
		flag.Set("spotify_tstype", "json")
		flag.Set("cache_enabled", "true")
	}

	require.NoError(t, os.Chdir("../.."))
	app, err := common.GetApp(ctx, username)
	require.NoError(t, err)
	return app
}
